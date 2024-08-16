package masa

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	shell "github.com/ipfs/go-ipfs-api"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/masa-finance/masa-oracle/internal/versioning"
	"github.com/masa-finance/masa-oracle/pkg/chain"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type OracleNode struct {
	Host              host.Host
	PrivKey           *ecdsa.PrivateKey
	Protocol          protocol.ID
	priorityAddrs     multiaddr.Multiaddr
	multiAddrs        []multiaddr.Multiaddr
	DHT               *dht.IpfsDHT
	Context           context.Context
	PeerChan          chan myNetwork.PeerEvent
	NodeTracker       *pubsub2.NodeEventTracker
	PubSubManager     *pubsub2.Manager
	Signature         string
	IsStaked          bool
	IsValidator       bool
	IsTwitterScraper  bool
	IsDiscordScraper  bool
	IsTelegramScraper bool
	IsWebScraper      bool
	IsLlmServer       bool
	StartTime         time.Time
	WorkerTracker     *pubsub2.WorkerEventTracker
	BlockTracker      *BlockEventTracker
	Blockchain        *chain.Chain
}

// GetMultiAddrs returns the priority multiaddr for this node.
// It first checks if the priority address is already set, and returns it if so.
// If not, it determines the priority address from the available multiaddrs using
// the GetPriorityAddress utility function, sets it, and returns it.
func (node *OracleNode) GetMultiAddrs() multiaddr.Multiaddr {
	if node.priorityAddrs == nil {
		pAddr := myNetwork.GetPriorityAddress(node.multiAddrs)
		node.priorityAddrs = pAddr
	}
	return node.priorityAddrs
}

// getOutboundIP is a function that returns the outbound IP address of the current machine as a string.
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("[-] Error getting outbound IP")
		return ""
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

// NewOracleNode creates a new OracleNode instance with the provided context and
// staking status. It initializes the libp2p host, DHT, pubsub manager, and other
// components needed for an Oracle node to join the network and participate.
func NewOracleNode(ctx context.Context, isStaked bool) (*OracleNode, error) {
	// Start with the default scaling limits.
	cfg := config.GetInstance()
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	resourceManager, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}

	var addrStr []string
	libp2pOptions := []libp2p.Option{
		libp2p.Identity(masacrypto.KeyManagerInstance().Libp2pPrivKey),
		libp2p.ResourceManager(resourceManager),
		libp2p.Ping(false), // disable built-in ping
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(), // Enable Circuit Relay v2 with hop
	}

	securityOptions := []libp2p.Option{
		libp2p.Security(noise.ID, noise.New),
	}
	// @note fix for increase buffer size warning on linux
	// sudo sysctl -w net.core.rmem_max=7500000
	// sudo sysctl -w net.core.wmem_max=7500000
	// sudo sysctl -p
	if cfg.UDP {
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", cfg.PortNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(quic.NewTransport))
	}
	if cfg.TCP {
		securityOptions = append(securityOptions, libp2p.Security(libp2ptls.ID, libp2ptls.New))
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.PortNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(tcp.NewTCPTransport))
		libp2pOptions = append(libp2pOptions, libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport))
	}
	libp2pOptions = append(libp2pOptions, libp2p.ChainOptions(securityOptions...))
	libp2pOptions = append(libp2pOptions, libp2p.ListenAddrStrings(addrStr...))

	hst, err := libp2p.New(libp2pOptions...)
	if err != nil {
		return nil, err
	}

	subscriptionManager, err := pubsub2.NewPubSubManager(ctx, hst)
	if err != nil {
		return nil, err
	}

	return &OracleNode{
		Host:              hst,
		PrivKey:           masacrypto.KeyManagerInstance().EcdsaPrivKey,
		Protocol:          config.ProtocolWithVersion(config.OracleProtocol),
		multiAddrs:        myNetwork.GetMultiAddressesForHostQuiet(hst),
		Context:           ctx,
		PeerChan:          make(chan myNetwork.PeerEvent),
		NodeTracker:       pubsub2.NewNodeEventTracker(versioning.ProtocolVersion, cfg.Environment, hst.ID().String()),
		PubSubManager:     subscriptionManager,
		IsStaked:          isStaked,
		IsValidator:       cfg.Validator,
		IsTwitterScraper:  cfg.TwitterScraper,
		IsDiscordScraper:  cfg.DiscordScraper,
		IsTelegramScraper: cfg.TelegramScraper,
		IsWebScraper:      cfg.WebScraper,
		IsLlmServer:       cfg.LlmServer,
		Blockchain:        &chain.Chain{},
	}, nil
}

