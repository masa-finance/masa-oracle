// Package keys provides functionality for managing cryptographic keys used within the Masa Oracle.
// It includes utilities for loading private and public keys from files or environment variables,
// constructing file paths for keys based on configuration, and deriving peer IDs from private keys.

package keys

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetPrivKeyFilePath constructs the file path for the private key using the directory and file name from Viper configuration.
func GetPrivKeyFilePath() string {
	masaDir := viper.GetString(masa.MasaDir)
	privKeyFile := viper.GetString(masa.PrivKeyFile)
	filePath := filepath.Join(masaDir, privKeyFile)
	logrus.WithFields(logrus.Fields{
		"masaDir":     masaDir,
		"privKeyFile": privKeyFile,
		"filePath":    filePath,
	}).Info("Constructed private key file path")
	return filePath
}

// This section of the code is dedicated to handling the loading of private keys from different sources.
// It includes two primary functions:
// 1. LoadPrivKeyFromFilePath: This function reads a hex-encoded private key from a file. The file path is constructed using configuration values for the directory and file name. It then decodes the hex-encoded string into bytes and unmarshals it into a crypto.PrivKey object.
// 2. LoadPrivKeyStringFromEnv: Similar to the first, but instead of reading from a file, it reads the hex-encoded private key from an environment variable. The name of the environment variable is passed as a parameter to the function. It then performs the same decoding and unmarshalling process to obtain the crypto.PrivKey object.

// LoadPrivKeyFromFilePath loads the hex-encoded private key as a crypto.PrivKey object from the configured file path.
func LoadPrivKeyFromFilePath() (crypto.PrivKey, error) {
	filePath := GetPrivKeyFilePath()

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

// LoadPrivKeyFromEnv loads the hex-encoded private key as a crypto.PrivKey object from the configured environment variable.
func LoadPrivKeyFromEnv(envVarName string) (crypto.PrivKey, error) {
	privKeyHex := viper.GetString(envVarName)
	if privKeyHex == "" {
		logrus.WithField("envVarName", envVarName).Error("Environment variable not set or has no value")
		return nil, fmt.Errorf("environment variable %s not set or has no value", envVarName)
	}
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		logrus.WithError(err).WithField("envVarName", envVarName).Error("Failed to decode private key from hex")
		return nil, fmt.Errorf("failed to decode private key from hex: %s", err)
	}
	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		logrus.WithError(err).WithField("envVarName", envVarName).Error("Failed to unmarshal private key")
		return nil, fmt.Errorf("failed to unmarshal private key: %s", err)
	}

	logrus.WithField("envVarName", envVarName).Info("Successfully loaded private key from environment variable")
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

// LoadPubKeyFromEnv derives the crypto.PubKey from the private key loaded from the specified environment variable.
func LoadPubKeyFromEnv(envVarName string) (crypto.PubKey, error) {
	privKey, err := LoadPrivKeyFromEnv(envVarName)
	if err != nil {
		return nil, err
	}

	logrus.WithField("envVarName", envVarName).Info("Successfully derived public key from private key loaded from environment variable")
	return privKey.GetPublic(), nil
}

// This section can be used to get the peerID from the private key
// 1. GetPeerIDFromPrivKeyFilePath: Retrieves and converts the private key from a file path to a peer ID.
// 2. GetPeerIDFromPrivKeyEnv: Does the same but sources the private key from an environment variable.
// Both functions load the private key, derive the public key, and generate a peer ID from it, returning the ID as a string.

// GetPeerIDFromPrivKeyFilePath derives the peer ID from the private key loaded from the configured file path.
func GetPeerIDFromPrivKeyFilePath() (string, error) {
	privKey, err := LoadPrivKeyFromFilePath()
	if err != nil {
		return "", err
	}

	pubKey := privKey.GetPublic()
	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		logrus.WithError(err).Error("Failed to derive peer ID from public key")
		return "", err
	}

	logrus.Info("Successfully derived peer ID from private key loaded from file path")
	return peerID.String(), nil
}

// GetPeerIDFromPrivKeyEnv derives the peer ID from the private key loaded from the specified environment variable.
func GetPeerIDFromPrivKeyEnv(envVarName string) (string, error) {
	privKey, err := LoadPrivKeyFromEnv(envVarName)
	if err != nil {
		return "", err
	}

	pubKey := privKey.GetPublic()
	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		logrus.WithError(err).WithField("envVarName", envVarName).Error("Failed to derive peer ID from public key")
		return "", err
	}

	logrus.WithField("envVarName", envVarName).Info("Successfully derived peer ID from private key loaded from environment variable")
	return peerID.String(), nil
}
