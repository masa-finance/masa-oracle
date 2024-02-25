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
	"github.com/spf13/viper"
)

// GetPrivKeyFilePath constructs the file path for the private key using the directory and file name from Viper configuration.
func GetPrivKeyFilePath() string {
	masaDir := viper.GetString(masa.MasaDir)
	privKeyFile := viper.GetString(masa.PrivKeyFile)

	return filepath.Join(masaDir, privKeyFile)
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
		return nil, fmt.Errorf("failed to read private key from file: %s", err)
	}
	privKeyBytes, err := hex.DecodeString(string(hexPrivKeyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key from hex: %s", err)
	}
	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %s", err)
	}

	return privKey, nil
}

// LoadPrivKeyFromEnv loads the hex-encoded private key as a crypto.PrivKey object from the configured environment variable.
func LoadPrivKeyFromEnv(envVarName string) (crypto.PrivKey, error) {
	privKeyHex := viper.GetString(envVarName)
	if privKeyHex == "" {
		return nil, fmt.Errorf("environment variable %s not set or has no value", envVarName)
	}
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key from hex: %s", err)
	}
	privKey, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %s", err)
	}

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

	return privKey.GetPublic(), nil
}

// LoadPubKeyFromEnv derives the crypto.PubKey from the private key loaded from the specified environment variable.
func LoadPubKeyFromEnv(envVarName string) (crypto.PubKey, error) {
	privKey, err := LoadPrivKeyFromEnv(envVarName)
	if err != nil {
		return nil, err
	}

	return privKey.GetPublic(), nil
}

// Developer Notes:
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
		return "", err
	}

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
		return "", err
	}

	return peerID.String(), nil
}
