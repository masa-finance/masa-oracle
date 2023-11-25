package staking

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

func VerifyStakingSignature(signature []byte, pubKey *ecdsa.PublicKey, data []byte) bool {
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	hash := crypto.Keccak256Hash(data)

	return ecdsa.Verify(pubKey, hash.Bytes(), r, s)
}
