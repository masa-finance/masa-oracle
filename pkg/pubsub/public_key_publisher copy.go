package pubsub

import (
	"errors"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	masaCrypto "github.com/masa-finance/masa-oracle/pkg/crypto"
)

// PublicKeyMessage represents the structure of the public key messages.
type PublicKeyMessage struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Data      string `json:"data"`
}

// PublishNodePublicKey can be called directly after creating the Manager instance,
// or at any other appropriate time when you need to publish or update the public key.
func (sm *Manager) PublishNodePublicKey(publicKey string, data, signature []byte) error {
	// Convert the publicKey string to a libp2p PubKey
	pubKey, err := libp2pCrypto.UnmarshalPublicKey([]byte(publicKey))
	if err != nil {
		return errors.New("failed to unmarshal public key")
	}

	// Verify the signature using the VerifySignature function from your custom crypto package
	isValid, err := masaCrypto.VerifySignature(pubKey, data, signature)
	if err != nil || !isValid {
		return errors.New("unauthorized: signature verification failed")
	}

	// Proceed to publish the public key if the signature is valid
	// Implementation of actual publishing logic goes here
	// For example, using a pubsub system to publish the PublicKeyMessage

	return nil
}
