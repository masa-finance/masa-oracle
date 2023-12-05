package pubsub

import (
	"encoding/json"
	"sort"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

type NodeEventTracker struct {
	NodeDataChan chan *NodeData
	nodeData     map[string]*NodeData
	dataMutex    sync.RWMutex
	changes      int
}

func NewNodeEventTracker() *NodeEventTracker {
	return &NodeEventTracker{
		nodeData:     make(map[string]*NodeData),
		NodeDataChan: make(chan *NodeData),
	}
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
	net.NodeDataChan <- NewNodeData(c.RemoteMultiaddr(), c.RemotePeer(), ActivityJoined)

	peerID := c.RemotePeer().String()
	net.dataMutex.Lock()
	nodeData, exists := net.nodeData[peerID]
	if !exists {
		nodeData = NewNodeData(c.RemoteMultiaddr(), c.RemotePeer(), ActivityJoined)
		net.nodeData[peerID] = nodeData
	} else {
		// If the node data exists, check if the multiaddress is already in the list
		addrExists := false
		for _, addr := range nodeData.Multiaddrs {
			if addr.Equal(c.RemoteMultiaddr()) {
				addrExists = true
				break
			}
		}
		if !addrExists {
			nodeData.Multiaddrs = append(nodeData.Multiaddrs, c.RemoteMultiaddr())
		}
	}
	nodeData.Joined()
	net.dataMutex.Unlock()
}

func (net *NodeEventTracker) Disconnected(n network.Network, c network.Conn) {
	// A node has left the network
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Disconnected")
	net.NodeDataChan <- NewNodeData(c.RemoteMultiaddr(), c.RemotePeer(), ActivityJoined)

	peerID := c.RemotePeer().String()
	net.dataMutex.Lock()
	nodeData, exists := net.nodeData[peerID]
	if exists {
		nodeData.Left()
	}
	net.dataMutex.Unlock()
}

func (net *NodeEventTracker) HandleMessage(msg *pubsub.Message) {
	var nodeData NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		logrus.Errorf("failed to unmarshal node data: %v", err)
	}
	// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
	net.HandleNodeData(&nodeData)
}

func (net *NodeEventTracker) HandleNodeData(data *NodeData) {
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
		net.nodeData[data.PeerId.String()] = data
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
	existingData.AccumulatedUptime = existingData.GetAccumulatedUptime()
}

func (net *NodeEventTracker) GetAllNodeData() []NodeData {
	logrus.Debug("Getting all node data")
	// Convert the map to a slice
	nodeDataSlice := make([]NodeData, 0, len(net.nodeData))
	for _, nodeData := range net.nodeData {
		nodeDataSlice = append(nodeDataSlice, *nodeData)
	}

	// Sort the slice based on the timestamp
	sort.Slice(nodeDataSlice, func(i, j int) bool {
		return nodeDataSlice[i].LastUpdated.Before(nodeDataSlice[j].LastUpdated)
	})
	return nodeDataSlice
}
