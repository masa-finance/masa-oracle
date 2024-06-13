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

// NewNodeEventTracker creates a new NodeEventTracker instance.
// It initializes the node data map, node data channel, node data file path,
// connect buffer map. It loads existing node data from file, starts a goroutine
// to clear expired buffer entries, and returns the initialized instance.
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

// Listen is called when the node starts listening on a new multiaddr.
// It logs the network and address that listening started on.
func (net *NodeEventTracker) Listen(n network.Network, a ma.Multiaddr) {
	// This method is called when the node starts listening on a multiaddr
	logrus.WithFields(logrus.Fields{
		"network": n,
		"address": a,
	}).Info("Started listening")
}

// ListenClose logs when the node stops listening on a multiaddr.
// It logs the network and multiaddr that was stopped listening on.
func (net *NodeEventTracker) ListenClose(n network.Network, a ma.Multiaddr) {
	// This method is called when the node stops listening on a multiaddr
	logrus.WithFields(logrus.Fields{
		"network": n,
		"address": a,
	}).Info("Stopped listening")
}

// Connected handles when a remote peer connects to this node.
// It checks if the connecting node already exists in nodeData.
// If not, it returns without doing anything.
// If it exists but is not active, it buffers the connect event.
// If it exists and is active, it marks the node as joined and
// saves the updated nodeData.
func (net *NodeEventTracker) Connected(n network.Network, c network.Conn) {
	// A node has joined the network
	remotePeer := c.RemotePeer()
	peerID := remotePeer.String()

	nodeData, exists := net.nodeData.Get(peerID)
	if !exists {
		return
	} else {
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
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("Connected")
}

// Disconnected handles when a remote peer disconnects from this node.
// It looks up the disconnected node's data. If it doesn't exist, it returns.
// If it exists but is not staked, it returns. Otherwise it handles
// disconnect logic like updating the node data, deleting buffered connect
// events, and sending updated node data through the channel.
func (net *NodeEventTracker) Disconnected(n network.Network, c network.Conn) {
	logrus.Debug("Disconnect")

	peerID := c.RemotePeer().String()

	nodeData, exists := net.nodeData.Get(peerID)
	if !exists {
		// this should never happen
		logrus.Debugf("Node data does not exist for disconnected node: %s", peerID)
		return
	}
	buffered := net.ConnectBuffer[peerID]
	if buffered.NodeData != nil {
		buffered.NodeData.Left()
		delete(net.ConnectBuffer, peerID)
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

// HandleMessage unmarshals the received pubsub message into a NodeData struct,
// and passes it to HandleNodeData for further processing. This allows the
// NodeEventTracker to handle incoming node data messages from the pubsub layer.
func (net *NodeEventTracker) HandleMessage(msg *pubsub.Message) {
	var nodeData NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		logrus.Errorf("failed to unmarshal node data: %v", err)
		return
	}
	// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
	net.HandleNodeData(nodeData)
}

// RefreshFromBoot updates the node data map with the provided NodeData
// when the node boots up. It associates the NodeData with the peer ID string
// as the map key.
func (net *NodeEventTracker) RefreshFromBoot(data NodeData) {
	net.nodeData.Set(data.PeerId.String(), &data)
}

// HandleNodeData processes incoming NodeData from the pubsub layer.
// It adds new NodeData to the nodeData map, checks for replay attacks,
// and reconciles discrepancies between incoming and existing NodeData.
// This allows the tracker to maintain an up-to-date view of the node
// topology based on pubsub messages.
func (net *NodeEventTracker) HandleNodeData(data NodeData) {
	logrus.Debugf("Handling node data for: %s", data.PeerId)
	// we want nodeData for status even if staked is false
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
	// Check for replay attacks using LastUpdated
	if !data.LastUpdated.After(existingData.LastUpdated) {
		if existingData.IsStaked {
			logrus.Debugf("Stale or replayed node data received for node: %s", data.PeerId)
			return
		} else {
			// this is the boot node and local data is incorrect, take the value from the boot node
			net.nodeData.Set(data.PeerId.String(), &data)
			return
		}
	}
	existingData.LastUpdated = data.LastUpdated

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
		logrus.Debugf("Received unstaked status for node: %s", data.PeerId)
	}
	err := net.AddOrUpdateNodeData(existingData, true)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// GetNodeData returns the NodeData for the node with the given peer ID,
// or nil if no NodeData exists for that peer ID.
func (net *NodeEventTracker) GetNodeData(peerID string) *NodeData {
	nodeData, exists := net.nodeData.Get(peerID)
	if !exists {
		return nil
	}
	return nodeData
}

// GetAllNodeData returns a slice containing the NodeData for all nodes currently tracked.
func (net *NodeEventTracker) GetAllNodeData() []NodeData {
	logrus.Debug("Getting all node data")
	return net.nodeData.GetStakedNodesSlice()
}

// GetUpdatedNodes returns a slice of NodeData for nodes that have been updated since the given time.
// It filters the full node data set to only those updated after the passed in time,
// sorts the filtered results by update timestamp, and returns the sorted slice.
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

// DumpNodeData writes the NodeData map to a JSON file. It determines the file path
// based on the configured data directory, defaulting to nodeDataFile if not set.
// It logs any errors writing the file. This allows periodically persisting the
// node data.
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

// LoadNodeData loads the node data from a JSON file. It determines the file path
// based on the configured data directory, defaulting to nodeDataFile if not set.
// It logs any errors reading or parsing the file. This allows initializing the
// node data tracker from persisted data.
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

// GetEthAddress returns the Ethereum address for the given remote peer.
// It gets the peer's public key from the network's peerstore, converts
// it to a hex string, and converts that to an Ethereum address.
// Returns an empty string if there is no public key for the peer.
func GetEthAddress(remotePeer peer.ID, n network.Network) string {
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

// IsStaked returns whether the node with the given peerID is marked as staked in the node data tracker.
// Returns false if no node data is found for the given peerID.
func (net *NodeEventTracker) IsStaked(peerID string) bool {
	peerNd := net.GetNodeData(peerID)
	if peerNd == nil {
		return false
	}
	return peerNd.IsStaked
}

// AddOrUpdateNodeData adds or updates the node data in the node event tracker.
// If the node data does not exist, it is added and marked as self-identified.
// If the node data exists, it updates the staked status, Ethereum address, and multiaddress if needed.
// It also sends the updated node data to the NodeDataChan if the data changed or forceGossip is true.
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
		dataChanged = true
		nd.BytesScraped = nodeData.BytesScraped
		nd.IsStaked = nodeData.IsStaked
		nd.IsDiscordScraper = nodeData.IsDiscordScraper
		nd.IsTwitterScraper = nodeData.IsTwitterScraper
		nd.IsWebScraper = nodeData.IsWebScraper
		nd.Records = nodeData.Records
		nd.LastLeft = nodeData.LastLeft
		nd.Multiaddrs = nodeData.Multiaddrs
		nd.EthAddress = nodeData.EthAddress
		nd.IsActive = nodeData.IsActive

		logrus.WithFields(logrus.Fields{
			"Peer": nd.PeerId.String(),
		}).Info("Connected")
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

// ClearExpiredBufferEntries periodically clears expired entries from the
// connect buffer cache. It loops forever sleeping for a configured interval.
// On each loop it checks the current time against the connect time for each
// entry, and if expired, processes the connect and removes the entry.
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