// Start initializes the OracleNode by setting up libp2p stream handlers,
// connecting to the DHT and bootnodes, and subscribing to topics. It launches
// goroutines to handle discovered peers, listen to the node tracker, and
// discover peers. If this is a bootnode, it adds itself to the node tracker.
func (node *OracleNode) Start() (err error) {
	logrus.Infof("[+] Starting node with ID: %s", node.GetMultiAddrs().String())

	bootNodeAddrs, err := myNetwork.GetBootNodesMultiAddress(config.GetInstance().Bootnodes)
	if err != nil {
		return err
	}

	node.Host.SetStreamHandler(node.Protocol, node.handleStream)
	node.Host.SetStreamHandler(config.ProtocolWithVersion(config.NodeDataSyncProtocol), node.ReceiveNodeData)

	if node.IsStaked {
		node.Host.SetStreamHandler(config.ProtocolWithVersion(config.NodeGossipTopic), node.GossipNodeData)
	}
	node.Host.Network().Notify(node.NodeTracker)

	go node.ListenToNodeTracker()
	go node.handleDiscoveredPeers()

	node.DHT, err = myNetwork.WithDht(node.Context, node.Host, bootNodeAddrs, node.Protocol, config.MasaPrefix, node.PeerChan, node.IsStaked)
	if err != nil {
		return err
	}
	err = myNetwork.WithMDNS(node.Host, config.Rendezvous, node.PeerChan)
	if err != nil {
		return err
	}

	go myNetwork.Discover(node.Context, node.Host, node.DHT, node.Protocol)

	nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
	if nodeData == nil {
		publicKeyHex := masacrypto.KeyManagerInstance().EthAddress
		ma := myNetwork.GetMultiAddressesForHostQuiet(node.Host)
		nodeData = pubsub2.NewNodeData(ma[0], node.Host.ID(), publicKeyHex, pubsub2.ActivityJoined)
		nodeData.IsStaked = node.IsStaked
		nodeData.SelfIdentified = true
		nodeData.IsDiscordScraper = node.IsDiscordScraper
		nodeData.IsTelegramScraper = node.IsTelegramScraper
		nodeData.IsTwitterScraper = node.IsTwitterScraper
		nodeData.IsWebScraper = node.IsWebScraper
		nodeData.IsValidator = node.IsValidator
	}

	nodeData.Joined()
	node.NodeTracker.HandleNodeData(*nodeData)

	// call SubscribeToTopics on startup
	if err := SubscribeToTopics(node); err != nil {
		return err
	}
	node.StartTime = time.Now()

	return nil
}

// handleDiscoveredPeers listens on the PeerChan for discovered peers from the
// network discovery routines. It handles connecting to new peers and closing
// connections to peers that disconnect. This runs continuously to handle
// discovered peers.
func (node *OracleNode) handleDiscoveredPeers() {
	for {
		select {
		case peer := <-node.PeerChan: // will block until we discover a peer
			logrus.Debugf("[+] Peer Event for: %s, Action: %s", peer.AddrInfo.ID.String(), peer.Action)
			// If the peer is a new peer, connect to it
			if peer.Action == myNetwork.PeerAdded {
				if err := node.Host.Connect(node.Context, peer.AddrInfo); err != nil {
					logrus.Errorf("[-] Connection failed for peer: %s %v", peer.AddrInfo.ID.String(), err)
					// close the connection
					err := node.Host.Network().ClosePeer(peer.AddrInfo.ID)
					if err != nil {
						logrus.Error("[-] Error closing peer: ", err)
					}
					continue
				}
			}
		case <-node.Context.Done():
			return
		}
	}
}

