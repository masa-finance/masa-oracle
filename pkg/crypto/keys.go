package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"

	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"
)

func getPrivateKeyFromEnv(envKey string) (privKey crypto.PrivKey, err error) {
	rawKey, err := hex.DecodeString(envKey)
	if err != nil {
		return nil, logAndReturnError("error decoding private key: %s", err)
	}
	privKey, err = crypto.UnmarshalPrivateKey(rawKey)
	if err != nil {
		return nil, logAndReturnError("error unmarshalling private key: %s", err)
	}
	logrus.Infof("Loaded private key from environment")
	return privKey, nil
}

func getPrivateKeyFromFile(keyFile string) (crypto.PrivKey, error) {
	// Check if the private key file exists
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, logAndReturnError("Error reading private key file: %s", err)
	}

	// Decode the private key from the file
	rawKey, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, logAndReturnError("Error decoding private key: %s", err)
	}

	// Unmarshal the private key
	privKey, err := crypto.UnmarshalPrivateKey(rawKey)
	if err != nil {
		return nil, logAndReturnError("Error unmarshalling private key: %s", err)
	}
	logrus.Infof("Loaded private key from %s", keyFile)
	return privKey, nil
}

func generateNewPrivateKey(keyFile string) (crypto.PrivKey, error) {
	// Generate a new private key
	privKey, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 2048)
	if err != nil {
		return nil, logAndReturnError("Error generating new private key: %s", err)
	}
	data, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, logAndReturnError("Error marshalling private key: %s", err)
	}
	encodedKey := hex.EncodeToString(data)

	// Save the private key to the file
	if err := os.WriteFile(keyFile, []byte(encodedKey), 0600); err != nil {
		return nil, logAndReturnError("Error saving private key to file: %s", err)
	}
	logrus.Infof("Generated and saved a new private key to %s: %s", keyFile, privKey)
	return privKey, nil
}

func GetOrCreatePrivateKey(keyFile string) (crypto.PrivKey, *ecdsa.PrivateKey, string, error) {
	var privKey crypto.PrivKey
	var err error
	envKey := os.Getenv("PRIVATE_KEY")
	if envKey != "" {
		privKey, err = getPrivateKeyFromEnv(envKey)
		if err != nil {
			return nil, nil, "", err
		}
	} else {
		// Check if the private key file exists
		privKey, err = getPrivateKeyFromFile(keyFile)
		if err != nil {
			privKey, err = generateNewPrivateKey(keyFile)
			if err != nil {
				return nil, nil, "", err
			}
		}
	}
	// After obtaining the libp2p privKey, convert it to an ECDSA private key
	ecdsaPrivKey, err := Libp2pPrivateKeyToEcdsa(privKey)
	if err != nil {
		return nil, nil, "", err
	}
	// Save the ECDSA private key in the same directory as the libp2p key
	ecdsaKeyFilePath := keyFile + ".ecdsa"
	ecdsaKeyBytes := ethCrypto.FromECDSA(ecdsaPrivKey)
	ecdsaKeyHex := hex.EncodeToString(ecdsaKeyBytes)
	if err := os.WriteFile(ecdsaKeyFilePath, []byte(ecdsaKeyHex), 0600); err != nil {
		logrus.Errorf("Error saving ECDSA private key to file: %s\n", err)
		return privKey, ecdsaPrivKey, "", err
	}
	logrus.Infof("Saved ECDSA private key to %s", ecdsaKeyFilePath)
	ethAddress := ethCrypto.PubkeyToAddress(ecdsaPrivKey.PublicKey).Hex()
	return privKey, ecdsaPrivKey, ethAddress, nil
}

func Libp2pPrivateKeyToEcdsa(privKey crypto.PrivKey) (*ecdsa.PrivateKey, error) {
	// After obtaining the libp2p privKey, convert it to an ECDSA private key
	raw, err := privKey.Raw()
	if err != nil {
		logrus.Errorf("Error getting raw private key: %s\n", err)
		return nil, err
	}

	ecdsaPrivKey, err := ethCrypto.ToECDSA(raw)
	if err != nil {
		logrus.Errorf("Error converting to ECDSA private key: %s\n", err)
		return nil, err
	}

	return ecdsaPrivKey, nil
}

func Libp2pPubKeyToEcdsa(pubKey crypto.PubKey) (*ecdsa.PublicKey, error) {
	raw, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}
	secpPubKey, err := secp.ParsePubKey(raw)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	ecdsaPubKey := secpPubKey.ToECDSA()
	return ecdsaPubKey, nil
}

func Libp2pPubKeyToEthAddress(pubKey crypto.PubKey) (string, error) {
	ecdsaPublic, err := Libp2pPubKeyToEcdsa(pubKey)
	if err != nil {
		return "", err
	}
	ethAddress := ethCrypto.PubkeyToAddress(*ecdsaPublic).Hex()

	return ethAddress, nil
}

// VerifyEthereumCompatibility Really only useful for printing key values, I think this could be removed.
func VerifyEthereumCompatibility(privKey crypto.PrivKey) (string, error) {
	ecdsaPrivKey, err := Libp2pPrivateKeyToEcdsa(privKey)
	if err != nil {
		return "", err
	}
	ethAddress := ethCrypto.PubkeyToAddress(ecdsaPrivKey.PublicKey).Hex()

	// Print the private key in hexadecimal format
	data, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return "", err
	}
	fmt.Printf("Private key: \n%s\n", hex.EncodeToString(data))
	// Print the public key and address
	fmt.Println("Ethereum public key:", ecdsaPrivKey.PublicKey)
	// Derive the Ethereum address from the private key
	fmt.Println("Ethereum address:", ethCrypto.PubkeyToAddress(ecdsaPrivKey.PublicKey).Hex())

	return ethAddress, nil
}

func logAndReturnError(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	logrus.Error(err)
	return err
}

func GetPublicKeyForHost(host host.Host) (publicKeyHex string, err error) {
	pubKey := host.Peerstore().PubKey(host.ID())
	if pubKey == nil {
		logrus.WithFields(logrus.Fields{
			"Peer": host.ID().String(),
		}).Warn("No public key found for peer")
	} else {
		publicKeyHex, err = Libp2pPubKeyToEthAddress(pubKey)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Peer": host.ID().String(),
			}).Warnf("Error getting public key %v", err)
		}
	}
	return
}
