package masa

import (
	"encoding/json"
	"log"

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
			err = node.topic.Publish(node.Context, jsonData)
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

func (node *OracleNode) SubscribeToTopic() error {
	sub, err := node.topic.Subscribe()
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(node.Context)
			if err != nil {
				log.Printf("Error reading from topic: %v", err)
				continue
			}

			var nodeData NodeData
			if err := json.Unmarshal(msg.Data, &nodeData); err != nil {
				log.Printf("Failed to unmarshal node data: %v", err)
				continue
			}
			// Handle the nodeData by calling NodeEventTracker.HandleIncomingData
			node.NodeTracker.HandleIncomingData(&nodeData)
		}
	}()
	return nil
}

func (node *OracleNode) SendNodeData(peerID peer.ID) {
	stream, err := node.Host.NewStream(node.Context, peerID, NodeDataSyncProtocol)
	if err != nil {
		log.Printf("Failed to open stream to %s: %v", peerID, err)
		return
	}

	nodeData := node.NodeTracker.nodeData
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
