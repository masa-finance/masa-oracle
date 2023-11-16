package masa

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
)

type NodeLite struct {
	Host       host.Host
	PrivKey    crypto.PrivKey
	Protocol   protocol.ID
	multiAddrs multiaddr.Multiaddr
	Context    context.Context
	PeerChan   chan peer.AddrInfo
}

func NewNodeLite(privKey crypto.PrivKey, ctx context.Context) (*NodeLite, error) {
	// Start with the default scaling limits.
	scalingLimits := rcmgr.DefaultLimits
	concreteLimits := scalingLimits.AutoScale()
	limiter := rcmgr.NewFixedLimiter(concreteLimits)

	resourceManager, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return nil, err
	}

	addrStr := []string{
		"/ip4/0.0.0.0/tcp/0",
	}
	if os.Getenv(PortNbr) != "" {
		addrStr = []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", os.Getenv(PortNbr)),
		}
	}

	host, err := libp2p.New(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.ListenAddrStrings(addrStr...),
		libp2p.Identity(privKey),
		libp2p.ResourceManager(resourceManager),
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
	return &NodeLite{
		Host:       host,
		PrivKey:    privKey,
		Protocol:   nodeLiteProtocol,
		multiAddrs: myNetwork.GetMultiAddressForHostQuiet(host),
		Context:    ctx,
		PeerChan:   make(chan peer.AddrInfo),
	}, nil
	return nil, nil
}

func (node *NodeLite) Start() (err error) {
	logrus.Infof("Starting node with ID: %s", node.multiAddrs.String())
	node.Host.SetStreamHandler(node.Protocol, node.handleStream)
	node.StartMDNSDiscovery("masa-chat")
	return nil
}

func (node *NodeLite) StartMDNSDiscovery(rendezvous string) {
	myNetwork.WithMDNS(node.Host, rendezvous, node.PeerChan)
	go func() {
		for {
			select {
			case peer := <-node.PeerChan: // will block until we discover a peer
				logrus.Info("Found peer:", peer, ", connecting")

				if err := node.Host.Connect(node.Context, peer); err != nil {
					logrus.Error("Connection failed:", err)
					continue
				}

				// open a stream, this stream will be handled by handleStream other end
				stream, err := node.Host.NewStream(node.Context, peer.ID, node.Protocol)

				if err != nil {
					logrus.Error("Stream open failed", err)
				} else {
					rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

					go node.writeData(rw)
					go node.readData(rw)
					logrus.Info("Connected to:", peer)
				}
			case <-node.Context.Done():
				return
			}
		}
	}()
}

func (node *NodeLite) handleStream(stream network.Stream) {
	logrus.Info("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go node.readData(rw)
	go node.writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func (node *NodeLite) readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			logrus.Error("Error reading from buffer:", err)
			return
		}
		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

func (node *NodeLite) writeData(rw *bufio.ReadWriter) {
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
		time.Sleep(time.Second * 10)
	}
}
