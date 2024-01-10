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
	NodeDataChan   chan *NodeData
	nodeData       map[string]*NodeData
	dataMutex      sync.RWMutex
	IsStakedStatus map[string]bool
	IsStakedCond   *sync.Cond
	version        string
}

func NewNodeEventTracker(version string) *NodeEventTracker {
	net := &NodeEventTracker{
		nodeData:       make(map[string]*NodeData),
		NodeDataChan:   make(chan *NodeData),
		IsStakedCond:   sync.NewCond(&sync.Mutex{}),
		IsStakedStatus: make(map[string]bool),
		version:        version,
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

	peerID := c.RemotePeer().String()
	isStaked, ok := net.waitForStakedStatus(peerID, time.Second*10)
	if !ok {
		logrus.Errorf("Timeout waiting for IsStaked status for node: %s", peerID)
		return
	}
	if !isStaked {
		logrus.Infof("Ignoring unstaked node: %s", peerID)
		return
	}

	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()

	ethAddress := getEthAddress(c.RemotePeer(), n)

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
		return
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
	// Check for replay attacks using LastUpdated
	if !data.LastUpdated.After(existingData.LastUpdated) {
		logrus.Warnf("Stale or replayed node data received for node: %s", data.PeerId)
		return
	}
	existingData.LastUpdated = data.LastUpdated

	maxDifference := time.Minute * 5

	// Handle discrepancies for existing nodes
	if !data.LastJoined.IsZero() &&
		data.LastJoined.Before(existingData.LastJoined) &&
		data.LastJoined.After(existingData.LastLeft) &&
		time.Since(data.LastJoined) < maxDifference {
		existingData.LastJoined = data.LastJoined
	}
	if !data.LastLeft.IsZero() &&
		data.LastLeft.After(existingData.LastLeft) &&
		data.LastLeft.Before(existingData.LastJoined) &&
		time.Since(data.LastLeft) < maxDifference {
		existingData.LastLeft = data.LastLeft
	}

	if existingData.EthAddress == "" && data.EthAddress != "" {
		existingData.EthAddress = data.EthAddress
	}
	if data.IsStaked && !existingData.IsStaked {
		existingData.IsStaked = data.IsStaked
	} else if !data.IsStaked && existingData.IsStaked {
		logrus.Warnf("Received unstaked status for node: %s", data.PeerId)
	}
}

func (net *NodeEventTracker) GetNodeData(peerID string) *NodeData {
	net.dataMutex.RLock()
	defer net.dataMutex.RUnlock()

	nodeData, exists := net.nodeData[peerID]
	if !exists {
		return nil
	}
	return nodeData
}

func (net *NodeEventTracker) GetAllNodeData() []NodeData {
	logrus.Debug("Getting all node data")
	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()

	// Convert the map to a slice
	nodeDataSlice := make([]NodeData, 0, len(net.nodeData))
	for _, nodeData := range net.nodeData {
		nd := *nodeData
		nd.CurrentUptime = nodeData.GetCurrentUptime()
		nd.AccumulatedUptime = nodeData.GetAccumulatedUptime()
		nd.CurrentUptimeStr = PrettyDuration(nd.CurrentUptime)
		nd.AccumulatedUptimeStr = PrettyDuration(nd.AccumulatedUptime)
		nodeDataSlice = append(nodeDataSlice, nd)
	}

	// Sort the slice based on the timestamp
	sort.Slice(nodeDataSlice, func(i, j int) bool {
		return nodeDataSlice[i].LastUpdated.Before(nodeDataSlice[j].LastUpdated)
	})
	return nodeDataSlice
}

func (net *NodeEventTracker) GetUpdatedNodes(since time.Time) []NodeData {
	net.dataMutex.RLock()
	defer net.dataMutex.RUnlock()

	// Filter allNodeData to only include nodes that have been updated since the given time
	var updatedNodeData []NodeData
	for _, nodeData := range net.nodeData {
		if nodeData.LastUpdated.After(since) {
			nd := *nodeData
			nd.CurrentUptime = nodeData.GetCurrentUptime()
			nd.AccumulatedUptime = nodeData.GetAccumulatedUptime()
			nd.CurrentUptimeStr = PrettyDuration(nd.CurrentUptime)
			nd.AccumulatedUptimeStr = PrettyDuration(nd.AccumulatedUptime)
			updatedNodeData = append(updatedNodeData, nd)
		}
	}
	// Sort the slice based on the timestamp
	sort.Slice(updatedNodeData, func(i, j int) bool {
		return updatedNodeData[i].LastUpdated.Before(updatedNodeData[j].LastUpdated)
	})
	return updatedNodeData
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
		filePath = fmt.Sprintf("%s_node_data.json", net.version)
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
		filePath = fmt.Sprintf("%s_node_data.json", net.version)
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

func PrettyDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	minute := int64(d / time.Minute)
	h := minute / 60
	minute %= 60
	days := h / 24
	h %= 24

	if days > 0 {
		return fmt.Sprintf("%d days %d hours %d minutes", days, h, minute)
	}
	if h > 0 {
		return fmt.Sprintf("%d hours %d minutes", h, minute)
	}
	return fmt.Sprintf("%d minutes", minute)
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

func (net *NodeEventTracker) waitForStakedStatus(peerID string, timeout time.Duration) (bool, bool) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		net.IsStakedCond.L.Lock()
		isStaked, ok := net.IsStakedStatus[peerID]
		net.IsStakedCond.L.Unlock()

		if ok {
			return isStaked, true
		}

		select {
		case <-time.After(time.Second):
			// Check the status again after a second
		case <-timer.C:
			return false, false // Timeout
		}
	}
}

func (net *NodeEventTracker) IsStaked(peerID string) bool {
	peerNd := net.GetNodeData(peerID)
	if peerNd == nil {
		return false
	}
	return peerNd.IsStaked
}
