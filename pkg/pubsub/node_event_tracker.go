package pubsub

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/crypto"
)

type NodeEventTracker struct {
	NodeDataChan chan *NodeData
	nodeData     map[string]*NodeData
	dataMutex    sync.RWMutex
	changes      int
}

func NewNodeEventTracker() *NodeEventTracker {
	net := &NodeEventTracker{
		nodeData:     make(map[string]*NodeData),
		NodeDataChan: make(chan *NodeData),
	}
	err := net.LoadNodeData()
	if err != nil {
		logrus.Error("Error loading node data", err)
	}
	return net
}

func (net *NodeEventTracker) Listen(n network.Network, a ma.Multiaddr) {
	// This method is called when the node starts listening on a multiaddr
	logrus.WithFields(logrus.Fields{
		"network": n,
		"address": a,
	}).Info("Started listening")
}

func (net *NodeEventTracker) ListenClose(n network.Network, a ma.Multiaddr) {
	// This method is called when the node stops listening on a multiaddr
	logrus.WithFields(logrus.Fields{
		"network": n,
		"address": a,
	}).Info("Stopped listening")
}

func (net *NodeEventTracker) Connected(n network.Network, c network.Conn) {
	// A node has joined the network
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Connected")

	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()

	ethAddress := getEthAddress(c.RemotePeer(), n)
	peerID := c.RemotePeer().String()

	nodeData, exists := net.nodeData[peerID]
	if !exists {
		nodeData = NewNodeData(c.RemoteMultiaddr(), c.RemotePeer(), ethAddress, ActivityJoined)
		net.nodeData[nodeData.PeerId.String()] = nodeData
	} else {
		if nodeData.EthAddress == "" {
			nodeData.EthAddress = ethAddress
		}
		// If the node data exists, check if the multiaddress is already in the list
		addrExists := false
		for _, addr := range nodeData.Multiaddrs {
			if addr.Equal(c.RemoteMultiaddr()) {
				addrExists = true
				break
			}
		}
		if !addrExists {
			nodeData.Multiaddrs = append(nodeData.Multiaddrs, JSONMultiaddr{c.RemoteMultiaddr()})
		}
	}
	net.NodeDataChan <- nodeData
	nodeData.Joined()
}

func (net *NodeEventTracker) Disconnected(n network.Network, c network.Conn) {
	// A node has left the network
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Disconnected")

	pubKeyHex := getEthAddress(c.RemotePeer(), n)
	peerID := c.RemotePeer().String()

	net.dataMutex.Lock()
	nodeData, exists := net.nodeData[peerID]
	if !exists {
		// this should never happen
		logrus.Warnf("Node data does not exist for disconnected node: %s", peerID)
		nodeData = NewNodeData(c.RemoteMultiaddr(), c.RemotePeer(), pubKeyHex, ActivityLeft)
	}
	net.NodeDataChan <- nodeData
	nodeData.Left()

	net.dataMutex.Unlock()
}

func (net *NodeEventTracker) HandleMessage(msg *pubsub.Message) {
	var nodeData NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		logrus.Errorf("failed to unmarshal node data: %v", err)
	}
	// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
	net.HandleNodeData(nodeData)
}

func (net *NodeEventTracker) HandleNodeData(data NodeData) {
	logrus.Debugf("Handling node data for: %s", data.PeerId)
	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()

	existingData, ok := net.nodeData[data.PeerId.String()]
	if !ok {
		// If the node data does not exist in the cache and the node has left, ignore it
		if data.LastLeft.After(data.LastJoined) {
			return
		}
		// Otherwise, add it
		logrus.Debugf("Adding new node data: %s", data.PeerId.String())
		net.nodeData[data.PeerId.String()] = &data
		return
	}
	// Handle discrepancies for existing nodes
	if data.LastJoined.Before(existingData.LastJoined) && data.LastJoined.After(existingData.LastLeft) {
		existingData.LastJoined = data.LastJoined
	}
	if data.LastLeft.After(existingData.LastLeft) && data.LastLeft.Before(existingData.LastJoined) {
		existingData.LastLeft = data.LastLeft
	}
	// Update accumulated uptime
	//existingData.AccumulatedUptime = existingData.GetAccumulatedUptime()
}

