package staking

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"math/big"
)

func VerifyStakingSignature(signature, pubKey, data string) bool {
	// Convert the public key to an ecdsa.PublicKey
	block, _ := pem.Decode([]byte(pubKey))
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubInterface.(*ecdsa.PublicKey)

	// Convert the signature to r and s values
	r := new(big.Int)
	s := new(big.Int)
	sigLen := len(signature)
	r.SetBytes([]byte(signature[:sigLen/2]))
	s.SetBytes([]byte(signature[sigLen/2:]))

	// Verify the signature
	hash := sha256.Sum256([]byte(data))
	isValid := ecdsa.Verify(pub, hash[:], r, s)

	return isValid
}
