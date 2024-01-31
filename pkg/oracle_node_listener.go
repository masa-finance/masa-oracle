package masa

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"

	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

func (node *OracleNode) ListenToNodeTracker() {
	for {
		select {
		case nodeData := <-node.NodeTracker.NodeDataChan:
			time.Sleep(1 * time.Second)
			if node.IsStaked && node.NodeTracker.IsStaked(node.Host.ID().String()) {
				// Marshal the nodeData into JSON
				jsonData, err := json.Marshal(nodeData)
				if err != nil {
					logrus.Errorf("Error marshaling node data: %v", err)
					continue
				}
				// Publish the JSON data on the node.topic
				err = node.PubSubManager.Publish(TopicWithVersion(NodeGossipTopic), jsonData)
				if err != nil {
					logrus.Errorf("Error publishing node data: %v", err)
				}
				// If the nodeData represents a join event and
				// the node is a boot node or (we don't want boot nodes to wait)
				// the node start time is greater than 5 minutes ago,
				// call SendNodeData in a separate goroutine
				if nodeData.Activity == pubsub2.ActivityJoined &&
					(os.Getenv(Peers) == "" || time.Now().Sub(node.StartTime) > 5*time.Minute) {
					go node.SendNodeData(nodeData.PeerId)
				}
			}
		case <-node.Context.Done():
			return
		}
	}
}

func (node *OracleNode) HandleMessage(msg *pubsub.Message) {
	var nodeData pubsub2.NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		logrus.Errorf("Failed to unmarshal node data: %v", err)
		return
	}
	// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
	node.NodeTracker.HandleNodeData(nodeData)
}

type NodeDataPage struct {
	Data         []pubsub2.NodeData `json:"data"`
	PageNumber   int                `json:"pageNumber"`
	TotalPages   int                `json:"totalPages"`
	TotalRecords int                `json:"totalRecords"`
}

func (node *OracleNode) SendNodeDataPage(allNodeData []pubsub2.NodeData, stream network.Stream, pageNumber int) {
	logrus.Debugf("SendNodeDataPage --> %s: Page: %d", stream.Conn().RemotePeer(), pageNumber)
	totalRecords := len(allNodeData)
	totalPages := int(math.Ceil(float64(totalRecords) / PageSize))

	startIndex := pageNumber * PageSize
	endIndex := startIndex + PageSize
	if endIndex > totalRecords {
		endIndex = totalRecords
	}
	nodeDataPage := NodeDataPage{
		Data:         allNodeData[startIndex:endIndex],
		PageNumber:   pageNumber,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}

	jsonData, err := json.Marshal(nodeDataPage)
	if err != nil {
		logrus.Errorf("Failed to marshal NodeDataPage: %v", err)
		return
	}

	_, err = stream.Write(append(jsonData, '\n'))
	if err != nil {
		logrus.Errorf("Failed to send NodeDataPage: %v", err)
	}
}

func (node *OracleNode) SendNodeData(peerID peer.ID) {
	if peerID == node.Host.ID() {
		return
	}
	// Check if the node is staked before proceeding
	if !node.NodeTracker.IsStaked(peerID.String()) {
		logrus.Debugf("Node %s is not staked. Aborting SendNodeData.", peerID)
		return
	}

	recipientNodeData := node.NodeTracker.GetNodeData(peerID.String())
	var nodeData []pubsub2.NodeData
	if recipientNodeData == nil {
		nodeData = node.NodeTracker.GetAllNodeData()
	} else {
		// set the time to LastLeft minus 5 minutes
		sinceTime := recipientNodeData.LastLeft.Add(-5 * time.Minute)
		nodeData = node.NodeTracker.GetUpdatedNodes(sinceTime)
	}
	totalRecords := len(nodeData)
	totalPages := int(math.Ceil(float64(totalRecords) / float64(PageSize)))

	stream, err := node.Host.NewStream(node.Context, peerID, ProtocolWithVersion(NodeDataSyncProtocol))
	if err != nil {
		logrus.Errorf("Failed to open stream to %s: %v", peerID, err)
		return
	}
	defer func(stream network.Stream) {
		err := stream.Close()
		if err != nil {
			logrus.Errorf("Failed to close stream: %v", err)
		}
	}(stream) // Ensure the stream is closed after sending the data
	logrus.Infof("Sending %d node data records to %s", totalRecords, peerID)
	for pageNumber := 0; pageNumber < totalPages; pageNumber++ {
		node.SendNodeDataPage(nodeData, stream, pageNumber)
	}
}

func (node *OracleNode) ReceiveNodeData(stream network.Stream) {
	logrus.Info("ReceiveNodeData")

	scanner := bufio.NewScanner(stream)
	//scanner.Scan() stops when it hits a new line
	for scanner.Scan() {
		data := scanner.Bytes()
		var page NodeDataPage
		if err := json.Unmarshal(data, &page); err != nil {
			logrus.Errorf("Failed to unmarshal NodeData page: %v", err)
			logrus.Errorf("%s", string(data))
			continue
		}

		for _, nd := range page.Data {
			node.NodeTracker.HandleNodeData(nd)
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.Errorf("Failed to read stream: %v", err)
	}
}

func (node *OracleNode) GossipNodeData(stream network.Stream) {
	logrus.Info("GossipNodeData")
	remotePeerId, nodeData, err := node.handleStreamData(stream)
	if err != nil {
		logrus.Errorf("Failed to read stream: %v", err)
		return
	}
	// Only allow gossip about a node from other nodes
	if remotePeerId.String() != nodeData.PeerId.String() {
		node.NodeTracker.HandleNodeData(nodeData)
	}
	err = stream.Close()
	if err != nil {
		logrus.Errorf("Failed to close stream: %v", err)
	}
}

func (node *OracleNode) handleStreamData(stream network.Stream) (peer.ID, pubsub2.NodeData, error) {
	// Log the peer.ID of the remote peer
	remotePeerID := stream.Conn().RemotePeer()
	logrus.Infof("received stream from %s", remotePeerID)
	jsonData := make([]byte, 4096)

	var buffer bytes.Buffer
	// Loop until all data is read from the stream
	for {
		n, err := stream.Read(jsonData)
		if err != nil && err != io.EOF {
			//try to read the data from the buffer, if it serializes to NodeData, return it
			var nodeData pubsub2.NodeData
			if err2 := json.Unmarshal(buffer.Bytes(), &nodeData); err2 != nil {
				logrus.Errorf("Failed to read stream from %s: %v", remotePeerID, err)
				return "", pubsub2.NodeData{}, err
			}
			return remotePeerID, nodeData, nil
		}
		// when the other side closes the connection right away we get the EOF right away, so you have to write
		// to the buffer before checking for the EOF
		buffer.Write(jsonData[:n])
		if err == io.EOF {
			// All data has been read
			break
		}
	}
	var nodeData pubsub2.NodeData
	if err := json.Unmarshal(buffer.Bytes(), &nodeData); err != nil {
		logrus.Errorf("Failed to unmarshal NodeData: %v", err)
		logrus.Errorf("%s", buffer.String())
		return "", pubsub2.NodeData{}, err
	}
	if !nodeData.IsStaked {
		return "", pubsub2.NodeData{}, errors.New(fmt.Sprintf("un-staked node is ignored: %s", nodeData.PeerId))
	}
	return remotePeerID, nodeData, nil
}