// handleStream handles an incoming libp2p stream from a remote peer.
// It reads the stream data, validates the remote peer ID, updates the node tracker
// with the remote peer's information, and logs the event.
func (node *OracleNode) handleStream(stream network.Stream) {
	remotePeer, nodeData, err := node.handleStreamData(stream)
	if err != nil {
		if strings.HasPrefix(err.Error(), "un-staked") {
			// just ignore the error
			return
		}
		logrus.Errorf("[-] Failed to read stream: %v", err)
		return
	}
	if remotePeer.String() != nodeData.PeerId.String() {
		logrus.Warnf("[-] Received data from unexpected peer %s", remotePeer)
		return
	}

	multiAddr := stream.Conn().RemoteMultiaddr()
	nodeData.Multiaddrs = []pubsub2.JSONMultiaddr{{Multiaddr: multiAddr}}

	// newNodeData := pubsub2.NewNodeData(multiAddr, remotePeer, nodeData.EthAddress, pubsub2.ActivityJoined)
	// newNodeData.IsStaked = nodeData.IsStaked
	err = node.NodeTracker.AddOrUpdateNodeData(&nodeData, false)
	if err != nil {
		logrus.Error("[-] Error adding or updating node data: ", err)
		return
	}
	logrus.Infof("[+] nodeStream -> Received data from: %s", remotePeer.String())
}

// IsWorker determines if the OracleNode is configured to act as an actor.
// An actor node is one that has at least one of the following scrapers enabled:
// TwitterScraper, DiscordScraper, or WebScraper.
// It returns true if any of these scrapers are enabled, otherwise false.
func (node *OracleNode) IsWorker() bool {
	// need to get this by node data
	cfg := config.GetInstance()
	if cfg.TwitterScraper || cfg.DiscordScraper || cfg.TelegramScraper || cfg.WebScraper {
		return true
	}
	return false
}

// IsPublisher returns true if this node is a publisher node.
// A publisher node is one that has a non-empty signature.
func (node *OracleNode) IsPublisher() bool {
	// Node is a publisher if it has a non-empty signature
	return node.Signature != ""
}

// FromUnixTime converts a Unix timestamp into a formatted string.
// The Unix timestamp is expected to be in seconds.
// The returned string is in the format "2006-01-02T15:04:05.000Z".
func (node *OracleNode) FromUnixTime(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02T15:04:05.000Z")
}

// ToUnixTime converts a formatted string time into a Unix timestamp.
// The input string is expected to be in the format "2006-01-02T15:04:05.000Z".
// The returned Unix timestamp is in seconds.
func (node *OracleNode) ToUnixTime(stringTime string) int64 {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", stringTime)
	return t.Unix()
}

// Version returns the current version string of the oracle node software.
func (node *OracleNode) Version() string {
	return config.GetInstance().Version
}

// LogActiveTopics logs the currently active topic names to the
// default logger. It gets the list of active topics from the
// PubSubManager and logs them if there are any, otherwise it logs
// that there are no active topics.
func (node *OracleNode) LogActiveTopics() {
	topicNames := node.PubSubManager.GetTopicNames()
	if len(topicNames) > 0 {
		logrus.Infof("[+] Active topics: %v", topicNames)
	} else {
		logrus.Info("[-] No active topics.")
	}
}

// Blockchain Implementation
var (
	blocksCh = make(chan *pubsub.Message)
)

type BlockData struct {
	Block            uint64      `json:"block"`
	InputData        interface{} `json:"input_data"`
	TransactionHash  string      `json:"transaction_hash"`
	PreviousHash     string      `json:"previous_hash"`
	TransactionNonce int         `json:"nonce"`
}

type Blocks struct {
	BlockData []BlockData `json:"blocks"`
}

type BlockEvents struct{}

type BlockEventTracker struct {
	BlockEvents []BlockEvents
	BlockTopic  *pubsub.Topic
	mu          sync.Mutex
}

