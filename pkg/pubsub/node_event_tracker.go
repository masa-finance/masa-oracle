package pubsub

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
)

type NodeEventTracker struct {
	NodeDataChan  chan *NodeData
	nodeData      *SafeMap
	nodeDataFile  string
	ConnectBuffer map[string]ConnectBufferEntry
}

type ConnectBufferEntry struct {
	NodeData    *NodeData
	ConnectTime time.Time
}

func NewNodeEventTracker(version, environment string) *NodeEventTracker {
	net := &NodeEventTracker{
		nodeData:      NewSafeMap(),
		NodeDataChan:  make(chan *NodeData),
		nodeDataFile:  fmt.Sprintf("%s_%s_node_data.json", version, environment),
		ConnectBuffer: make(map[string]ConnectBufferEntry),
	}
	err := net.LoadNodeData()
	if err != nil {
		logrus.Error("Error loading node data", err)
	}
	go net.ClearExpiredBufferEntries()
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

	nodeData, exists := net.nodeData.Get(peerID)
	if !exists {
		return
	} else {
		if nodeData.IsStaked {
			if nodeData.IsActive {
				// Node appears already connected, buffer this connect event
				net.ConnectBuffer[peerID] = ConnectBufferEntry{NodeData: nodeData, ConnectTime: time.Now()}
			} else {
				nodeData.Joined()
				err := net.AddOrUpdateNodeData(nodeData, true)
				if err != nil {
					logrus.Error(err)
					return
				}
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Connected")
}

func (net *NodeEventTracker) Disconnected(n network.Network, c network.Conn) {
	logrus.Debug("Disconnect")

	peerID := c.RemotePeer().String()

	nodeData, exists := net.nodeData.Get(peerID)
	if !exists {
		// this should never happen
		logrus.Warnf("Node data does not exist for disconnected node: %s", peerID)
		return
	} else if !nodeData.IsStaked {
		return
	}
	buffered := net.ConnectBuffer[peerID]
	if buffered.NodeData != nil {
		buffered.NodeData.Left()
		delete(net.ConnectBuffer, peerID)
		// net.nodeData.Delete(peerID)
		// Now process the buffered connect
		buffered.NodeData.Joined()
		net.NodeDataChan <- buffered.NodeData
	} else {
		nodeData.Left()
		net.NodeDataChan <- nodeData
	}
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Disconnected")
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

func (net *NodeEventTracker) RefreshFromBoot(data NodeData) {
	net.nodeData.Set(data.PeerId.String(), &data)
}

func (net *NodeEventTracker) HandleNodeData(data NodeData) {
	logrus.Debugf("Handling node data for: %s", data.PeerId)
	if !data.IsStaked {
		return
	}
	existingData, ok := net.nodeData.Get(data.PeerId.String())
	if !ok {
		// If the node data does not exist in the cache and the node has left, ignore it
		if data.LastLeft.After(data.LastJoined) {
			return
		}
		// Otherwise, add it
		logrus.Debugf("Adding new node data: %s", data.PeerId.String())
		net.nodeData.Set(data.PeerId.String(), &data)
		return
	}
	// Check for replay attacks using LastUpdated -- @TODO check why is this considered a replay?
	if !data.LastUpdated.After(existingData.LastUpdated) {
		if existingData.IsStaked {
			logrus.Warnf("Stale or replayed node data received for node: %s", data.PeerId)
			return
		} else {
			// this is the boot node and local data is incorrect, take the value from the boot node
			net.nodeData.Set(data.PeerId.String(), &data)
			return
		}
	}
	existingData.LastUpdated = data.LastUpdated

	// maxDifference := time.Minute * 5
	maxDifference := time.Millisecond * 15

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
	err := net.AddOrUpdateNodeData(existingData, true)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func (net *NodeEventTracker) GetNodeData(peerID string) *NodeData {
	nodeData, exists := net.nodeData.Get(peerID)
	if !exists {
		return nil
	}
	return nodeData
}

func (net *NodeEventTracker) GetAllNodeData() []NodeData {
	logrus.Debug("Getting all node data")
	return net.nodeData.GetStakedNodesSlice()
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
	// Write the JSON data to a file
	var filePath string
	dataDir := config.GetInstance().MasaDir
	if dataDir == "" {
		filePath = net.nodeDataFile
	} else {
		filePath = filepath.Join(dataDir, net.nodeDataFile)
	}
	logrus.Infof("writing node data to file: %s", filePath)
	err := net.nodeData.DumpNodeData(filePath)
	if err != nil {
		logrus.Error("could not dump node data", err)
	}
}

func (net *NodeEventTracker) LoadNodeData() error {
	// Read the JSON data from a file
	var filePath string
	dataDir := config.GetInstance().MasaDir
	if dataDir == "" {
		filePath = net.nodeDataFile
	} else {
		filePath = filepath.Join(dataDir, net.nodeDataFile)
	}
	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logrus.Warn(fmt.Sprintf("file does not exist: %s", filePath))
		return nil
	}
	err := net.nodeData.LoadNodeData(filePath)
	if err != nil {
		logrus.Error("could not load node data", err)
		return err
	}
	return nil
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
		publicKeyHex, err = masacrypto.Libp2pPubKeyToEthAddress(pubKey)
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

func (net *NodeEventTracker) AddOrUpdateNodeData(nodeData *NodeData, forceGossip bool) error {
	logrus.Debug("Adding self identity")
	dataChanged := false

	nd, exists := net.nodeData.Get(nodeData.PeerId.String())
	if !exists {
		nodeData.SelfIdentified = true
		net.nodeData.Set(nodeData.PeerId.String(), nodeData)
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
		if dataChanged || forceGossip {
			net.NodeDataChan <- nd
		}
	}
	return nil
}

func (net *NodeEventTracker) ClearExpiredBufferEntries() {
	for {
		time.Sleep(30 * time.Second) // E.g., every 5 seconds
		now := time.Now()
		for peerID, entry := range net.ConnectBuffer {
			if now.Sub(entry.ConnectTime) > time.Minute*1 {
				// Buffer period expired without a disconnect, process connect
				entry.NodeData.Joined()
				net.NodeDataChan <- entry.NodeData
				delete(net.ConnectBuffer, peerID)
				// net.nodeData.Delete(peerID)
			}
		}
	}
}
