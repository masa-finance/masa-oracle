package consensus

// Developer Notes:
// This package provides cryptographic functionalities for signing and verifying data within the consensus mechanism of the system.
// It leverages the go-libp2p core crypto library to handle cryptographic operations, ensuring secure data handling.
// - SignData: Signs arbitrary data using a private key, ensuring the data integrity and source authenticity.
// - VerifySignature: Validates the signature of the data against a public key, confirming the data's integrity and origin.

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"
)

// SignData signs the data using the provided private key and returns the signature.
func SignData(privKey crypto.PrivKey, data []byte) ([]byte, error) {
	if privKey == nil {
		logrus.Error("Private key is nil")
		return nil, fmt.Errorf("private key is nil")
	}

	logrus.Info("Signing data")

	signature, err := privKey.Sign(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to sign data")
		return nil, err
	}

	logrus.Info("Data signed successfully")
	return signature, nil
}

// VerifySignature verifies the signature of the data using the signers public key, the data that was signed, and the signature.
func VerifySignature(pubKey crypto.PubKey, data []byte, signature []byte) (bool, error) {
	if pubKey == nil {
		logrus.Error("Public key is nil")
		return false, fmt.Errorf("public key is nil")
	}

	logrus.Info("Verifying signature")

	verified, err := pubKey.Verify(data, signature)
	if err != nil {
		logrus.WithError(err).Error("Failed to verify signature")
		return false, err
	}

	if verified {
		logrus.Info("Signature verified successfully")
	} else {
		logrus.Info("Signature verification failed")
	}
	return verified, nil
}
