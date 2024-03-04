package masa

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/masa-finance/masa-oracle/pkg/ad"
	crypto2 "github.com/masa-finance/masa-oracle/pkg/crypto"
	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type OracleNode struct {
	Host                  host.Host
	PrivKey               *ecdsa.PrivateKey
	Protocol              protocol.ID
	priorityAddrs         multiaddr.Multiaddr
	multiAddrs            []multiaddr.Multiaddr
	DHT                   *dht.IpfsDHT
	Context               context.Context
	PeerChan              chan myNetwork.PeerEvent
	NodeTracker           *pubsub2.NodeEventTracker
	PubSubManager         *pubsub2.Manager
	Signature             string
	IsStaked              bool
	StartTime             time.Time
	AdSubscriptionHandler *ad.SubscriptionHandler
}

func (node *OracleNode) GetMultiAddrs() multiaddr.Multiaddr {
	if node.priorityAddrs == nil {
		pAddr := myNetwork.GetPriorityAddress(node.multiAddrs)
		node.priorityAddrs = pAddr
	}
	return node.priorityAddrs
}

func NewOracleNode(ctx context.Context, privKey crypto.PrivKey, portNbr int, useUdp, useTcp bool, isStaked bool) (*OracleNode, error) {
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
		libp2p.Identity(privKey),
		libp2p.ResourceManager(resourceManager),
		libp2p.Ping(false), // disable built-in ping
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(), // Enable Circuit Relay v2 with hop
	}

	securityOptions := []libp2p.Option{
		libp2p.Security(noise.ID, noise.New),
	}
	if useUdp {
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", portNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(quic.NewTransport))
	}
	if useTcp {
		securityOptions = append(securityOptions, libp2p.Security(libp2ptls.ID, libp2ptls.New))
		addrStr = append(addrStr, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", portNbr))
		libp2pOptions = append(libp2pOptions, libp2p.Transport(tcp.NewTCPTransport))
		libp2pOptions = append(libp2pOptions, libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport))
	}
	libp2pOptions = append(libp2pOptions, libp2p.ChainOptions(securityOptions...))
	libp2pOptions = append(libp2pOptions, libp2p.ListenAddrStrings(addrStr...))

	hst, err := libp2p.New(libp2pOptions...)
	if err != nil {
		return nil, err
	}

	// Extract the public key from the private key
	pubKey := privKey.GetPublic()
	if err != nil {
		return nil, err
	}

	// Pass the public key as the third argument to NewPubSubManager
	subscriptionManager, err := pubsub2.NewPubSubManager(ctx, hst, pubKey)
	if err != nil {
		return nil, err
	}

	ecdsaPrivKey, err := crypto2.Libp2pPrivateKeyToEcdsa(privKey)
	if err != nil {
		return nil, err
	}
	return &OracleNode{
		Host:          hst,
		PrivKey:       ecdsaPrivKey,
		Protocol:      ProtocolWithVersion(oracleProtocol),
		multiAddrs:    myNetwork.GetMultiAddressesForHostQuiet(hst),
		Context:       ctx,
		PeerChan:      make(chan myNetwork.PeerEvent),
		NodeTracker:   pubsub2.NewNodeEventTracker(Version, viper.GetString(Environment)),
		PubSubManager: subscriptionManager,
		IsStaked:      isStaked,
	}, nil
}

func (node *OracleNode) Start() (err error) {
	logrus.Infof("Starting node with ID: %s", node.GetMultiAddrs().String())

	bootNodeAddrs, err := myNetwork.GetBootNodesMultiAddress(viper.GetString(BootNodes))
	if err != nil {
		return err
	}

	node.Host.SetStreamHandler(node.Protocol, node.handleStream)
	node.Host.SetStreamHandler(ProtocolWithVersion(NodeDataSyncProtocol), node.ReceiveNodeData)
	// if node.IsStaked then allow them to be added to the NodeData -- move to node tracker
	if node.IsStaked {
		node.Host.SetStreamHandler(ProtocolWithVersion(NodeGossipTopic), node.GossipNodeData)
	}
	node.Host.Network().Notify(node.NodeTracker)

	go node.ListenToNodeTracker()
	go node.handleDiscoveredPeers()

	node.DHT, err = myNetwork.WithDht(node.Context, node.Host, bootNodeAddrs, node.Protocol, masaPrefix, node.PeerChan, node.IsStaked)
	if err != nil {
		return err
	}
	err = myNetwork.WithMDNS(node.Host, rendezvous, node.PeerChan)
	if err != nil {
		return err
	}

	go myNetwork.Discover(node.Context, node.Host, node.DHT, node.Protocol)
	// if this is the original boot node then add it to the node tracker
	if viper.GetString(BootNodes) == "" {
		nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
		if nodeData == nil {
			publicKeyHex, _ := crypto2.GetPublicKeyForHost(node.Host)
			nodeData = pubsub2.NewNodeData(node.GetMultiAddrs(), node.Host.ID(), publicKeyHex, pubsub2.ActivityJoined)
			nodeData.IsStaked = node.IsStaked
			nodeData.SelfIdentified = true
		}
		nodeData.Joined()
		node.NodeTracker.HandleNodeData(*nodeData)
	}
	// call SubscribeToTopics on startup
	if err := SubscribeToTopics(node); err != nil {
		return err
	}
	node.StartTime = time.Now()

	return nil
}

func (node *OracleNode) handleDiscoveredPeers() {
	for {
		select {
		case peer := <-node.PeerChan: // will block until we discover a peer
			logrus.Debugf("Peer Event for: %s, Action: %s", peer.AddrInfo.ID.String(), peer.Action)
			// If the peer is a new peer, connect to it
			if peer.Action == myNetwork.PeerAdded {
				if err := node.Host.Connect(node.Context, peer.AddrInfo); err != nil {
					logrus.Errorf("Connection failed for peer: %s %v", peer.AddrInfo.ID.String(), err)
					// close the connection
					err := node.Host.Network().ClosePeer(peer.AddrInfo.ID)
					if err != nil {
						logrus.Error(err)
					}
					continue
				}
			}
		case <-node.Context.Done():
			return
		}
	}
}

func (node *OracleNode) handleStream(stream network.Stream) {
	remotePeer, nodeData, err := node.handleStreamData(stream)
	if err != nil {
		if strings.HasPrefix(err.Error(), "un-staked") {
			// just ignore the error
			return
		}
		logrus.Errorf("Failed to read stream: %v", err)
		return
	}
	if remotePeer.String() != nodeData.PeerId.String() {
		logrus.Warnf("Received data from unexpected peer %s", remotePeer)
		return
	}
	multiAddr := stream.Conn().RemoteMultiaddr()
	newNodeData := pubsub2.NewNodeData(multiAddr, remotePeer, nodeData.EthAddress, pubsub2.ActivityJoined)
	newNodeData.IsStaked = nodeData.IsStaked
	err = node.NodeTracker.AddOrUpdateNodeData(newNodeData, false)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("handleStream -> Received data from:", remotePeer.String())
}

func (node *OracleNode) IsPublisher() bool {
	// Node is a publisher if it has a non-empty signature
	return node.Signature != ""
}

func (node *OracleNode) Version() string {
	return Version
}

func (node *OracleNode) LogActiveTopics() {
	topicNames := node.PubSubManager.GetTopicNames()
	if len(topicNames) > 0 {
		logrus.Infof("Active topics: %v", topicNames)
	} else {
		logrus.Info("No active topics.")
	}
}
