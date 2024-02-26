package pubsub

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"
)

// PublicKeyMessage represents the structure of the public key messages.
type PublicKeyMessage struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Data      string `json:"data"`
}

// PublicKeyPublisher uses the existing Manager to publish public keys.
type PublicKeyPublisher struct {
	pubSubManager *Manager
	pubKey        libp2pCrypto.PubKey
}

// topicPublicKeyMap maps topics to their associated public keys.
var topicPublicKeyMap = make(map[string]string)

// NewPublicKeyPublisher creates a new instance of PublicKeyPublisher.
func NewPublicKeyPublisher(manager *Manager, privKeyPath string, pubKey libp2pCrypto.PubKey) *PublicKeyPublisher {
	logrus.Info("Creating new PublicKeyPublisher")
	return &PublicKeyPublisher{
		pubSubManager: manager,
		pubKey:        pubKey,
	}
}

// PublishNodePublicKey publishes the node's public key to the designated topic.
func (p *PublicKeyPublisher) PublishNodePublicKey(publicKey string, data, signature []byte) error {
	topic := "bootNodePublicKey"
	logrus.Infof("Publishing public key to topic: %s", topic)

	// Check if a public key has already been published to the topic
	existingPubKey, exists := topicPublicKeyMap[topic]

	if exists {
		logrus.Infof("Public key already exists for topic: %s. Verifying signature.", topic)
		// If a public key exists, verify the signature against the existing public key
		pubKeyBytes, err := hex.DecodeString(existingPubKey)
		if err != nil {
			logrus.WithError(err).Error("Failed to decode existing public key")
			return err
		}
		pubKey, err := libp2pCrypto.UnmarshalPublicKey(pubKeyBytes)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal existing public key")
			return err
		}
		isValid, err := pubKey.Verify(data, signature)
		if err != nil || !isValid {
			logrus.WithError(err).Error("Unauthorized: Signature verification failed or signature is invalid")
			return errors.New("unauthorized: only the owner of the public key can publish changes")
		}
		logrus.Info("Signature verified successfully")
	} else {
		logrus.Infof("No existing public key for topic: %s. This is the initial publication.", topic)
		// If no public key is associated with the topic, this is the initial publication
		topicPublicKeyMap[topic] = publicKey
	}

	// Serialize the public key message
	msg := PublicKeyMessage{
		PublicKey: publicKey,
		Signature: hex.EncodeToString(signature),
		Data:      string(data),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal public key message")
		return errors.New("failed to marshal message")
	}

	// Use the existing Manager to publish the message
	logrus.Info("Publishing message to pubSubManager")
	return p.pubSubManager.Publish(topic, msgBytes)
}
