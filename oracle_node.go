package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

type OracleNode struct {
	Host       host.Host
	PrivKey    crypto.PrivKey
	DHT        *dht.IpfsDHT
	Protocol   protocol.ID
	multiAddrs multiaddr.Multiaddr
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
	addrStr := "/ip4/0.0.0.0/tcp/0/ws"
	if os.Getenv(portNbr) != "" {
		addrStr = fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", os.Getenv(portNbr))
	}
	host, err := libp2p.New(
		libp2p.Transport(websocket.New),
		libp2p.ListenAddrStrings(addrStr),
		libp2p.ResourceManager(rm),
		libp2p.Identity(privKey),
		libp2p.Ping(false), // disable built-in ping
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
	)
	if err != nil {
		return nil, err
	}
	nodeProtocol := protocol.ID(oracleProtocol)
	return &OracleNode{
		Host:     host,
		PrivKey:  privKey,
		Protocol: nodeProtocol,
		ctx:      ctx,
	}, nil
}

func (node *OracleNode) Start() error {
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
		node.Host.Close()
	}()
	go func() {
		router := gin.Default()

		// Use the auth middleware for the /webhook route
		router.POST("/webhook", node.authMiddleware(), node.webhookHandler)

		// Paths to the certificate and key files
		certFile := os.Getenv(cert)
		keyFile := os.Getenv(certPem)

		if err := router.RunTLS(":8080", certFile, keyFile); err != nil {
			logrus.Error("Failed to start HTTPS server:", err)
		}
	}()

	peersStr := os.Getenv(peers)
	bootstrapPeers := strings.Split(peersStr, ",")
	addrs := make([]multiaddr.Multiaddr, 0)
	for _, peerAddr := range bootstrapPeers {
		addr, err := multiaddr.NewMultiaddr(peerAddr)
		if err != nil {
			return err
		}
		addrs = append(addrs, addr)
	}

	node.DiscoverAndJoin(addrs)
	go node.sendMessageToRandomPeer()

	return nil
}

func (node *OracleNode) handleMessage(stream network.Stream) {
	logrus.Infof("handleMessage: %s", node.Host.ID().String())
	buf := bufio.NewReader(stream)
	message, err := buf.ReadString('\n')
	if err != nil {
		logrus.Error(err)
		stream.Reset()
	}
	connection := stream.Conn()

	logrus.Infof("Message from '%s': %s, remote: %s", connection.RemotePeer().String(), message, stream.Conn().RemotePeer())
	// Send an acknowledgement
	_, err = stream.Write([]byte("ACK\n"))
	if err != nil {
		logrus.Error("Error writing to stream:", err)
	}
}

func (node *OracleNode) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		token := c.GetHeader("Authorization")

		// Check the token
		if token != "your_expected_token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func (node *OracleNode) webhookHandler(c *gin.Context) {
	// Handle the webhook request here
	// ...

	c.JSON(http.StatusOK, gin.H{"message": "Webhook called"})
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
	return node.Host.ID().Pretty()
}

func (node *OracleNode) Addresses() string {
	addressesString := make([]string, 0)
	for _, address := range node.Host.Addrs() {
		addressesString = append(addressesString, address.String())
	}
	return strings.Join(addressesString, ", ")
}

func (node *OracleNode) DiscoverAndJoin(bootstrapPeers []multiaddr.Multiaddr) error {
	kademliaDHT, err := dht.New(node.ctx, node.Host)
	if err != nil {
		return err
	}
	if err = kademliaDHT.Bootstrap(node.ctx); err != nil {
		return err
	}
	node.DHT = kademliaDHT

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			logrus.Error(err)
		}
		if peerinfo.ID == node.Host.ID() {
			logrus.Info("Skipping connect to self")
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := node.Host.Connect(node.ctx, *peerinfo); err != nil {
				logrus.Warning(err)
			} else {
				logrus.Info("Connection established with bootstrap node:", *peerinfo)
				stream, err := node.Host.NewStream(node.ctx, peerinfo.ID, node.Protocol)
				if err != nil {
					logrus.Error("Error opening stream:", err)
				}
				_, err = stream.Write([]byte(fmt.Sprintf("Hello from %s\n", node.multiAddrs.String())))
				if err != nil {
					logrus.Error("Error writing to stream:", err)
				}
			}
		}()
	}
	wg.Wait()

	routingDiscovery := routing.NewRoutingDiscovery(kademliaDHT)
	logrus.Info("Announcing ourselves...")

	logrus.Debug("Searching for other peers...")
	// Use the routing discovery to find peers.
	peerChan, err := routingDiscovery.FindPeers(node.ctx, string(node.Protocol))
	if err != nil {
		return err
	}
	for peer := range peerChan {
		if peer.ID == node.Host.ID() {
			continue
		}
		logrus.Infof("Found peer: %s", peer.String())
		// Send a message with this node's multi address string to each peer that is found
		stream, err := node.Host.NewStream(node.ctx, peer.ID, node.Protocol)
		if err != nil {
			logrus.Error("Error opening stream:", err)
			continue
		}
		_, err = stream.Write([]byte(fmt.Sprintf("Hello from %s", node.multiAddrs.String())))
		if err != nil {
			logrus.Error("Error writing to stream:", err)
			continue
		}
	}
	logrus.Infof("found %d peers", len(node.Host.Network().Peers()))
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
					logrus.Error("Error opening stream:", err)
					continue
				}

				// Send a message to this peer
				_, err = stream.Write([]byte(fmt.Sprintf("Hello from %s\n", node.multiAddrs.String())))
				if err != nil {
					logrus.Error("Error writing to stream:", err)
				}
			}
		case <-node.ctx.Done():
			return
		}
	}
}
