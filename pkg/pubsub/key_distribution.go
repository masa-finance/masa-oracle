package pubsub

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

const PublicKeyTopic = "public-key-topic"

// PublicKeyMessage represents the structure of the public key messages.
type PublicKeyMessage struct {
	PublicKey string `json:"publicKey"`
}

// PublishPublicKey publishes the given public key to the public key topic.
func (sm *Manager) PublishPublicKey(publicKey string) error {
	msg := PublicKeyMessage{PublicKey: publicKey}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("Failed to marshal public key message: %v", err)
		return err
	}
	if err := sm.gossipSub.Publish(PublicKeyTopic, msgBytes); err != nil {
		logrus.Errorf("Failed to publish public key to topic %s: %v", PublicKeyTopic, err)
		return err
	}
	logrus.Infof("Successfully published public key to topic %s", PublicKeyTopic)
	return nil
}

// SubscribeToPublicKeyTopic subscribes to the public key topic to receive updates.
func (sm *Manager) SubscribeToPublicKeyTopic() error {
	sub, err := sm.gossipSub.Subscribe(PublicKeyTopic)
	if err != nil {
		logrus.Errorf("Failed to subscribe to public key topic %s: %v", PublicKeyTopic, err)
		return err
	}
	logrus.Infof("Successfully subscribed to public key topic %s", PublicKeyTopic)
	go func() {
		for {
			msg, err := sub.Next(sm.ctx)
			if err != nil {
				logrus.Errorf("Error reading from public key topic: %v", err)
				continue
			}
			var publicKeyMsg PublicKeyMessage
			if err := json.Unmarshal(msg.Data, &publicKeyMsg); err != nil {
				logrus.Errorf("Error unmarshalling public key message: %v", err)
				continue
			}
			logrus.Infof("Successfully received and unmarshalled public key message from topic %s", PublicKeyTopic)
			// Process the received public key, e.g., verify and update local copy
		}
	}()
	return nil
}
