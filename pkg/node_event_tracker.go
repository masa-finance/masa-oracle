package masa

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/mudler/edgevpn/pkg/blockchain"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

type NodeEventTracker struct {
	ledger    *blockchain.Ledger
	nodeData  map[string]*NodeData
	dataMutex sync.RWMutex
	changes   int
}

func NewNodeEventTracker(ledger *blockchain.Ledger) *NodeEventTracker {
	return &NodeEventTracker{
		ledger:   ledger,
		nodeData: make(map[string]*NodeData),
	}
}

func (net *NodeEventTracker) Listen(n network.Network, a ma.Multiaddr) {
	// This method is called when the node starts listening on a multiaddr
	log.WithFields(log.Fields{
		"network": n,
		"address": a,
	}).Info("Started listening")
}

func (net *NodeEventTracker) ListenClose(n network.Network, a ma.Multiaddr) {
	// This method is called when the node stops listening on a multiaddr
	log.WithFields(log.Fields{
		"network": n,
		"address": a,
	}).Info("Stopped listening")
}

func (net *NodeEventTracker) Connected(n network.Network, c network.Conn) {
	// A node has joined the network
	log.WithFields(log.Fields{
		"network": n,
		"conn":    c,
	}).Info("Connected")

	peerID := c.RemotePeer().String()
	net.dataMutex.Lock()
	nodeData, exists := net.nodeData[peerID]
	if !exists {
		multiaddr := c.RemoteMultiaddr()
		nodeData = NewNodeData(multiaddr)
		net.nodeData[peerID] = nodeData
	}
	nodeData.Joined()
	net.dataMutex.Unlock()
}

func (net *NodeEventTracker) Disconnected(n network.Network, c network.Conn) {
	// A node has left the network
	log.WithFields(log.Fields{
		"network": n,
		"conn":    c,
	}).Info("Disconnected")

	peerID := c.RemotePeer().String()
	net.dataMutex.Lock()
	nodeData, exists := net.nodeData[peerID]
	if exists {
		nodeData.Left()
	}
	net.dataMutex.Unlock()
}

func (net *NodeEventTracker) WriteToLedger() {
	net.dataMutex.RLock()
	// Get the timestamp of the last block in the ledger
	lastBlockTime, _ := time.Parse(time.RFC3339, net.ledger.LastBlock().Timestamp)
	for peerID, nodeData := range net.nodeData {
		// Check if the NodeData has been updated since the last block was added to the ledger
		if nodeData.LastUpdated.After(lastBlockTime) {
			// Convert NodeData to JSON
			data, _ := json.Marshal(nodeData)
			net.ledger.Add(peerID, map[string]interface{}{
				"nodeData": string(data),
			})
		}
	}
	net.dataMutex.RUnlock()
}
