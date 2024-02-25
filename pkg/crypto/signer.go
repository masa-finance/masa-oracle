package crypto

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// SignData signs the data using the private key directly obtained from keys.go.
// This function assumes that the private key is already loaded or generated by keys.go.
func SignData(privKey crypto.PrivKey, data []byte) ([]byte, error) {
	if privKey == nil {
		// Use fmt.Errorf or errors.New to return an error object
		return nil, fmt.Errorf("private key is nil")
	}
	return privKey.Sign(data)
}

// VerifySignature verifies the signature of the data using the public key.
// This function directly utilizes the public key obtained from keys.go or similar mechanisms.
func VerifySignature(pubKey crypto.PubKey, data []byte, signature []byte) (bool, error) {
	if pubKey == nil {
		// Use fmt.Errorf or errors.New to return an error object
		return false, fmt.Errorf("public key is nil")
	}
	return pubKey.Verify(data, signature)
}
