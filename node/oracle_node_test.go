package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNodeSignature(t *testing.T) {
	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Create a data string
	data := "This is some test data"

	// Hash the data
	hash := crypto.Keccak256Hash([]byte(data))

	// Sign the hash with the private key
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new OracleNode
	node := &OracleNode{
		Signature: hex.EncodeToString(sig),
	}

	// Check if the node is a publisher
	if !node.IsPublisher() {
		t.Errorf("Expected node to be a publisher, but it's not")
	}
}
