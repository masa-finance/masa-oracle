package interfaces

import "github.com/libp2p/go-libp2p/core/crypto"

// KeyLoader defines the interface for loading private and public keys.
type KeyLoader interface {
	LoadPrivKey() (crypto.PrivKey, error)
	LoadPubKey() (crypto.PubKey, error)
}

// SignatureVerifier defines the interface for verifying signatures.
type SignatureVerifier interface {
	VerifySignature(data []byte, signature []byte) (bool, error)
}
