package masa

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/security/noise"

	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
)

type OracleNode struct {
	Host       host.Host
	PrivKey    crypto.PrivKey
	DHT        *dht.IpfsDHT
	Protocol   protocol.ID
	multiAddrs multiaddr.Multiaddr
	topic      *pubsub.Topic
	ctx        context.Context
}

func NewOracleNode(privKey crypto.PrivKey, ctx context.Context) (*OracleNode, error) {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	rm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}
	addrStr := []string{
		"/ip4/0.0.0.0/udp/0/quic-v1",
		"/ip4/0.0.0.0/tcp/0",
	}
	// There should not be a case where UDP port is not set but TCP port is set
	if os.Getenv(UdpPortNbr) != "" {
		if os.Getenv(TcpPortNbr) != "" {
			addrStr = []string{
				fmt.Sprintf("/ip4/0.0.0.0/udp/%s/quic-v1", os.Getenv(UdpPortNbr)),
				fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", os.Getenv(TcpPortNbr)),
			}
		} else {
			addrStr = []string{
				fmt.Sprintf("/ip4/0.0.0.0/udp/%s/quic-v1", os.Getenv(UdpPortNbr)),
			}
		}
	}

	newHost, err := libp2p.New(
		libp2p.Transport(quic.NewTransport),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.ListenAddrStrings(addrStr...),
		libp2p.ResourceManager(rm),
		libp2p.Identity(privKey),
		libp2p.Ping(false), // disable built-in ping
		libp2p.ChainOptions(
			libp2p.Security(noise.ID, noise.New),
			libp2p.Security(libp2ptls.ID, libp2ptls.New),
		),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(), // Enable Circuit Relay v2 with hop
	)
	if err != nil {
		return nil, err
	}
	nodeProtocol := protocol.ID(oracleProtocol)
	topic, err := myNetwork.NewPubSub(ctx, newHost, oracleProtocol)
	if err != nil {
		return nil, err
	}
	return &OracleNode{
		Host:     newHost,
		PrivKey:  privKey,
		Protocol: nodeProtocol,
		ctx:      ctx,
		topic:    topic,
	}, nil
}

func (node *OracleNode) Start() (err error) {
	node.Host.SetStreamHandler(node.Protocol, node.handleMessage)
	node.Host.Network().Notify(&ConnectionLogger{})

	peerInfo := peer.AddrInfo{
		ID:    node.Host.ID(),
		Addrs: node.Host.Addrs(),
	}
	multiaddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		return err
	}
	node.multiAddrs = multiaddrs[0]
	fmt.Println("libp2p host address:", multiaddrs[0])

	go func() {
		<-node.ctx.Done()
		err = node.Host.Close()
		if err != nil {
			return
		}
	}()

	peersStr := os.Getenv(Peers)
	bootstrapPeers := strings.Split(peersStr, ",")
	addrs := make([]multiaddr.Multiaddr, 0)
	for _, peerAddr := range bootstrapPeers {
		if peerAddr == "" {
			continue
		}
		addr, err := multiaddr.NewMultiaddr(peerAddr)
		if err != nil {
			return err
		}
		addrs = append(addrs, addr)
	}

	err = node.DiscoverAndJoin(addrs)
	if err != nil {
		return err
	}

	go node.sendMessageToRandomPeer()
	return
}

func (node *OracleNode) handleMessage(stream network.Stream) {
	logrus.Infof("handleMessage: %s", node.Host.ID().String())
	buf := bufio.NewReader(stream)
	message, err := buf.ReadString('\n')
	if err != nil {
		logrus.Errorf("handleMessage: %s", err.Error())
	}
	connection := stream.Conn()

	logrus.Infof("Message from '%s': %s, remote: %s", connection.RemotePeer().String(), message, connection.RemoteMultiaddr())
	//peerinfo, err := peer.AddrInfoFromP2pAddr(connection.RemoteMultiaddr())
	// Send an acknowledgement
	_, err = stream.Write([]byte("ACK\n"))
	if err != nil {
		if err == network.ErrReset {
			logrus.Info("Stream was reset, skipping write operation")
		} else {
			logrus.Error("Error writing to stream:", err)
		}
	}
}

// Connect is useful for testing in a local environment, It should probably be removed
func (node *OracleNode) Connect(targetNode *OracleNode) error {
	targetNodeAddressInfo := host.InfoFromHost(targetNode.Host)

	err := node.Host.Connect(context.Background(), *targetNodeAddressInfo)
	if err != nil {
		return err
	}
	logrus.Infof("connected: to %s", targetNode.Id())
	return nil
}

func (node *OracleNode) Id() string {
	return node.Host.ID().String()
}

func (node *OracleNode) Addresses() string {
	addressesString := make([]string, 0)
	for _, address := range node.Host.Addrs() {
		addressesString = append(addressesString, address.String())
	}
	return strings.Join(addressesString, ", ")
}

func (node *OracleNode) DiscoverAndJoin(bootstrapPeers []multiaddr.Multiaddr) error {
	var err error
	node.DHT, err = myNetwork.NewDht(node.ctx, node.Host, bootstrapPeers, node.Protocol, node.multiAddrs)
	if err != nil {
		return err
	}
	go myNetwork.Discover(node.ctx, node.Host, node.DHT, node.Protocol, node.multiAddrs)
	node.SetupAutoNAT()
	return nil
}

func (node *OracleNode) SetupAutoNAT() error {
	privKey, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 2048)
	if err != nil {
		return err
	}
	newHost, err := libp2p.New(
		libp2p.Transport(quic.NewTransport),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/0/quic-v1"),
		libp2p.Identity(privKey),
		libp2p.Ping(false), // disable built-in ping
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
	)
	if err != nil {
		return err
	}
	opts := []autonat.Option{
		autonat.EnableService(newHost.Network()),
		autonat.WithReachability(network.ReachabilityUnknown),
		autonat.WithoutThrottling(),
	}
	autoNat, err := autonat.New(node.Host, opts...)
	if err != nil {
		logrus.Fatal(err)
	}
	// Wait a bit for the service to bootstrap
	time.Sleep(5 * time.Second)

	// Get the public address
	reachability := autoNat.Status()
	logrus.Info("Reachability status:", reachability.String())
	return nil
}

func (node *OracleNode) sendMessageToRandomPeer() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			peers := node.Host.Network().Peers()
			if len(peers) > 0 {
				// Choose a random peer
				randPeer := peers[rand.Intn(len(peers))]

				// Create a new stream with this peer
				stream, err := node.Host.NewStream(node.ctx, randPeer, node.Protocol)
				if err != nil {
					// Check if the error is about the protocol
					if strings.Contains(err.Error(), "failed to negotiate protocol") {
						logrus.Info("Skipping peer due to protocol negotiation error:", err)
						continue
					}
					logrus.Error("Error opening stream:", err)
					continue
				}

				// Send a message to this peer
				_, err = stream.Write([]byte(fmt.Sprintf("ticker Hello from %s\n", node.multiAddrs.String())))
				if err != nil {
					if errors.Is(err, network.ErrReset) {
						logrus.Info("Stream was reset, skipping write operation")
					} else {
						logrus.Error("Error writing to stream:", err)
					}
					continue
				}
				//publish a message on the Topic
				err = node.topic.Publish(node.ctx, []byte(fmt.Sprintf("topic Hello from %s\n", node.multiAddrs.String())))
				if err != nil {
					logrus.Error("Error publishing to topic:", err)
				}
			}
		case <-node.ctx.Done():
			return
		}
	}
}