// HandleMessage processes incoming pubsub messages containing block events.
// It unmarshals the message data into a slice of BlockEvents and appends them
// to the tracker's BlockEvents slice.
func (b *BlockEventTracker) HandleMessage(m *pubsub.Message) {
	var blockEvents any

	// Try to decode as base64 first
	decodedData, err := base64.StdEncoding.DecodeString(string(m.Data))
	if err == nil {
		m.Data = decodedData
	}

	// Try to unmarshal as JSON
	err = json.Unmarshal(m.Data, &blockEvents)
	if err != nil {
		// If JSON unmarshal fails, try to interpret as string
		blockEvents = string(m.Data)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	switch v := blockEvents.(type) {
	case []BlockEvents:
		b.BlockEvents = append(b.BlockEvents, v...)
	case BlockEvents:
		b.BlockEvents = append(b.BlockEvents, v)
	case map[string]interface{}:
		// Convert map to BlockEvents struct
		newBlockEvent := BlockEvents{}
		// You might need to add logic here to properly convert the map to BlockEvents
		b.BlockEvents = append(b.BlockEvents, newBlockEvent)
	case []interface{}:
		// Convert each item in the slice to BlockEvents
		for _, item := range v {
			if be, ok := item.(BlockEvents); ok {
				b.BlockEvents = append(b.BlockEvents, be)
			}
		}
	case string:
		// Handle string data
		newBlockEvent := BlockEvents{}
		// You might need to add logic here to properly convert the string to BlockEvents
		b.BlockEvents = append(b.BlockEvents, newBlockEvent)
	default:
		logrus.Warnf("[-] Unexpected data type in message: %v", reflect.TypeOf(v))
	}

	blocksCh <- m
}

func updateBlocks(ctx context.Context, node *OracleNode) error {

	var existingBlocks Blocks
	blocks := chain.GetBlockchain(node.Blockchain)

	for _, block := range blocks {
		var inputData interface{}
		err := json.Unmarshal(block.Data, &inputData)
		if err != nil {
			inputData = string(block.Data) // Fallback to string if unmarshal fails
		}

		blockData := BlockData{
			Block:            block.Block,
			InputData:        base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", inputData))),
			TransactionHash:  fmt.Sprintf("%x", block.Hash),
			PreviousHash:     fmt.Sprintf("%x", block.Link),
			TransactionNonce: int(block.Nonce),
		}
		existingBlocks.BlockData = append(existingBlocks.BlockData, blockData)
	}
	jsonData, err := json.Marshal(existingBlocks)
	if err != nil {
		return err
	}

	err = node.DHT.PutValue(ctx, "/db/blocks", jsonData)
	if err != nil {
		logrus.Warningf("[-] Unable to store block on DHT: %v", err)
	}

	if os.Getenv("IPFS_URL") != "" {

		infuraURL := fmt.Sprintf("https://%s:%s@%s", os.Getenv("PID"), os.Getenv("PS"), os.Getenv("IPFS_URL"))
		sh := shell.NewShell(infuraURL)

		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			logrus.Errorf("[-] Error marshalling JSON: %s", err)
		}

		reader := bytes.NewReader(jsonBytes)

		hash, err := sh.AddWithOpts(reader, true, true)
		if err != nil {
			logrus.Errorf("[-] Error persisting to IPFS: %s", err)
		} else {
			logrus.Printf("[+] Ledger persisted with IPFS hash: https://dwn.infura-ipfs.io/ipfs/%s\n", hash)
			_ = node.DHT.PutValue(ctx, "/db/ipfs", []byte(fmt.Sprintf("https://dwn.infura-ipfs.io/ipfs/%s", hash)))

		}
	}

	return nil
}

func SubscribeToBlocks(ctx context.Context, node *OracleNode) {
	if !node.IsValidator {
		return
	}

	go func() {
		err := node.Blockchain.Init()
		if err != nil {
			logrus.Error(err)
		}
	}()

	updateTicker := time.NewTicker(time.Second * 60)
	defer updateTicker.Stop()

	for {
		select {
		case block, ok := <-blocksCh:
			if !ok {
				logrus.Error("[-] Block channel closed")
				return
			}
			if err := processBlock(node, block); err != nil {
				logrus.Errorf("[-] Error processing block: %v", err)
				// Consider adding a retry mechanism or circuit breaker here
			}

		case <-updateTicker.C:
			logrus.Info("[+] blockchain tick")
			if err := updateBlocks(ctx, node); err != nil {
				logrus.Errorf("[-] Error updating blocks: %v", err)
				// Consider adding a retry mechanism or circuit breaker here
			}

		case <-ctx.Done():
			logrus.Info("[+] Context cancelled, stopping block subscription")
			return
		}
	}
}

func processBlock(node *OracleNode, block *pubsub.Message) error {
	blocks := chain.GetBlockchain(node.Blockchain)
	for _, b := range blocks {
		if string(b.Data) == string(block.Data) {
			return nil // Block already exists
		}
	}

	if err := node.Blockchain.AddBlock(block.Data); err != nil {
		return fmt.Errorf("[-] failed to add block: %w", err)
	}

	if node.Blockchain.LastHash != nil {
		b, err := node.Blockchain.GetBlock(node.Blockchain.LastHash)
		if err != nil {
			return fmt.Errorf("[-] failed to get last block: %w", err)
		}
		b.Print()
	}

	return nil
}
