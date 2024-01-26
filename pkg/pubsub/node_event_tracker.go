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
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/crypto"
)

type NodeEventTracker struct {
	NodeDataChan chan *NodeData
	nodeData     map[string]*NodeData
	dataMutex    sync.RWMutex
	version      string
}

func NewNodeEventTracker(version string) *NodeEventTracker {
	net := &NodeEventTracker{
		nodeData:     make(map[string]*NodeData),
		NodeDataChan: make(chan *NodeData),
		version:      version,
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
	remotePeer := c.RemotePeer()
	peerID := remotePeer.String()
	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()

	nodeData, exists := net.nodeData[peerID]
	if !exists {
		return
	} else {
		if nodeData.IsStaked {
			err := net.AddOrUpdateNodeData(nodeData)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Connected")
	nodeData.Joined()
}

func (net *NodeEventTracker) Disconnected(n network.Network, c network.Conn) {
	logrus.Debug("Disconnect")

	pubKeyHex := getEthAddress(c.RemotePeer(), n)
	peerID := c.RemotePeer().String()

	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()
	nodeData, exists := net.nodeData[peerID]
	if !nodeData.IsStaked {
		return
	}
	if !exists {
		// this should never happen
		logrus.Warnf("Node data does not exist for disconnected node: %s", peerID)
		nodeData = NewNodeData(c.RemoteMultiaddr(), c.RemotePeer(), pubKeyHex, ActivityLeft)
	}
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Disconnected")
	net.NodeDataChan <- nodeData
	nodeData.Left()
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
	if !data.IsStaked {
		return
	}
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
		if existingData.IsStaked {
			logrus.Warnf("Stale or replayed node data received for node: %s", data.PeerId)
			return
		} else {
			//this is the boot node and local data is incorrect, take the value from the boot node
			net.nodeData[data.PeerId.String()] = &data
			return
		}
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
	net.dataMutex.RLock()
	defer net.dataMutex.RUnlock()

	// Convert the map to a slice
	nodeDataSlice := make([]NodeData, 0, len(net.nodeData))
	for _, nodeData := range net.nodeData {
		if !nodeData.IsStaked {
			continue
		}
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
	// Filter allNodeData to only include nodes that have been updated since the given time
	var updatedNodeData []NodeData
	for _, nodeData := range net.GetAllNodeData() {
		if nodeData.LastUpdated.After(since) {
			updatedNodeData = append(updatedNodeData, nodeData)
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

func (net *NodeEventTracker) IsStaked(peerID string) bool {
	peerNd := net.GetNodeData(peerID)
	if peerNd == nil {
		return false
	}
	return peerNd.IsStaked
}

func (net *NodeEventTracker) AddOrUpdateNodeData(nodeData *NodeData) error {
	logrus.Debug("Adding self identity")
	net.dataMutex.Lock()
	defer net.dataMutex.Unlock()
	dataChanged := false

	nd, exists := net.nodeData[nodeData.PeerId.String()]
	if !exists {
		nodeData.SelfIdentified = true
		net.nodeData[nodeData.PeerId.String()] = nodeData
		nodeData.Joined()
		net.NodeDataChan <- nodeData
	} else {
		if !nd.SelfIdentified {
			dataChanged = true
			nd.SelfIdentified = true
		}
		if !nd.IsStaked && nodeData.IsStaked {
			dataChanged = true
			nd.IsStaked = nodeData.IsStaked
			logrus.WithFields(logrus.Fields{
				"Peer": nd.PeerId.String(),
			}).Info("Connected")
		}
		if nd.EthAddress == "" && nodeData.EthAddress != "" {
			dataChanged = true
			nd.EthAddress = nodeData.EthAddress
		}
		// If the node data exists, check if the multiaddress is already in the list
		multiAddress := nodeData.Multiaddrs[0].Multiaddr
		addrExists := false
		for _, addr := range nodeData.Multiaddrs {
			if addr.Equal(multiAddress) {
				addrExists = true
				break
			}
		}
		if !addrExists {
			nodeData.Multiaddrs = append(nodeData.Multiaddrs, JSONMultiaddr{multiAddress})
		}
		if dataChanged {
			net.NodeDataChan <- nd
		}
	}
	return nil
}
