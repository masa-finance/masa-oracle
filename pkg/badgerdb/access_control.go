package badgerdb

import (
	"encoding/hex"
	"fmt"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Load the allowedPeerID and its public key from Viper configuration
func getAllowedPeerIDAndKey() (string, libp2pCrypto.PubKey, error) {
	allowedPeerID := viper.GetString("ALLOWED_PEER_ID")
	// Assuming the public key is stored in a configuration key "ALLOWED_PEER_PUBKEY"
	allowedPeerPubKeyString := viper.GetString("ALLOWED_PEER_PUBKEY")

	if allowedPeerID == "" || allowedPeerPubKeyString == "" {
		logrus.Warn("Allowed peer ID or public key not found in configuration")
		return "", nil, fmt.Errorf("allowed peer ID or public key not configured")
	}

	allowedPeerPubKeyBytes, err := hex.DecodeString(allowedPeerPubKeyString)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode allowed peer public key")
		return "", nil, err
	}

	allowedPeerPubKey, err := libp2pCrypto.UnmarshalPublicKey(allowedPeerPubKeyBytes)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal allowed peer public key")
		return "", nil, err
	}

	return allowedPeerID, allowedPeerPubKey, nil
}

// CanWrite checks if the given host is allowed to write to the database
// Now also requires data and signature for verification
func CanWrite(h host.Host, data []byte, signature []byte) bool {
	allowedPeerID, allowedPeerPubKey, err := getAllowedPeerIDAndKey()
	if err != nil {
		logrus.WithError(err).Error("Failed to load allowed peer ID and public key")
		return false
	}

	if h.ID().String() != allowedPeerID {
		logrus.WithFields(logrus.Fields{
			"hostID":        h.ID().String(),
			"allowedPeerID": allowedPeerID,
		}).Warn("Host ID does not match allowed peer ID")
		return false
	}

	// Verify the signature
	isValid, err := consensus.VerifySignature(allowedPeerPubKey, data, signature)
	if err != nil || !isValid {
		logrus.WithFields(logrus.Fields{
			"hostID":        h.ID().String(),
			"allowedPeerID": allowedPeerID,
			"error":         err,
		}).Warn("Failed to verify signature or signature is invalid")
		return false
	}

	logrus.WithFields(logrus.Fields{
		"hostID":        h.ID().String(),
		"allowedPeerID": allowedPeerID,
	}).Info("Host is allowed to write to the database")
	return true
}
