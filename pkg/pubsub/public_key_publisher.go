package pubsub

import (
	"encoding/json"
	"errors"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/masa-finance/masa-oracle/pkg/interfaces"
	"github.com/masa-finance/masa-oracle/pkg/keys"
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
	keyLoader     interfaces.KeyLoader
}

// NewPublicKeyPublisher creates a new instance of PublicKeyPublisher.
func NewPublicKeyPublisher(manager *Manager) *PublicKeyPublisher {
	return &PublicKeyPublisher{
		pubSubManager: manager,
		keyLoader:     &keys.KeyManager{},
	}
}

// PublishNodePublicKey publishes the node's public key to the designated topic.
func (p *PublicKeyPublisher) PublishNodePublicKey(publicKey string, data, signature []byte) error {
	// Convert the publicKey string to a libp2p PubKey
	_, err := libp2pCrypto.UnmarshalPublicKey([]byte(publicKey))
	if err != nil {
		return errors.New("failed to unmarshal public key")
	}

	// Verify the signature using the VerifySignature function from your custom crypto package
	isValid, err := consensus.VerifySignature(p.keyLoader, data, signature)
	if err != nil || !isValid {
		return errors.New("unauthorized: signature verification failed")
	}

	// Additional step: Check if the publicKey is authorized as the first node
	// Assume isFirstNodeAuthorized is defined elsewhere in your package
	isAuthorized := isFirstNodeAuthorized(publicKey)
	if !isAuthorized {
		return errors.New("unauthorized: not the first node or already published")
	}

	// Serialize the public key message
	msg := PublicKeyMessage{
		PublicKey: publicKey,
		Signature: string(signature),
		Data:      string(data),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return errors.New("failed to marshal message")
	}

	// Use the existing Manager to publish the message
	return p.pubSubManager.Publish("nodePublicKey", msgBytes)
}
