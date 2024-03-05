package masacrypto

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

// KeyManager is meant to simplify the management of cryptographic keys used in the application.
// It is a singleton, ensuring that there is only one instance of KeyManager throughout the
// application lifecycle, which provides a consistent and secure access point to cryptographic
// keys across different components of the application.
//
// The KeyManager is designed to encapsulate the complexity of cryptographic operations,
// such as key generation, loading, saving, and conversion between different key formats.
// It leverages private helper functions defined in keys.go for performing these operations,
// thereby abstracting the underlying cryptographic mechanisms from the rest of the application.
//
// Features and Responsibilities:
// - Manages Libp2p and ECDSA cryptographic keys, including private and public keys.
// - Supports loading private keys from environment variables or files, with a fallback
//   to generating new keys if necessary.
// - Provides functionality to marshal and unmarshal private keys, as well as to convert
//   keys between different formats (e.g., from Libp2p to ECDSA format, and vice versa).
// - Offers methods to get hex-encoded representations of keys and to convert public keys
//   to Ethereum address format.
// - Ensures thread-safe initialization and access to the cryptographic keys through the
//   use of the sync.Once mechanism.
//
// Usage:
// To access the KeyManager and its functionalities, use the KeyManagerInstance() function
// which returns the singleton instance of KeyManager. This instance can then be used to
// perform various key management tasks, such as retrieving the application's cryptographic
// keys, converting key formats, and more.
// Example:
//     keyManager := crypto.KeyManagerInstance()

var (
	keyManagerInstance *KeyManager
	once               sync.Once
)

// KeyManager holds all the cryptographic entities used in the application.
type KeyManager struct {
	Libp2pPrivKey crypto.PrivKey    // Libp2p private key
	Libp2pPubKey  crypto.PubKey     // Libp2p public key
	EcdsaPrivKey  *ecdsa.PrivateKey // ECDSA private key
	EcdsaPubKey   *ecdsa.PublicKey  // ECDSA public key
	HexPrivKey    string            // Hex-encoded private key
	HexPubKey     string            // Hex-encoded public key
	EthAddress    string            // Ethereum format address
}

// KeyManagerInstance returns the singleton instance of KeyManager, initializing it if necessary.
func KeyManagerInstance() *KeyManager {
	once.Do(func() {
		keyManagerInstance = &KeyManager{}
		if err := keyManagerInstance.loadPrivateKey(); err != nil {
			logrus.Fatal("Failed to initialize keys:", err)
		}
	})
	return keyManagerInstance
}

func (km *KeyManager) loadPrivateKey() (err error) {
	var keyFile string
	cfg := config.GetInstance()
	if cfg.PrivateKey != "" {
		km.Libp2pPrivKey, err = getPrivateKeyFromEnv(cfg.PrivateKey)
	} else {
		keyFile = config.GetInstance().PrivateKeyFile
		// Check if the private key file exists
		km.Libp2pPrivKey, err = getPrivateKeyFromFile(keyFile)
		if err != nil {
			km.Libp2pPrivKey, err = generateNewPrivateKey(keyFile)
			if err != nil {
				return err
			}
		}
	}
	km.Libp2pPubKey = km.Libp2pPrivKey.GetPublic()
	// After obtaining the libp2p privKey, convert it to an ECDSA private key
	km.EcdsaPrivKey, err = libp2pPrivateKeyToEcdsa(km.Libp2pPrivKey)
	if err != nil {
		return err
	}
	// TODO: do we still need the ecdsa version saved to a file?
	err = saveEcdesaPrivateKeyToFile(km.EcdsaPrivKey, fmt.Sprintf("%s.ecdsa", keyFile))
	if err != nil {
		return err
	}
	km.HexPrivKey, err = getHexEncodedPrivateKey(km.Libp2pPrivKey)
	if err != nil {
		return err
	}
	km.EcdsaPubKey = &km.EcdsaPrivKey.PublicKey
	km.HexPubKey, err = getHexEncodedPublicKey(km.Libp2pPubKey)
	if err != nil {
		return err
	}
	km.EthAddress = ethCrypto.PubkeyToAddress(km.EcdsaPrivKey.PublicKey).Hex()
	return nil
}
