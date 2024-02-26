package pubsub

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
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
	return &PublicKeyPublisher{
		pubSubManager: manager,
		pubKey:        pubKey,
	}
}

// PublishNodePublicKey publishes the node's public key to the designated topic.
func (p *PublicKeyPublisher) PublishNodePublicKey(publicKey string, data, signature []byte) error {
	topic := "bootNodePublicKey"

	// Check if a public key has already been published to the topic
	existingPubKey, exists := topicPublicKeyMap[topic]

	if exists {
		// If a public key exists, verify the signature against the existing public key
		pubKeyBytes, err := hex.DecodeString(existingPubKey)
		if err != nil {
			return err
		}
		pubKey, err := libp2pCrypto.UnmarshalPublicKey(pubKeyBytes)
		if err != nil {
			return err
		}
		isValid, err := pubKey.Verify(data, signature)
		if err != nil || !isValid {
			return errors.New("unauthorized: only the owner of the public key can publish changes")
		}
	} else {
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
		return errors.New("failed to marshal message")
	}

	// Use the existing Manager to publish the message
	return p.pubSubManager.Publish(topic, msgBytes)
}
