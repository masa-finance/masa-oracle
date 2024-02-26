package pubsub

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"

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
	privKeyPath   string // Assuming the private key path is provided during initialization
}

// NewPublicKeyPublisher creates a new instance of PublicKeyPublisher.
func NewPublicKeyPublisher(manager *Manager, privKeyPath string) *PublicKeyPublisher {
	return &PublicKeyPublisher{
		pubSubManager: manager,
		privKeyPath:   privKeyPath,
	}
}

// loadPrivateKey loads the private key from the file path.
func (p *PublicKeyPublisher) loadPrivateKey() (libp2pCrypto.PrivKey, error) {
	hexPrivKeyBytes, err := os.ReadFile(p.privKeyPath)
	if err != nil {
		return nil, err
	}
	privKeyBytes, err := hex.DecodeString(string(hexPrivKeyBytes))
	if err != nil {
		return nil, err
	}
	return libp2pCrypto.UnmarshalPrivateKey(privKeyBytes)
}

// verifySignature verifies the signature of the data using the public key.
func (p *PublicKeyPublisher) verifySignature(publicKeyStr string, data, signature []byte) (bool, error) {
	pubKeyBytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return false, err
	}

	pubKey, err := libp2pCrypto.UnmarshalPublicKey(pubKeyBytes)
	if err != nil {
		return false, err
	}

	isValid, err := pubKey.Verify(data, signature)
	if err != nil {
		return false, err
	}

	return isValid, nil
}

// PublishNodePublicKey publishes the node's public key to the designated topic.
func (p *PublicKeyPublisher) PublishNodePublicKey(publicKey string, data, signature []byte) error {
	// Verify the signature
	isValid, err := p.verifySignature(publicKey, data, signature)
	if err != nil || !isValid {
		return errors.New("unauthorized: signature verification failed")
	}

	// Additional checks and message publishing logic remains the same
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
	return p.pubSubManager.Publish("nodePublicKey", msgBytes)
}
