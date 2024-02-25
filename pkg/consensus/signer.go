package consensus

import (
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/masa-finance/masa-oracle/pkg/keys"
	"github.com/sirupsen/logrus"
)

// SignData signs the data using the private key loaded from a configured source.
func SignData(data []byte, useEnv bool, envVarName string) ([]byte, error) {
	var privKey crypto.PrivKey
	var err error

	if useEnv {
		// Load the private key from an environment variable
		privKey, err = keys.LoadPrivKeyFromEnv(envVarName)
		logrus.WithFields(logrus.Fields{"envVarName": envVarName}).Info("Loading private key from environment variable")
	} else {
		// Load the private key from the file path
		privKey, err = keys.LoadPrivKeyFromFilePath()
		logrus.Info("Loading private key from file path")
	}

	if err != nil {
		logrus.WithError(err).Error("Failed to load private key")
		return nil, err
	}

	signedData, err := privKey.Sign(data)
	if err != nil {
		logrus.WithError(err).Error("Failed to sign data")
		return nil, err
	}

	logrus.Info("Data signed successfully")
	return signedData, nil
}

// VerifySignature verifies the signature of the data using the public key loaded from a configured source.
func VerifySignature(data []byte, signature []byte, useEnv bool, envVarName string) (bool, error) {
	var pubKey crypto.PubKey
	var err error

	if useEnv {
		// Load the public key from an environment variable
		pubKey, err = keys.LoadPubKeyFromEnv(envVarName)
		logrus.WithFields(logrus.Fields{"envVarName": envVarName}).Info("Loading public key from environment variable")
	} else {
		// Load the public key from the file path
		pubKey, err = keys.LoadPubKeyFromFilePath()
		logrus.Info("Loading public key from file path")
	}

	if err != nil {
		logrus.WithError(err).Error("Failed to load public key")
		return false, err
	}

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
