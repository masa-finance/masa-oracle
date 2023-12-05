package masa

import (
	"encoding/json"
	"log"
	"math"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

func (node *OracleNode) ListenToNodeTracker() {
	for {
		select {
		case nodeData := <-node.NodeTracker.NodeDataChan:
			// Marshal the nodeData into JSON
			jsonData, err := json.Marshal(nodeData)
			if err != nil {
				log.Printf("Error marshaling node data: %v", err)
				continue
			}
			// Publish the JSON data on the node.topic
			err = node.PubSubManager.Publish(NodeTopic, jsonData)
			if err != nil {
				log.Printf("Error publishing node data: %v", err)
			}
			// If the nodeData represents a join event, call OnJoin in a separate goroutine
			if nodeData.Activity == pubsub2.ActivityJoined {
				go node.OnJoinEvent(nodeData.PeerId)
			}

		case <-node.Context.Done():
			return
		}
	}
}

func (node *OracleNode) HandleMessage(msg *pubsub.Message) {
	var nodeData pubsub2.NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		log.Printf("Failed to unmarshal node data: %v", err)
		return
	}
	// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
	node.NodeTracker.HandleNodeData(&nodeData)
}

type NodeDataPage struct {
	Data         []pubsub2.NodeData `json:"data"`
	PageNumber   int                `json:"pageNumber"`
	TotalPages   int                `json:"totalPages"`
	TotalRecords int                `json:"totalRecords"`
}

func (node *OracleNode) SendNodeDataPage(peerID peer.ID, pageNumber int) {
	stream, err := node.Host.NewStream(node.Context, peerID, NodeDataSyncProtocol)
	if err != nil {
		log.Printf("Failed to open stream to %s: %v", peerID, err)
		return
	}

	allNodeData := node.NodeTracker.GetAllNodeData()
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
		log.Printf("Failed to marshal NodeDataPage: %v", err)
		return
	}

	_, err = stream.Write(jsonData)
	if err != nil {
		log.Printf("Failed to send NodeDataPage to %s: %v", peerID, err)
	}
}

func (node *OracleNode) SendNodeData(peerID peer.ID) {
	allNodeData := node.NodeTracker.GetAllNodeData()
	totalRecords := len(allNodeData)
	totalPages := int(math.Ceil(float64(totalRecords) / float64(PageSize)))

	for pageNumber := 0; pageNumber < totalPages; pageNumber++ {
		node.SendNodeDataPage(peerID, pageNumber)
	}
}

func (node *OracleNode) ReceiveNodeData(stream network.Stream) {
	defer stream.Close()

	jsonData := make([]byte, 1024)
	n, err := stream.Read(jsonData)
	if err != nil {
		log.Printf("Failed to read NodeData from %s: %v", stream.Conn().RemotePeer(), err)
		return
	}

	var nodeData []pubsub2.NodeData
	if err := json.Unmarshal(jsonData[:n], &nodeData); err != nil {
		log.Printf("Failed to unmarshal NodeData: %v", err)
		return
	}

	for _, data := range nodeData {
		node.NodeTracker.HandleNodeData(&data)
	}
}

func (node *OracleNode) OnJoinEvent(peerID peer.ID) {
	// Send NodeData to the new node
	node.SendNodeData(peerID)
}
