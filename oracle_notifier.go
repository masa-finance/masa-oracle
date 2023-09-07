package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/sirupsen/logrus"
)

type OracleNotifier struct {
	Host  host.Host
	Nodes []peer.ID
}

func NewOracleNotifier(privKey crypto.PrivKey, ctx context.Context) (*OracleNotifier, error) {
	host, err := libp2p.New(
		libp2p.Transport(websocket.New),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0/ws"),
		libp2p.Identity(privKey),
		libp2p.Ping(false),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
	)
	if err != nil {
		return nil, err
	}

	return &OracleNotifier{
		Host:  host,
		Nodes: make([]peer.ID, 0),
	}, nil
}

func (notifier *OracleNotifier) Start() {
	notifier.Host.SetStreamHandler("/announce", notifier.handleAnnounce)

	router := gin.Default()

	// Use the auth middleware for the /webhook route
	router.POST("/webhook", notifier.authMiddleware(), notifier.handleWebhook)

	// Paths to the certificate and key files
	certFile := os.Getenv(cert)
	keyFile := os.Getenv(certPem)

	if err := router.RunTLS(":8080", certFile, keyFile); err != nil {
		logrus.Error("Failed to start HTTPS server:", err)
	}
}

func (notifier *OracleNotifier) handleAnnounce(stream network.Stream) {
	nodeID := stream.Conn().RemotePeer()
	notifier.Nodes = append(notifier.Nodes, nodeID)
}

func (notifier *OracleNotifier) authMiddleware() gin.HandlerFunc {
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

func (notifier *OracleNotifier) handleWebhook(c *gin.Context) {
	if len(notifier.Nodes) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No nodes available"})
		return
	}

	// Choose a random node
	randNode := notifier.Nodes[rand.Intn(len(notifier.Nodes))]

	// Create a new stream with this node
	stream, err := notifier.Host.NewStream(context.Background(), randNode, "/message")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to node"})
		return
	}

	// Send a message to this node
	_, err = stream.Write([]byte("Webhook called\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook called"})
}
