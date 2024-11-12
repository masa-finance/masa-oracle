package masacrypto

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/crypto"
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

// NewKeyManager returns an initialized KeyManager. It first checks for a
// private key set via the PrivateKey config. If not found, it tries to
// load the key from the PrivateKeyFile. As a last resort, it generates
// a new key and saves it to the private key file.
// The private key is loaded into both Libp2p and ECDSA formats for use by
// different parts of the system. The public key and hex-encoded key representations
// are also derived.
func NewKeyManager(privateKey string, privateKeyFile string) (*KeyManager, error) {
	km := &KeyManager{}

	var err error

	if len(privateKey) > 0 {
		km.Libp2pPrivKey, err = getPrivateKeyFromEnv(privateKey)
		if err != nil {
			return nil, err
		}
	} else {
		// Check if the private key file exists
		km.Libp2pPrivKey, err = getPrivateKeyFromFile(privateKeyFile)
		if err != nil {
			km.Libp2pPrivKey, err = generateNewPrivateKey(privateKeyFile)
			if err != nil {
				return nil, err
			}
		}
	}

	km.Libp2pPubKey = km.Libp2pPrivKey.GetPublic()

	// After obtaining the libp2p privKey, convert it to an ECDSA private key
	km.EcdsaPrivKey, err = libp2pPrivateKeyToEcdsa(km.Libp2pPrivKey)
	if err != nil {
		return nil, err
	}

	err = saveEcdesaPrivateKeyToFile(km.EcdsaPrivKey, fmt.Sprintf("%s.ecdsa", privateKeyFile))
	if err != nil {
		return nil, err
	}

	km.HexPrivKey, err = getHexEncodedPrivateKey(km.Libp2pPrivKey)
	if err != nil {
		return nil, err
	}

	km.EcdsaPubKey = &km.EcdsaPrivKey.PublicKey
	km.HexPubKey, err = getHexEncodedPublicKey(km.Libp2pPubKey)
	if err != nil {
		return nil, err
	}

	km.EthAddress = ethCrypto.PubkeyToAddress(km.EcdsaPrivKey.PublicKey).Hex()
	return km, nil
}
