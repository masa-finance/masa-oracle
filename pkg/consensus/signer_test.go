package consensus

import (
	"encoding/hex"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignData(t *testing.T) {
	privKey, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	assert.NoError(t, err)

	data := []byte("test data")

	signature, err := SignData(privKey, data)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	// Test with nil private key
	signature, err = SignData(nil, data)
	assert.Error(t, err)
	assert.Nil(t, signature)
}

func TestVerifySignature(t *testing.T) {
	privKey, pubKey, err := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	assert.NoError(t, err)

	data := []byte("test data")

	signature, err := SignData(privKey, data)
	assert.NoError(t, err)

	valid, err := VerifySignature(pubKey, data, hex.EncodeToString(signature))
	assert.NoError(t, err)
	assert.True(t, valid)

	// Test with invalid signature
	valid, err = VerifySignature(pubKey, data, "invalid")
	assert.Error(t, err)
	assert.False(t, valid)

	// Test with nil public key
	valid, err = VerifySignature(nil, data, hex.EncodeToString(signature))
	assert.Error(t, err)
	assert.False(t, valid)
}
