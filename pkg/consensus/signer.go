package consensus

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"
)

// SignData signs the data using the provided private key.
func SignData(privKey crypto.PrivKey, data []byte) ([]byte, error) {
	if privKey == nil {
		logrus.Error("Private key is nil")
		return nil, fmt.Errorf("private key is nil")
	}

	logrus.Info("Signing data")

	signedData, err := privKey.Sign(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to sign data")
		return nil, err
	}

	logrus.Info("Data signed successfully")
	return signedData, nil
}

// VerifySignature verifies the signature of the data using the provided public key.
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
