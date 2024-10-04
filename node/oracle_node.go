package node

import (
	"context"
	"fmt"
	"strings"
	"time"

	ethereumCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/internal/versioning"
	"github.com/masa-finance/masa-oracle/pkg/chain"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type OracleNode struct {
	Host          host.Host
	Protocol      protocol.ID
	priorityAddrs multiaddr.Multiaddr
	multiAddrs    []multiaddr.Multiaddr
	DHT           *dht.IpfsDHT
	PeerChan      chan myNetwork.PeerEvent
	NodeTracker   *pubsub.NodeEventTracker
	PubSubManager *pubsub.Manager
	Signature     string
	StartTime     time.Time
	WorkerTracker *pubsub.WorkerEventTracker
	Blockchain    *chain.Chain
	Options       NodeOption
	Context       context.Context
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

// GetP2PMultiAddrs returns the multiaddresses for the host in P2P format.
func (node *OracleNode) GetP2PMultiAddrs() ([]multiaddr.Multiaddr, error) {
	addrs := node.Host.Addrs()
	pi := peer.AddrInfo{
		ID:    node.Host.ID(),
		Addrs: addrs,
	}

	return peer.AddrInfoToP2pAddrs(&pi)
}

// NewOracleNode creates a new OracleNode instance with the provided context and
// staking status. It initializes the libp2p host, DHT, pubsub manager, and other
// components needed for an Oracle node to join the network and participate.
func NewOracleNode(ctx context.Context, opts ...Option) (*OracleNode, error) {
	o := &NodeOption{}
	o.Apply(opts...)

	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	resourceManager, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}

	var addrStr []string
	libp2pOptions := []libp2p.Option{
		libp2p.ResourceManager(resourceManager),
		libp2p.Ping(false), // disable built-in ping
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(), // Enable Circuit Relay v2 with hop
	}

	if o.RandomIdentity {
		libp2pOptions = append(libp2pOptions, libp2p.RandomIdentity)
	} else {
		libp2pOptions = append(libp2pOptions, libp2p.Identity(masacrypto.KeyManagerInstance().Libp2pPrivKey))
	}

	securityOptions := []libp2p.Option{
		libp2p.Security(noise.ID, noise.New),
	}
	// @note fix for increase buffer size warning on linux
	// sudo sysctl -w net.core.rmem_max=7500000
	// sudo sysctl -w net.core.wmem_max=7500000
	// sudo sysctl -p
	if o.UDP {
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", o.PortNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(quic.NewTransport))
	}
	if o.TCP {
		securityOptions = append(securityOptions, libp2p.Security(libp2ptls.ID, libp2ptls.New))
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", o.PortNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(tcp.NewTCPTransport))
		libp2pOptions = append(libp2pOptions, libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport))
	}
	libp2pOptions = append(libp2pOptions, libp2p.ChainOptions(securityOptions...))
	libp2pOptions = append(libp2pOptions, libp2p.ListenAddrStrings(addrStr...))

	hst, err := libp2p.New(libp2pOptions...)
	if err != nil {
		return nil, err
	}

	subscriptionManager, err := pubsub.NewPubSubManager(ctx, hst)
	if err != nil {
		return nil, err
	}

	n := &OracleNode{
		Host:          hst,
		multiAddrs:    myNetwork.GetMultiAddressesForHostQuiet(hst),
		PeerChan:      make(chan myNetwork.PeerEvent),
		NodeTracker:   pubsub.NewNodeEventTracker(versioning.ProtocolVersion, o.Environment, hst.ID().String()),
		Context:       ctx,
		PubSubManager: subscriptionManager,
		Blockchain:    &chain.Chain{},
		Options:       *o,
	}

	n.Protocol = n.protocolWithVersion(config.OracleProtocol)
	return n, nil
}

func (node *OracleNode) generateEthHexKeyForRandomIdentity() (string, error) {
	// If it's a random identity, get the pubkey from Libp2p
	// and convert these to Ethereum public Hex types
	pubkey, err := node.Host.ID().ExtractPublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to extract public key from p2p identity: %w", err)
	}
	rawKey, err := pubkey.Raw()
	if err != nil {
		return "", fmt.Errorf("failed to extract public key from p2p identity: %w", err)
	}
	return common.BytesToAddress(ethereumCrypto.Keccak256(rawKey[1:])[12:]).Hex(), nil
}

