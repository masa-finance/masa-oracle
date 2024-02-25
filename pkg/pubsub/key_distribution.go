package pubsub

import (
	"encoding/json"
	"errors"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	masaCrypto "github.com/masa-finance/masa-oracle/pkg/crypto"
	"github.com/sirupsen/logrus"
)

const PublicKeyTopic = "public-key-topic"

// PublicKeyMessage represents the structure of the public key messages.
type PublicKeyMessage struct {
	PublicKey string `json:"publicKey"`
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

// This method can be called directly after creating the Manager instance,
// or at any other appropriate time when you need to publish or update the public key.
func (sm *Manager) PublishNodePublicKey(publicKey string, data, signature []byte) error {
	// Convert the publicKey string to a libp2p PubKey
	pubKey, err := libp2pCrypto.UnmarshalPublicKey([]byte(publicKey)) // Use libp2pCrypto for UnmarshalPublicKey
	if err != nil {
		return errors.New("failed to unmarshal public key")
	}

	// Verify the signature using the VerifySignature function from your custom crypto package
	isValid, err := masaCrypto.VerifySignature(pubKey, data, signature) // Use masaCrypto for VerifySignature
	if err != nil || !isValid {
		return errors.New("unauthorized: signature verification failed")
	}

	// Proceed to publish the public key if the signature is valid
	return sm.PublishNodePublicKey(publicKey, data, signature)
}
