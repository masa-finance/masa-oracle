package consensus

import (
	"github.com/masa-finance/masa-oracle/pkg/interfaces"
	"github.com/sirupsen/logrus"
)

// SignData signs the data using the private key loaded using a KeyLoader.
func SignData(kl interfaces.KeyLoader, data []byte) ([]byte, error) {
	privKey, err := kl.LoadPrivKey()
	if err != nil {
		logrus.WithError(err).Error("Failed to load private key")
		return nil, err
	}

	logrus.Info("Private key loaded")

	signedData, err := privKey.Sign(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to sign data")
		return nil, err
	}

	logrus.Info("Data signed successfully")
	return signedData, nil
}

// VerifySignature verifies the signature of the data using the public key loaded using a KeyLoader.
func VerifySignature(kl interfaces.KeyLoader, data []byte, signature []byte) (bool, error) {
	pubKey, err := kl.LoadPubKey()
	if err != nil {
		logrus.WithError(err).Error("Failed to load public key")
		return false, err
	}

	logrus.Info("Public key loaded")

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
