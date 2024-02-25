package util

import (
	"encoding/hex"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	masa "github.com/masa-finance/masa-oracle/pkg"
	masaCrypto "github.com/masa-finance/masa-oracle/pkg/crypto"
	"github.com/spf13/viper"
)

// This function constructs the file path for the private key using the directory and file name keys defined in init.go.
func GetKeyFilePath() string {
	return filepath.Join(viper.GetString(masa.MasaDir), viper.GetString(masa.PrivKeyFile))
}

// LoadPrivateKey loads or generates the private key based on the configuration.
func LoadPrivateKey() (crypto.PrivKey, error) {
	keyFilePath := GetKeyFilePath()
	privKey, _, _, err := masaCrypto.GetOrCreatePrivateKey(keyFilePath)
	return privKey, err
}

// GetHexEncodedPublicKey derives the public key from the given private key and returns its hex-encoded string.
func GetHexEncodedPublicKey(privKey crypto.PrivKey) (string, error) {
	pubKey := privKey.GetPublic()
	pubKeyBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pubKeyBytes), nil
}

// GetPeerIDFromPubKey converts a public key to a peer ID.
func GetPeerIDFromPubKey(pubKey crypto.PubKey) (string, error) {
	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return peerID.String(), nil
}
