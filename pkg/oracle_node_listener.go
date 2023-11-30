package masa

import (
	"encoding/json"
	"log"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
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
