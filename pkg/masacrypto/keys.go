package masacrypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"

	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"
)

// Package crypto provides cryptographic utilities used by the KeyManager.
// This file, keys.go, contains a set of private helper functions that are essential
// for managing cryptographic operations related to keys. These functions include
// generating new private keys, loading private keys from environment variables or files,
// marshalling and unmarshalling private keys, and converting keys between different formats.
//
// The functions defined in this file are designed to be private to the package and are
// primarily intended to support the operations of the KeyManager, defined in key_manager.go.
// The KeyManager acts as a centralized entity for managing cryptographic keys within the application,
// ensuring secure and consistent access to these keys across different components.
//
// The emphasis on keeping these functions private is to encapsulate the cryptographic operations
// within the package, allowing for a controlled interface exposed by the KeyManager. This approach
// helps in maintaining the integrity and security of key management processes, by restricting direct
// access to the underlying cryptographic operations and promoting the use of the KeyManager as the
// primary interface for key-related activities.

func getPrivateKeyFromEnv(envKey string) (privKey crypto.PrivKey, err error) {
	rawKey, err := hex.DecodeString(envKey)
	if err != nil {
		return nil, logAndReturnError("error decoding private key: %s", err)
	}
	privKey, err = crypto.UnmarshalPrivateKey(rawKey)
	if err != nil {
		return nil, logAndReturnError("error unmarshalling private key: %s", err)
	}
	logrus.Infof("[+] Loaded private key from environment")
	return privKey, nil
}

func getPrivateKeyFromFile(keyFile string) (crypto.PrivKey, error) {
	// Check if the private key file exists
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, logAndReturnError("[-] Error reading private key file: %s", err)
	}

	// Decode the private key from the file
	rawKey, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, logAndReturnError("[-] Error decoding private key: %s", err)
	}

	// Unmarshal the private key
	privKey, err := crypto.UnmarshalPrivateKey(rawKey)
	if err != nil {
		return nil, logAndReturnError("[-] Error unmarshalling private key: %s", err)
	}
	logrus.Infof("[+] Loaded private key from %s", keyFile)
	return privKey, nil
}

func generateNewPrivateKey(keyFile string) (crypto.PrivKey, error) {
	// Generate a new private key
	privKey, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 2048)
	if err != nil {
		return nil, logAndReturnError("[-] Error generating new private key: %s", err)
	}

	encodedKey, err := getHexEncodedPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	// Save the private key to the file
	if err := os.WriteFile(keyFile, []byte(encodedKey), 0600); err != nil {
		return nil, logAndReturnError("[-] Error saving private key to file: %s", err)
	}
	logrus.Infof("[+] Generated and saved a new private key to %s: %s", keyFile, privKey)
	return privKey, nil
}

func getHexEncodedPrivateKey(privKey crypto.PrivKey) (string, error) {
	data, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return "", logAndReturnError("[-] Error marshalling private key: %s", err)
	}
	encodedKey := hex.EncodeToString(data)
	return encodedKey, nil
}

func getHexEncodedPublicKey(publicKey crypto.PubKey) (string, error) {
	pubKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	if err != nil {
		return "", logAndReturnError("[-] Error marshalling public key: %s", err)
	}
	encodedKey := hex.EncodeToString(pubKeyBytes)
	return encodedKey, nil
}

func libp2pPrivateKeyToEcdsa(privKey crypto.PrivKey) (*ecdsa.PrivateKey, error) {
	// After obtaining the libp2p privKey, convert it to an ECDSA private key
	raw, err := privKey.Raw()
	if err != nil {
		logrus.Errorf("[-] Error getting raw private key: %s\n", err)
		return nil, err
	}

	ecdsaPrivKey, err := ethCrypto.ToECDSA(raw)
	if err != nil {
		logrus.Errorf("[-] Error converting to ECDSA private key: %s\n", err)
		return nil, err
	}

	return ecdsaPrivKey, nil
}

func libp2pPubKeyToEcdsa(pubKey crypto.PubKey) (*ecdsa.PublicKey, error) {
	raw, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	secpPubKey, err := secp.ParsePubKey(raw)
	if err != nil {
		return nil, fmt.Errorf("[-] error parsing public key: %v", err)
	}

	ecdsaPubKey := secpPubKey.ToECDSA()
	return ecdsaPubKey, nil
}

func Libp2pPubKeyToEthAddress(pubKey crypto.PubKey) (string, error) {
	ecdsaPublic, err := libp2pPubKeyToEcdsa(pubKey)
	if err != nil {
		return "", err
	}
	ethAddress := ethCrypto.PubkeyToAddress(*ecdsaPublic).Hex()

	return ethAddress, nil
}

func logAndReturnError(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	logrus.Error("[-] logAndReturnError ", err)
	return err
}

func saveEcdesaPrivateKeyToFile(ecdsaPrivKey *ecdsa.PrivateKey, keyFile string) error {
	ecdsaKeyBytes := ethCrypto.FromECDSA(ecdsaPrivKey)
	ecdsaKeyHex := hex.EncodeToString(ecdsaKeyBytes)
	if err := os.WriteFile(keyFile, []byte(ecdsaKeyHex), 0600); err != nil {
		logrus.Errorf("[-] Error saving ECDSA private key to file: %s", err)
		return err
	}
	logrus.Infof("[+] Saved ECDSA private key to %s", keyFile)
	return nil
}
