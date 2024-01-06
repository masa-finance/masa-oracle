package masa

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/ad"
	crypto2 "github.com/masa-finance/masa-oracle/pkg/crypto"
	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type OracleNode struct {
	Host          host.Host
	PrivKey       *ecdsa.PrivateKey
	Protocol      protocol.ID
	priorityAddrs multiaddr.Multiaddr
	multiAddrs    []multiaddr.Multiaddr
	DHT           *dht.IpfsDHT
	Context       context.Context
	PeerChan      chan myNetwork.PeerEvent
	NodeTracker   *pubsub2.NodeEventTracker
	PubSubManager *pubsub2.Manager
	Signature     string
	IsStaked      bool
	StartTime     time.Time
	IDService     identify.IDService
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

	subscriptionManager, err := pubsub2.NewPubSubManager(ctx, hst)
	if err != nil {
		return nil, err
	}

	// Create a new Identify service
	ids, err := identify.NewIDService(hst)
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
		NodeTracker:   pubsub2.NewNodeEventTracker(Version),
		PubSubManager: subscriptionManager,
		IsStaked:      isStaked,
		IDService:     ids,
	}, nil
}

func (node *OracleNode) Start() (err error) {
	logrus.Infof("Starting node with ID: %s", node.GetMultiAddrs().String())
	node.Host.SetStreamHandler(node.Protocol, node.handleStream)
	node.Host.SetStreamHandler(ProtocolWithVersion(NodeDataSyncProtocol), node.ReceiveNodeData)
	//if node.IsStaked then allow them to be added to the NodeData -- move to node tracker
	if node.IsStaked {
		node.Host.SetStreamHandler(ProtocolWithVersion(NodeGossipTopic), node.GossipNodeData)
	}
	node.Host.Network().Notify(node.NodeTracker)

	go node.ListenToNodeTracker()
	go node.handleDiscoveredPeers()

	err = myNetwork.WithMDNS(node.Host, rendezvous, node.PeerChan)
	if err != nil {
		return err
	}

	bootNodeAddrs, err := myNetwork.GetBootNodesMultiAddress(os.Getenv(Peers))
	if err != nil {
		return err
	}

	node.DHT, err = myNetwork.WithDht(node.Context, node.Host, bootNodeAddrs, node.Protocol, masaPrefix, node.PeerChan, node.IsStaked)
	if err != nil {
		return err
	}

	go myNetwork.Discover(node.Context, node.Host, node.DHT, node.Protocol, node.GetMultiAddrs())

	// Subscribe to a topics
	err = node.PubSubManager.AddSubscription(TopiclWithVersion(NodeGossipTopic), node.NodeTracker)
	if err != nil {
		return err
	}
	err = node.PubSubManager.AddSubscription(TopiclWithVersion(AdTopic), &ad.SubscriptionHandler{})
	node.StartTime = time.Now()
	return nil
}

func (node *OracleNode) handleDiscoveredPeers() {
	for {
		select {
		case peer := <-node.PeerChan: // will block until we discover a peer
			logrus.Info("Peer Event for:", peer, ", Action:", peer.Action)

			if err := node.Host.Connect(node.Context, peer.AddrInfo); err != nil {
				logrus.Error("Connection failed:", err)
				continue
			}

			//open a stream, this stream will be handled by handleStream other end
			stream, err := node.Host.NewStream(node.Context, peer.AddrInfo.ID, node.Protocol)
			if err != nil {
				logrus.Error("Stream open failed", err)
			}
			sendData := pubsub2.GetSelfNodeDataJson(node.Host, node.IsStaked)
			_, err = stream.Write(sendData)
			if err != nil {
				logrus.Error("Stream write failed", err)
				return
			}
		case <-node.Context.Done():
			return
		}
	}
}

func (node *OracleNode) handleStream(stream network.Stream) {
	message, err := node.handleStreamData(stream)
	if err != nil {
		logrus.Error("Error handling stream data:", err)
		return
	}
	// Deserialize the JSON message to a NodeData struct
	var nodeData pubsub2.NodeData
	err = json.Unmarshal(message, &nodeData)
	if err != nil {
		logrus.Errorf("Failed to unmarshal NodeData: %v", err)
		logrus.Errorf("%s", string(message))
		return
	}
	remotePeer := stream.Conn().RemotePeer()
	if remotePeer.String() != nodeData.PeerId.String() {
		logrus.Warnf("Received data from unexpected peer %s", remotePeer)
		return
	}
	// Store the IsStaked status in the map
	node.NodeTracker.IsStakedCond.L.Lock()
	node.NodeTracker.IsStakedStatus[stream.Conn().RemotePeer().String()] = nodeData.IsStaked
	node.NodeTracker.IsStakedCond.L.Unlock()

	// Signal that the IsStaked status is available
	node.NodeTracker.IsStakedCond.Signal()
	node.NodeTracker.HandleNodeData(nodeData)
	logrus.Info("handleStream -> Received data:", string(message))
}

func (node *OracleNode) IsPublisher() bool {
	// Node is a publisher if it has a non-empty signature
	return node.Signature != ""
}

func (node *OracleNode) Version() string {
	return Version
}
