// This file is part of the pubsub package and contains the implementation for subscribing to the public key topic on the gossip network.
// It is responsible for managing the subscription to the public key distribution topic, enabling the node to receive and process public key updates from other nodes in the network.
// This ensures the node has the latest public keys necessary for secure communication.

package pubsub

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// PublicKeyTopic is the name of the topic used for public key distribution.
const PublicKeyTopic = "public-key-topic"

// SubscribeToPublicKeyTopic subscribes the node to the public key topic on the gossip network.
// This function is crucial for the node to receive updates about new public keys from other nodes.
// It ensures that the node stays updated with the latest public keys which are essential for secure communication.
func (sm *Manager) SubscribeToPublicKeyTopic() error {
	// Attempt to join the public key topic on the gossip network.
	topic, err := sm.gossipSub.Join(PublicKeyTopic)
	if err != nil {
		logrus.Errorf("Failed to join public key topic %s: %v", PublicKeyTopic, err)
		return err
	}
	// Subscribe to the topic to start receiving messages.
	sub, err := topic.Subscribe()
	if err != nil {
		logrus.Errorf("Failed to subscribe to public key topic %s: %v", PublicKeyTopic, err)
		return err
	}

	logrus.Infof("Successfully subscribed to public key topic %s", PublicKeyTopic)

	// Start a goroutine to continuously read messages from the topic.
	go func() {
		for {
			// Wait for the next message.
			msg, err := sub.Next(sm.ctx)
			if err != nil {
				logrus.Errorf("Error reading from public key topic: %v", err)
				continue
			}
			// Attempt to unmarshal the message into a PublicKeyMessage struct.
			var publicKeyMsg PublicKeyMessage
			if err := json.Unmarshal(msg.Data, &publicKeyMsg); err != nil {
				logrus.Errorf("Error unmarshalling public key message: %v", err)
				continue
			}
			logrus.Infof("Successfully received and unmarshalled public key message from topic %s", PublicKeyTopic)
		}
	}()
	return nil
}
