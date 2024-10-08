package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
)

type NodeEventTracker struct {
	NodeDataChan  chan *NodeData
	nodeData      *SafeMap
	nodeDataFile  string
	ConnectBuffer map[string]ConnectBufferEntry
	nodeVersion   string
}

type ConnectBufferEntry struct {
	NodeData    *NodeData
	ConnectTime time.Time
}

// NodeSorter provides methods for sorting NodeData slices
type NodeSorter struct {
	nodes []NodeData
	less  func(i, j NodeData) bool
}

// Len returns the length of the nodes slice
func (s NodeSorter) Len() int { return len(s.nodes) }

// Swap swaps the nodes at indices i and j
func (s NodeSorter) Swap(i, j int) { s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i] }

// Less compares nodes at indices i and j using the provided less function
func (s NodeSorter) Less(i, j int) bool { return s.less(s.nodes[i], s.nodes[j]) }

// SortNodesByTwitterReliability sorts the given nodes based on their Twitter reliability.
// It uses multiple criteria to determine the reliability and performance of nodes:
//  1. Prioritizes nodes that have been found more often (lower NotFoundCount)
//  2. Considers the last time a node was not found (earlier LastNotFoundTime is better)
//  3. Sorts by higher number of returned tweets
//  4. Then by more recent last returned tweet
//  5. Then by lower number of timeouts
//  6. Then by less recent last timeout
//  7. Finally, sorts by PeerId for stability when no performance data is available
//
// The function modifies the input slice in-place, sorting the nodes from most to least reliable.
func SortNodesByTwitterReliability(nodes []NodeData) {
	sorter := NodeSorter{
		nodes: nodes,
		less: func(i, j NodeData) bool {
			// First, prioritize nodes that have been found more often
			if i.NotFoundCount != j.NotFoundCount {
				return i.NotFoundCount < j.NotFoundCount
			}
			// Then, consider the last time they were not found
			if !i.LastNotFoundTime.Equal(j.LastNotFoundTime) {
				return i.LastNotFoundTime.Before(j.LastNotFoundTime)
			}
			// Primary sort: Higher number of returned tweets
			if i.ReturnedTweets != j.ReturnedTweets {
				return i.ReturnedTweets > j.ReturnedTweets
			}
			// Secondary sort: More recent last returned tweet
			if !i.LastReturnedTweet.Equal(j.LastReturnedTweet) {
				return i.LastReturnedTweet.After(j.LastReturnedTweet)
			}
			// Tertiary sort: Lower number of timeouts
			if i.TweetTimeouts != j.TweetTimeouts {
				return i.TweetTimeouts < j.TweetTimeouts
			}
			// Quaternary sort: Less recent last timeout
			if !i.LastTweetTimeout.Equal(j.LastTweetTimeout) {
				return i.LastTweetTimeout.Before(j.LastTweetTimeout)
			}
			// Default sort: By PeerId (ensures stable sorting when no performance data is available)
			return i.PeerId.String() < j.PeerId.String()
		},
	}
	sort.Sort(sorter)
}

// NewNodeEventTracker creates a new NodeEventTracker instance.
// It initializes the node data map, node data channel, node data file path,
// connect buffer map. It loads existing node data from file, starts a goroutine
// to clear expired buffer entries, and returns the initialized instance.
func NewNodeEventTracker(version, environment, hostId string) *NodeEventTracker {
	net := &NodeEventTracker{
		nodeData:      NewSafeMap(),
		nodeVersion:   version,
		NodeDataChan:  make(chan *NodeData),
		nodeDataFile:  fmt.Sprintf("%s_%s_node_data.json", version, environment),
		ConnectBuffer: make(map[string]ConnectBufferEntry),
	}
	go net.ClearExpiredBufferEntries()
	go net.StartCleanupRoutine(context.Background(), hostId)
	return net
}

// Listen is called when the node starts listening on a new multiaddr.
// It logs the network and address that listening started on.
func (net *NodeEventTracker) Listen(n network.Network, a ma.Multiaddr) {
	// This method is called when the node starts listening on a multiaddr
	logrus.WithFields(logrus.Fields{
		"network": n,
		"address": a,
	}).Info("[+] Started listening")
}

// ListenClose logs when the node stops listening on a multiaddr.
// It logs the network and multiaddr that was stopped listening on.
func (net *NodeEventTracker) ListenClose(n network.Network, a ma.Multiaddr) {
	// This method is called when the node stops listening on a multiaddr
	logrus.WithFields(logrus.Fields{
		"network": n,
		"address": a,
	}).Info("[-]Stopped listening")
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
			nodeData.Joined(net.nodeVersion)
			err := net.AddOrUpdateNodeData(nodeData, true)
			if err != nil {
				logrus.Error("[-] Error adding or updating node data: ", err)
				return
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("[+] Connected")
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
		buffered.NodeData.Joined(net.nodeVersion)
		net.NodeDataChan <- buffered.NodeData
	} else {
		nodeData.Left()
		net.NodeDataChan <- nodeData
	}
	logrus.WithFields(logrus.Fields{
		"Peer":    c.RemotePeer().String(),
		"network": n,
		"conn":    c,
	}).Info("[+] Disconnected")
}