func (node *OracleNode) getNodeData(host host.Host, addr multiaddr.Multiaddr, publicEthAddress string) *pubsub.NodeData {
	// GetSelfNodeData converts the local node's data into a JSON byte array.
	// It populates a NodeData struct with the node's ID, staking status, and Ethereum address.
	// The NodeData struct is then marshalled into a JSON byte array.
	// Returns nil if there is an error marshalling to JSON.
	// Create and populate NodeData
	nodeData := pubsub.NewNodeData(addr, host.ID(), publicEthAddress, pubsub.ActivityJoined)
	nodeData.MultiaddrsString = addr.String()
	nodeData.IsStaked = node.Options.IsStaked
	nodeData.IsTwitterScraper = node.Options.IsTwitterScraper
	nodeData.IsDiscordScraper = node.Options.IsDiscordScraper
	nodeData.IsTelegramScraper = node.Options.IsLlmServer
	nodeData.IsWebScraper = node.Options.IsWebScraper
	nodeData.IsValidator = node.Options.IsValidator
	nodeData.IsActive = true
	nodeData.Version = versioning.ProtocolVersion

	return nodeData
}

// Start initializes the OracleNode by setting up libp2p stream handlers,
// connecting to the DHT and bootnodes, and subscribing to topics. It launches
// goroutines to handle discovered peers, listen to the node tracker, and
// discover peers. If this is a bootnode, it adds itself to the node tracker.
func (node *OracleNode) Start() (err error) {
	logrus.Infof("[+] Starting node with ID: %s", node.GetMultiAddrs().String())

	node.Host.SetStreamHandler(node.Protocol, node.handleStream)
	node.Host.SetStreamHandler(node.protocolWithVersion(config.NodeDataSyncProtocol), node.ReceiveNodeData)

	for pid, n := range node.Options.ProtocolHandlers {
		node.Host.SetStreamHandler(pid, n)
	}

	for protocol, n := range node.Options.MasaProtocolHandlers {
		node.Host.SetStreamHandler(node.protocolWithVersion(protocol), n)
	}

	if node.Options.IsStaked {
		node.Host.SetStreamHandler(node.protocolWithVersion(config.NodeGossipTopic), node.GossipNodeData)
	}

	node.Host.Network().Notify(node.NodeTracker)

	go node.ListenToNodeTracker()
	go node.handleDiscoveredPeers()
	go node.NodeTracker.ClearExpiredWorkerTimeouts()

	var publicKeyHex string
	if node.Options.RandomIdentity {
		publicKeyHex, _ = node.generateEthHexKeyForRandomIdentity()
	} else {
		publicKeyHex = masacrypto.KeyManagerInstance().EthAddress
	}

	myNodeData := node.getNodeData(node.Host, node.priorityAddrs, publicKeyHex)

	bootstrapNodes, err := myNetwork.GetBootNodesMultiAddress(node.Options.Bootnodes)
	if err != nil {
		return err
	}

	node.DHT, err = myNetwork.WithDHT(node.Context, node.Host, bootstrapNodes, node.Protocol, masaPrefix, node.PeerChan, myNodeData)
	if err != nil {
		return err
	}

	err = myNetwork.WithMDNS(node.Host, config.Rendezvous, node.PeerChan)
	if err != nil {
		return err
	}

	for _, p := range node.Options.Services {
		go p(node.Context, node)
	}

	go myNetwork.Discover(node.Context, node.Host, node.DHT, node.Protocol)

	nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
	if nodeData == nil {
		nodeData = myNodeData
		nodeData.SelfIdentified = true
	}
	nodeData.Joined(node.Options.Version)
	node.NodeTracker.HandleNodeData(*nodeData)

	// call SubscribeToTopics on startup
	if err := node.subscribeToTopics(); err != nil {
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
	defer func(stream network.Stream) {
		err := stream.Close()
		if err != nil {
			logrus.Infof("[-] Error closing stream: %v", err)
		}
	}(stream)

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
	nodeData.MergeMultiaddresses(stream.Conn().RemoteMultiaddr())

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
