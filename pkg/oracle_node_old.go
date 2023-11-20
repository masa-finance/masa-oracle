package masa

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/mudler/edgevpn/pkg/blockchain"
	"github.com/mudler/edgevpn/pkg/hub"

	"github.com/libp2p/go-libp2p/p2p/host/autonat"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
)

type OracleNodeOld struct {
	Host       host.Host
	PrivKey    crypto.PrivKey
	DHT        *dht.IpfsDHT
	Protocol   protocol.ID
	multiAddrs multiaddr.Multiaddr
	topic      *pubsub.Topic
	ctx        context.Context
	sync.Mutex
	ledger    *blockchain.Ledger
	tempCh    chan *hub.Message
	inputCh   chan *NodeData
	nodeData  map[string]*NodeData
	dataMutex sync.RWMutex
	changes   int
	PeerChan  chan myNetwork.PeerEvent
}

func NewOracleNodeOld(privKey crypto.PrivKey, ctx context.Context) (*OracleNodeOld, error) {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	rm, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}
	addrStr := []string{
		//"/ip4/0.0.0.0/udp/0/quic-v1",
		"/ip4/0.0.0.0/tcp/0",
	}
	if os.Getenv(PortNbr) != "" {
		addrStr = []string{
			//fmt.Sprintf("/ip4/0.0.0.0/udp/%s/quic-v1", os.Getenv(PortNbr)),
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", os.Getenv(PortNbr)),
		}
	}

	newHost, err := libp2p.New(
		//libp2p.Transport(quic.NewTransport),
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
	return &OracleNodeOld{
		Host:     newHost,
		PrivKey:  privKey,
		Protocol: nodeProtocol,
		ctx:      ctx,
		inputCh:  make(chan *NodeData),
		tempCh:   make(chan *hub.Message),
		nodeData: make(map[string]*NodeData),
		PeerChan: make(chan myNetwork.PeerEvent),
	}, nil
}

func (node *OracleNodeOld) Start() (err error) {
	err = node.Gossip(oracleProtocol)
	if err != nil {
		return err
	}
	go node.handleEvents()

	mw, err := node.messageWriter()
	if err != nil {
		return err
	}
	// Create a new ledger
	node.ledger = blockchain.New(mw, &blockchain.MemoryStore{})

	node.Host.SetStreamHandler(node.Protocol, node.handleMessage)

	// important to order:  make sure the ledger is already initialized
	node.Host.Network().Notify(NewNodeEventTracker(node.inputCh))

	node.multiAddrs, err = myNetwork.GetMultiAddressForHost(node.Host)
	if err != nil {
		return err
	}
	fmt.Println("libp2p host address:", node.multiAddrs.String())

	bootNodeAddrs, err := myNetwork.GetBootNodesMultiAddress(os.Getenv(Peers))
	if err != nil {
		return err
	}
	err = node.DiscoverAndJoin(bootNodeAddrs)
	if err != nil {
		return err
	}
	go node.sendMessageToRandomPeer()
	return
}

