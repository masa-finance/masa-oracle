package main

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
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

const nodeLiteProtocol = "masa_node_lite_protocol/1.0.0"

type NodeLite struct {
	Host       host.Host
	PrivKey    crypto.PrivKey
	Protocol   protocol.ID
	multiAddrs multiaddr.Multiaddr
	ctx        context.Context
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
		ctx:        ctx,
	}, nil
	return nil, nil
}

func (node *NodeLite) Start() (err error) {
	node.Host.SetStreamHandler(node.Protocol, handleStream)
	node.StartMDNSDiscovery("masa-chat")
	return nil
}

func (node *NodeLite) StartMDNSDiscovery(rendezvous string) {
	peerChan := myNetwork.StartMDNS(node.Host, rendezvous)
	go func() {
		for {
			select {
			case peer := <-peerChan: // will block until we discover a peer
				logrus.Info("Found peer:", peer, ", connecting")

				if err := node.Host.Connect(node.ctx, peer); err != nil {
					logrus.Error("Connection failed:", err)
					continue
				}

				// open a stream, this stream will be handled by handleStream other end
				stream, err := node.Host.NewStream(node.ctx, peer.ID, node.Protocol)

				if err != nil {
					logrus.Error("Stream open failed", err)
				} else {
					rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

					go node.writeData(rw)
					go node.readData(rw)
					logrus.Info("Connected to:", peer)
				}
			case <-node.ctx.Done():
				return
			}
		}
	}()
}

func handleStream(stream network.Stream) {
	logrus.Info("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func (node *NodeLite) readData(rw *bufio.ReadWriter) {
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
		time.Sleep(time.Second * 5)
	}
}
