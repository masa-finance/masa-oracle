package pubsub

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

const PublicKeyTopic = "public-key-topic"

func (sm *Manager) SubscribeToPublicKeyTopic() error {
	topic, err := sm.gossipSub.Join(PublicKeyTopic)
	if err != nil {
		logrus.Errorf("Failed to join public key topic %s: %v", PublicKeyTopic, err)
		return err
	}
	sub, err := topic.Subscribe()
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
		}
	}()
	return nil
}
