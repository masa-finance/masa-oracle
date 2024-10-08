package tee

/*

This is a wrapper package just to ease out adding logics that
should apply to all callers of the sealer.

XXX: Currently it is equivalent as calling the library directly,
and provides just syntax sugar.

*/

import "github.com/edgelesssys/ego/ecrypto"

// Seal uses the TEE Product Key to encrypt the plaintext
// The Product key is the one bound to the signer pubkey
func Seal(plaintext []byte) ([]byte, error) {
	return ecrypto.SealWithProductKey(plaintext, nil)
}

func Unseal(encryptedText []byte) ([]byte, error) {
	return ecrypto.Unseal(encryptedText, nil)
}
