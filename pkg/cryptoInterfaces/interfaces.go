package cryptoInterfaces

import "github.com/libp2p/go-libp2p/core/crypto"

// KeyLoader defines the interface for loading private and public keys.
type KeyLoader interface {
	LoadPrivKey() (crypto.PrivKey, error)
	LoadPubKey() (crypto.PubKey, error)
}
