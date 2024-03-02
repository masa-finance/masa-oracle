// This package provides functionality for managing cryptographic keys used within the Masa Oracle.
// It includes utilities for loading private and public keys from files or environment variables,
// constructing file paths for keys based on configuration, and deriving peer IDs from private keys.

package masa

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

// KeyManager implements the KeyLoader interface used in the interfaces package to load keys from file path.

// @bob - can you help add in a case where it load from the ENV fist if the filepath isn't availble if you think this is necessarey?
type KeyManager struct{}

func (km *KeyManager) LoadPrivKey() (crypto.PrivKey, error) {
	return LoadPrivKeyFromFilePath()
}

func (km *KeyManager) LoadPubKey() (crypto.PubKey, error) {
	return LoadPubKeyFromFilePath()
}

// This section of the code is dedicated to handling the loading of private keys from different sources.
// It includes two primary functions:
// 1. LoadPrivKeyFromFilePath: This function reads a hex-encoded private key from a file. The file path is constructed using configuration values for the directory and file name. It then decodes the hex-encoded string into bytes and unmarshals it into a crypto.PrivKey object.
// 2. LoadPrivKeyStringFromEnv: Similar to the first, but instead of reading from a file, it reads the hex-encoded private key from an environment variable. The name of the environment variable is passed as a parameter to the function. It then performs the same decoding and unmarshalling process to obtain the crypto.PrivKey object.

// LoadPrivKeyFromFilePath loads the hex-encoded private key as a crypto.PrivKey object from the configured file path.
func LoadPrivKeyFromFilePath() (crypto.PrivKey, error) {
	filePath := config.GetInstance().PrivateKeyFile

	hexPrivKeyBytes, err := os.ReadFile(filePath)
	if err != nil {
		logrus.WithError(err).WithField("filePath", filePath).Error("Failed to read private key from file")
		return nil, fmt.Errorf("failed to read private key from file: %s", err)
	}
	privKeyBytes, err := hex.DecodeString(string(hexPrivKeyBytes))
	if err != nil {
		logrus.WithError(err).Error("Failed to decode private key from hex")
		return nil, fmt.Errorf("failed to decode private key from hex: %s", err)
	}
	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal private key")
		return nil, fmt.Errorf("failed to unmarshal private key: %s", err)
	}

	logrus.Info("Successfully loaded private key from file path")
	return privKey, nil
}

// This section provides two methods for accessing the crypto.PubKey:
// 1. LoadPubKeyFromFilePath: This method utilizes the LoadPrivKeyFromFilePath function to load the private key from a file path specified in the configuration. It then derives the public key from this private key.
// 2. LoadPubKeyFromEnv: This method, on the other hand, loads the private key from an environment variable whose name is passed as a parameter. It similarly derives the public key from the loaded private key.

// LoadPubKeyFromFilePath derives the crypto.PubKey from the private key loaded from the configured file path.
func LoadPubKeyFromFilePath() (crypto.PubKey, error) {
	privKey, err := LoadPrivKeyFromFilePath()
	if err != nil {
		return nil, err
	}

	logrus.Info("Successfully derived public key from private key loaded from file path")
	return privKey.GetPublic(), nil
}

// PubKeyToString converts a crypto.PubKey to its string representation.
func PubKeyToString(pubKey crypto.PubKey) (string, error) {
	pubKeyBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal public key")
		return "", fmt.Errorf("failed to marshal public key: %s", err)
	}
	return hex.EncodeToString(pubKeyBytes), nil
}