func (node *OracleNodeOld) Gossip(topicName string) error {
	gossipSub, err := pubsub.NewGossipSub(node.ctx, node.Host)
	if err != nil {
		return err
	}

	// Subscribe to a topic
	node.topic, err = gossipSub.Join(topicName)
	if err != nil {
		return err
	}
	go myNetwork.StreamConsoleTo(node.ctx, node.topic)

	sub, err := node.topic.Subscribe()
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(node.ctx)
			if err != nil {
				logrus.Errorf("sub.Next: %s", err.Error())
				if err.Error() == "context canceled" {
					return
				}
				continue
			}
			// Skip messages from the same node
			if msg.ReceivedFrom == node.Host.ID() {
				//logrus.Debugf("message received from: %s", msg.ReceivedFrom.String())
				//logrus.Debug(string(msg.Message.Data))
				continue
			}
			var addrs multiaddr.Multiaddr
			connectedness := node.Host.Network().Connectedness(msg.ReceivedFrom)
			if connectedness == network.Connected {
				peerInfo := node.Host.Peerstore().PeerInfo(msg.ReceivedFrom)
				if len(peerInfo.Addrs) == 0 {
					continue
				}
				addrString := fmt.Sprintf("%s/p2p/%s", peerInfo.Addrs[0].String(), peerInfo.ID.String())
				logrus.Infof("%s : %s : %s", msg.ReceivedFrom, string(msg.Message.Data), addrString)
				addrs, err = multiaddr.NewMultiaddr(addrString)
				if err != nil {
					logrus.Error("Failed to create multiaddress:", err)
					continue
				}
			} else {
				logrus.Info(msg.ReceivedFrom, ": ", string(msg.Message.Data))
			}

			// Check if the peer is already in the DHT and Peerstore
			peerInfo2, err := peer.AddrInfoFromP2pAddr(addrs)
			if err != nil {
				logrus.Errorf("Failed to get AddrInfo from multiaddress: %s, %v", addrs.String(), err)
				continue
			}

			if node.Host.Peerstore().PeerInfo(peerInfo2.ID).ID == "" {
				// The peer is not in the Peerstore, add it
				node.Host.Peerstore().AddAddrs(peerInfo2.ID, peerInfo2.Addrs, peerstore.PermanentAddrTTL)
			}

			if node.DHT.RoutingTable().Find(peerInfo2.ID) != "" {
				// The peer is not in the DHT, add it
				node.DHT.RoutingTable().TryAddPeer(peerInfo2.ID, false, false)
			}
		}
	}()
	return nil
}

func (node *OracleNodeOld) handleEvents() {
	for {
		select {
		case data := <-node.inputCh:
			if data == nil {
				continue
			}
			_, exists := node.nodeData[data.PeerId.String()]
			switch data.Activity {
			case ActivityLeft:
				if exists {
					data.Left()
					node.changes++
				}
			case ActivityJoined:
				if !exists {
					node.nodeData[data.PeerId.String()] = data
					data.Joined()
					node.changes++
					jsnBytes, err := json.Marshal(data)
					if err != nil {
						logrus.Errorf("Failed to marshal NodeData: %v", err)
						continue
					}
					multiAddr, err := multiaddr.NewMultiaddr(data.Address())
					if err != nil {
						logrus.Error("Failed to create multiaddress:", err)
						continue
					}
					peerInfo, err := peer.AddrInfoFromP2pAddr(multiAddr)
					if err != nil {
						logrus.Errorf("Failed to get AddrInfo from multiaddress: %s, %v", multiAddr.String(), err)
					}

					if node.Host.Peerstore().PeerInfo(peerInfo.ID).ID == "" {
						// The peer is not in the Peerstore, add it
						node.Host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)
					}

					logrus.Debugf("Publishing: %s", string(jsnBytes))
					err = node.topic.Publish(node.ctx, jsnBytes)
					if err != nil {
						logrus.Error("Error publishing to topic:", err)
					}
				}
			default:
				logrus.Errorf("Unknown activity: %d", data.Activity)
				continue
			}
			// if node.changes > 10 {
			//	node.WriteToLedger()
			// }
			node.WriteToLedger()

		case <-node.ctx.Done():
			<-node.ctx.Done()
			err := node.Host.Close()
			if err != nil {
				return
			}
		}
	}
}

func (node *OracleNodeOld) handleMessage(stream network.Stream) {
	logrus.Infof("handleMessage: %s", node.Host.ID().String())
	buf := bufio.NewReader(stream)
	message, err := buf.ReadString('\n')
	if err != nil {
		logrus.Errorf("handleMessage: %s", err.Error())
	}
	connection := stream.Conn()

	logrus.Infof("Message from '%s': %s, remote: %s", connection.RemotePeer().String(), message, connection.RemoteMultiaddr())
	// peerinfo, err := peer.AddrInfoFromP2pAddr(connection.RemoteMultiaddr())
	// Send an acknowledgement
	_, err = stream.Write([]byte("ACK\n"))
	if err != nil {
		if errors.Is(err, network.ErrReset) {
			logrus.Info("Stream was reset, skipping write operation")
		} else {
			logrus.Error("Error writing to stream:", err)
		}
	}
}