// HandleMessage unmarshals the received pubsub message into a NodeData struct,
// and passes it to HandleNodeData for further processing. This allows the
// NodeEventTracker to handle incoming node data messages from the pubsub layer.
func (net *NodeEventTracker) HandleMessage(msg *pubsub.Message) {
	var nodeData NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		logrus.Errorf("[-] Failed to unmarshal node data: %v", err)
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
		if data.LastLeftUnix > data.LastJoinedUnix {
			return
		}
		// Otherwise, add it
		logrus.Debugf("Adding new node data: %s", data.PeerId.String())
		net.nodeData.Set(data.PeerId.String(), &data)
		return
	}
	// Check for replay attacks using LastUpdated
	if data.LastUpdatedUnix < existingData.LastUpdatedUnix {
		if existingData.IsStaked {
			logrus.Debugf("Stale or replayed node data received for node: %s", data.PeerId)
			return
		} else {
			// this is the boot node and local data is incorrect, take the value from the boot node
			net.nodeData.Set(data.PeerId.String(), &data)
			return
		}
	}
	existingData.LastUpdatedUnix = data.LastUpdatedUnix

	maxDifference := time.Millisecond * 15

	// Handle discrepancies for existing nodes
	if data.LastJoinedUnix > 0 &&
		data.LastJoinedUnix < existingData.LastJoinedUnix &&
		data.LastJoinedUnix > existingData.LastLeftUnix &&
		time.Since(time.Unix(data.LastJoinedUnix, 0)) < maxDifference {
		existingData.LastJoinedUnix = data.LastJoinedUnix
	}
	if data.LastLeftUnix > 0 &&
		data.LastLeftUnix > existingData.LastLeftUnix &&
		data.LastLeftUnix < existingData.LastJoinedUnix &&
		time.Since(time.Unix(data.LastLeftUnix, 0)) < maxDifference {
		existingData.LastLeftUnix = data.LastLeftUnix
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
		logrus.Error("[-] Error adding or updating node data: ", err)
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
		if nodeData.LastUpdatedUnix > since.Unix() {
			updatedNodeData = append(updatedNodeData, nodeData)
		}
	}
	// Sort the slice based on the timestamp
	sort.Slice(updatedNodeData, func(i, j int) bool {
		return updatedNodeData[i].LastUpdatedUnix < updatedNodeData[j].LastUpdatedUnix
	})
	return updatedNodeData
}

// GetEthAddress returns the Ethereum address for the given remote peer.
// It gets the peer's public key from the network's peer store, converts
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

// GetEligibleWorkerNodes returns a slice of NodeData for nodes that are eligible to perform a specific category of work.
func (net *NodeEventTracker) GetEligibleWorkerNodes(category WorkerCategory) []NodeData {
	logrus.Debugf("Getting eligible worker nodes for category: %s", category)
	result := make([]NodeData, 0)
	for _, nodeData := range net.GetAllNodeData() {
		if nodeData.CanDoWork(category) {
			result = append(result, nodeData)
		}
	}

	// Sort the eligible nodes based on the worker category
	switch category {
	case CategoryTwitter:
		SortNodesByTwitterReliability(result)
		// Add cases for other categories as needed such as
		// web
		// discord
		// telegram
	}

	return result
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
	logrus.Debugf("Handling node data for: %s", nodeData.PeerId)
	dataChanged := false

	nd, ok := net.nodeData.Get(nodeData.PeerId.String())
	if !ok {
		nodeData.SelfIdentified = true
		nodeData.Joined(net.nodeVersion)
		net.NodeDataChan <- nodeData
		net.nodeData.Set(nodeData.PeerId.String(), nodeData)
	} else {
		if !nd.SelfIdentified {
			dataChanged = true
			nd.SelfIdentified = true
		}
		dataChanged = true
		nd.IsStaked = nodeData.IsStaked
		nd.IsDiscordScraper = nodeData.IsDiscordScraper
		nd.IsTelegramScraper = nodeData.IsTelegramScraper
		nd.IsTwitterScraper = nodeData.IsTwitterScraper
		nd.IsWebScraper = nodeData.IsWebScraper
		nd.Records = nodeData.Records
		nd.Multiaddrs = nodeData.Multiaddrs
		nd.EthAddress = nodeData.EthAddress
		if nd.EthAddress == "" && nodeData.EthAddress != "" {
			dataChanged = true
			nd.EthAddress = nodeData.EthAddress
		}

		if len(nodeData.Multiaddrs) > 0 {
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

			nd.LastUpdatedUnix = nodeData.LastUpdatedUnix
			net.nodeData.Set(nodeData.PeerId.String(), nd)

		}

		// If the node data exists, check if the multiaddress is already in the list
		// multiAddress := nodeData.Multiaddrs[0].Multiaddr
		// addrExists := false
		// for _, addr := range nodeData.Multiaddrs {
		// 	if addr.Equal(multiAddress) {
		// 		addrExists = true
		// 		break
		// 	}
		// }
		// if !addrExists {
		// 	nodeData.Multiaddrs = append(nodeData.Multiaddrs, JSONMultiaddr{multiAddress})
		// }
		// if dataChanged || forceGossip {
		// 	net.NodeDataChan <- nd
		// }

		// nd.LastUpdatedUnix = nodeData.LastUpdatedUnix
		// net.nodeData.Set(nodeData.PeerId.String(), nd)
	}
	return nil
}

