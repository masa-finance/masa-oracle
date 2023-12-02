package masa

import (
	"encoding/json"
	"log"
	"math"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
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
			err = node.PubSubManager.Publish(masaNodeTopic, jsonData)
			if err != nil {
				log.Printf("Error publishing node data: %v", err)
			}
		case <-node.Context.Done():
			return
		}
	}
}

func (node *OracleNode) HandleMessage(msg *pubsub.Message) {
	var nodeData NodeData
	if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
		log.Printf("Failed to unmarshal node data: %v", err)
		return
	}
	// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
	node.NodeTracker.HandleIncomingData(&nodeData)
}

type NodeDataPage struct {
	Data         []NodeData `json:"data"`
	PageNumber   int        `json:"pageNumber"`
	TotalPages   int        `json:"totalPages"`
	TotalRecords int        `json:"totalRecords"`
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
	stream, err := node.Host.NewStream(node.Context, peerID, NodeDataSyncProtocol)
	if err != nil {
		log.Printf("Failed to open stream to %s: %v", peerID, err)
		return
	}

	nodeData := node.NodeTracker.GetAllNodeData()
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		log.Printf("Failed to marshal NodeData: %v", err)
		return
	}

	_, err = stream.Write(jsonData)
	if err != nil {
		log.Printf("Failed to send NodeData to %s: %v", peerID, err)
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

	var nodeData []NodeData
	if err := json.Unmarshal(jsonData[:n], &nodeData); err != nil {
		log.Printf("Failed to unmarshal NodeData: %v", err)
		return
	}

	for _, data := range nodeData {
		node.NodeTracker.HandleIncomingData(&data)
	}
}

func (node *OracleNode) OnJoinEvent(peerID peer.ID) {
	// Existing join event handling code...

	// Send NodeData to the new node
	node.SendNodeData(peerID)
}