// Connect is useful for testing in ai local environment, It should probably be removed
func (node *OracleNodeOld) Connect(targetNode *OracleNodeOld) error {
	targetNodeAddressInfo := host.InfoFromHost(targetNode.Host)

	err := node.Host.Connect(context.Background(), *targetNodeAddressInfo)
	if err != nil {
		return err
	}
	logrus.Infof("connected: to %s", targetNode.Id())
	return nil
}

func (node *OracleNodeOld) Id() string {
	return node.Host.ID().String()
}

func (node *OracleNodeOld) Addresses() string {
	addressesString := make([]string, 0)
	for _, address := range node.Host.Addrs() {
		addressesString = append(addressesString, address.String())
	}
	return strings.Join(addressesString, ", ")
}

func (node *OracleNodeOld) DiscoverAndJoin(bootstrapPeers []multiaddr.Multiaddr) error {
	var err error
	node.DHT, err = myNetwork.WithDht(node.ctx, node.Host, bootstrapPeers, node.Protocol, node.PeerChan)
	if err != nil {
		return err
	}
	go myNetwork.Discover(node.ctx, node.Host, node.DHT, node.Protocol, node.multiAddrs)
	err = node.SetupAutoNAT()
	if err != nil {
		return err
	}
	return nil
}

func (node *OracleNodeOld) SetupAutoNAT() error {
	privKey, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 2048)
	if err != nil {
		return err
	}
	newHost, err := libp2p.New(
		//libp2p.Transport(quic.NewTransport),
		//libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/0/quic-v1"),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", os.Getenv(PortNbr))),
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

func (node *OracleNodeOld) sendMessageToRandomPeer() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			peers := node.Host.Network().Peers()
			if len(peers) > 0 {
				logrus.Info("******************************************************************************")
				for _, peer := range peers {
					logrus.Info("Peer ID: ", peer)
				}
				logrus.Info("******************************************************************************")
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
					if err == network.ErrReset {
						logrus.Info("Stream was reset, skipping write operation")
					} else {
						logrus.Error("Error writing to stream:", err)
					}
					continue
				}
				// publish a message on the Topic
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

// messageWriter returns a new MessageWriter bound to the edgevpn instance
// with the given options
func (node *OracleNodeOld) messageWriter(opts ...hub.MessageOption) (*messageWriter, error) {
	mess := &hub.Message{}
	mess.Apply(opts...)

	return &messageWriter{
		// input: node.inputCh,
		mess: mess,
	}, nil
}

func (node *OracleNodeOld) WriteToLedger() {
	logrus.Debug("WriteToLedger")
	node.dataMutex.RLock()
	// Get the timestamp of the last block in the ledger
	lastBlockTime, _ := time.Parse(time.RFC3339, node.ledger.LastBlock().Timestamp)
	for peerID, nodeData := range node.nodeData {
		// Check if the NodeData has been updated since the last block was added to the ledger
		if lastBlockTime.IsZero() || nodeData.LastUpdated.After(lastBlockTime) {
			// Convert NodeData to JSON
			data, _ := json.Marshal(nodeData)
			node.ledger.Add(peerID, map[string]interface{}{
				"nodeData": string(data),
			})
		}
	}
	node.changes = 0
	node.dataMutex.RUnlock()
}

func (node *OracleNodeOld) readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			logrus.Error("Error reading from buffer:", err)
			return
		}

		if str != "" && str != "\n" {
			logrus.Infof("MDNS Received message: %s from %s", str, node.multiAddrs.String())
		}
	}
}

func (node *OracleNodeOld) writeData(rw *bufio.ReadWriter) {
	for {
		// Generate a message including the multiaddress of the sender
		sendData := fmt.Sprintf("MDNS Hello from %s\n", node.multiAddrs.String())

		_, err := rw.WriteString(sendData)
		if err != nil {
			logrus.Error("Error writing to buffer:", err)
			return
		}
		err = rw.Flush()
		if err != nil {
			logrus.Error("Error flushing buffer:", err)
			return
		}
		// Sleep for a while before sending the next message
		time.Sleep(time.Second * 5)
	}
}