// ClearExpiredBufferEntries periodically clears expired entries from the
// connect buffer cache. It loops forever sleeping for a configured interval.
// On each loop it checks the current time against the connect time for each
// entry, and if expired, processes the connect and removes the entry.
func (net *NodeEventTracker) ClearExpiredBufferEntries() {
	for {
		time.Sleep(1 * time.Minute)
		now := time.Now()
		for peerID, entry := range net.ConnectBuffer {
			if now.Sub(entry.ConnectTime) > time.Minute*1 {
				// first force a leave event so that timestamps are updated properly
				entry.NodeData.Left()
				// Buffer period expired without a disconnect, process connect
				entry.NodeData.Joined(net.nodeVersion)
				net.NodeDataChan <- entry.NodeData
				delete(net.ConnectBuffer, peerID)
			}
		}
	}
}

// RemoveNodeData removes the node data associated with the given peer ID from the NodeEventTracker.
// It deletes the node data from the internal map and removes any corresponding entry
// from the connect buffer. This function is typically called when a peer disconnects
// or is no longer part of the network.
//
// Parameters:
//   - peerID: A string representing the ID of the peer to be removed.
//
// TODO: we should never remove node data from the internal map. Otherwise we lose all tracking of activity.
//func (net *NodeEventTracker) RemoveNodeData(peerID string) {
//	net.nodeData.Delete(peerID)
//	delete(net.ConnectBuffer, peerID)
//	logrus.Infof("[+] Removed peer %s from NodeTracker", peerID)
//}

// ClearExpiredWorkerTimeouts periodically checks and clears expired worker timeouts.
// It runs in an infinite loop, sleeping for 5 minutes between each iteration.
// For each node in the network, it checks if the worker timeout has expired (after 60 minutes).
// If a timeout has expired, it resets the WorkerTimeout to zero and updates the node data.
// This function helps manage the availability of workers in the network by clearing
// temporary timeout states.
func (net *NodeEventTracker) ClearExpiredWorkerTimeouts() {
	for {
		time.Sleep(5 * time.Minute) // Check every 5 minutes
		now := time.Now()

		for _, nodeData := range net.GetAllNodeData() {
			if !nodeData.WorkerTimeout.IsZero() && now.Sub(nodeData.WorkerTimeout) >= 16*time.Minute {
				nodeData.WorkerTimeout = time.Time{} // Reset to zero value
				err := net.AddOrUpdateNodeData(&nodeData, true)
				if err != nil {
					logrus.Warnf("Error adding worker timeout %v", err)
				}
			}
		}
	}
}

const (
	maxDisconnectionTime = 1 * time.Minute
	cleanupInterval      = 2 * time.Minute
)

// StartCleanupRoutine starts a goroutine that periodically checks for and removes stale peers
func (net *NodeEventTracker) StartCleanupRoutine(ctx context.Context, hostId string) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			net.cleanupStalePeers(hostId)
		case <-ctx.Done():
			return
		}
	}
}

// cleanupStalePeers checks for and removes stale peers from both the routing table and node data
func (net *NodeEventTracker) cleanupStalePeers(hostId string) {
	now := time.Now()

	for _, nodeData := range net.GetAllNodeData() {
		if now.Sub(time.Unix(nodeData.LastUpdatedUnix, 0)) > maxDisconnectionTime {
			if nodeData.PeerId.String() != hostId {
				logrus.Infof("Removing stale peer: %s", nodeData.PeerId)
				delete(net.ConnectBuffer, nodeData.PeerId.String())

				// Notify about peer removal
				net.NodeDataChan <- &NodeData{
					PeerId:          nodeData.PeerId,
					Activity:        ActivityLeft,
					LastUpdatedUnix: now.Unix(),
				}
			}

			// Use the node parameter to access OracleNode methods if needed
			// For example:
			// node.SomeMethod(nodeData.PeerId)
		}
	}
}

func (net *NodeEventTracker) UpdateNodeDataTwitter(peerID string, updates NodeData) error {
	nodeData, exists := net.nodeData.Get(peerID)
	if !exists {
		return fmt.Errorf("node data not found for peer ID: %s", peerID)
	}

	// Update fields based on non-zero values
	nodeData.UpdateTwitterFields(updates)

	// Save the updated node data
	err := net.AddOrUpdateNodeData(nodeData, true)
	if err != nil {
		return fmt.Errorf("error updating node data: %v", err)
	}
	return nil
}
