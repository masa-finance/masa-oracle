package crypto

import (
	"crypto/rand"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// GenerateKeyPair generates a new Ed25519 key pair for the writer node.
func GenerateKeyPair() (crypto.PrivKey, crypto.PubKey, error) {
	return crypto.GenerateEd25519Key(rand.Reader)
}

// SignData signs the data using the given private key.
func SignData(privKey crypto.PrivKey, data []byte) ([]byte, error) {
	return privKey.Sign(data)
}

// VerifySignature verifies the signature of the data using the given public key.
func VerifySignature(pubKey crypto.PubKey, data []byte, signature []byte) (bool, error) {
	return pubKey.Verify(data, signature)
}