func (net *NodeEventTracker) GetAllNodeData() []NodeData {
	logrus.Debug("Getting all node data")
	// Convert the map to a slice
	nodeDataSlice := make([]NodeData, 0, len(net.nodeData))
	for _, nodeData := range net.nodeData {
		nd := *nodeData
		nd.CurrentUptime = nodeData.GetCurrentUptime()
		nd.AccumulatedUptime = nodeData.GetAccumulatedUptime()
		nd.CurrentUptimeStr = prettyDuration(nd.CurrentUptime)
		nd.AccumulatedUptimeStr = prettyDuration(nd.AccumulatedUptime)
		nodeDataSlice = append(nodeDataSlice, nd)
	}

	// Sort the slice based on the timestamp
	sort.Slice(nodeDataSlice, func(i, j int) bool {
		return nodeDataSlice[i].LastUpdated.Before(nodeDataSlice[j].LastUpdated)
	})
	return nodeDataSlice
}

func (net *NodeEventTracker) DumpNodeData() {
	// Lock the nodeData map for concurrent read
	net.dataMutex.RLock()
	defer net.dataMutex.RUnlock()

	// Convert the nodeData map to JSON
	data, err := json.Marshal(net.nodeData)
	if err != nil {
		// Handle error
	}

	// Write the JSON data to a file
	filePath := os.Getenv("nodeBackupPath")
	if filePath == "" {
		filePath = "node_data.json"
	}
	logrus.Infof("writing node data to file: %s", filePath)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		logrus.Error(fmt.Sprintf("could not write to file: %s", filePath), err)
	}
}

func (net *NodeEventTracker) LoadNodeData() error {
	// Read the JSON data from a file
	filePath := os.Getenv("nodeBackupPath")
	if filePath == "" {
		filePath = "node_data.json"
	}
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logrus.Warn(fmt.Sprintf("file does not exist: %s", filePath))
		return nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Error(fmt.Sprintf("could not read from file: %s", filePath), err)
		return err
	}
	// Convert the JSON data to a map
	nodeData := make(map[string]*NodeData)
	err = json.Unmarshal(data, &nodeData)
	if err != nil {
		logrus.Error("could not unmarshal JSON data", err)
		return err
	}
	// Lock the nodeData map for concurrent write
	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()
	// remove invalids from an earlier bug
	for key, value := range nodeData {
		if key != value.PeerId.String() {
			logrus.Warnf("peer ID mismatch: %s != %s", key, value.PeerId)
			delete(nodeData, key)
		}
	}
	// Replace the nodeData map with the new map
	logrus.Info("Loaded node data from file")
	net.nodeData = nodeData
	return nil
}

func prettyDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	min := int64(d / time.Minute)
	h := min / 60
	min %= 60
	days := h / 24
	h %= 24

	if days > 0 {
		return fmt.Sprintf("%d days %d hours %d minutes", days, h, min)
	}
	if h > 0 {
		return fmt.Sprintf("%d hours %d minutes", h, min)
	}
	return fmt.Sprintf("%d minutes", min)
}

func getEthAddress(remotePeer peer.ID, n network.Network) string {
	var publicKeyHex string
	var err error

	// Get the public key of the remote peer
	pubKey := n.Peerstore().PubKey(remotePeer)
	if pubKey == nil {
		logrus.WithFields(logrus.Fields{
			"Peer": remotePeer.String(),
		}).Warn("No public key found for peer")
	} else {
		publicKeyHex, err = crypto.Libp2pPubKeyToEthAddress(pubKey)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Peer": remotePeer.String(),
			}).Warnf("Error getting public key %v", err)
		}
	}
	return publicKeyHex
}
