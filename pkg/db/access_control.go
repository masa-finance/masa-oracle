// Developer Note: This code is part of an ongoing development effort. It is crucial for managing access control within a DHT database environment.
// Specifically, it involves verifying if a peer, identified through its host ID, has the authorization to write data to the database. This process is
// accomplished by retrieving an allowed peer ID and its corresponding public key from the application's configuration settings. Subsequently, a signature
// provided by the peer is verified to ensure it matches the allowed credentials. This mechanism is vital for maintaining the database's security by
// ensuring that only authorized peers can perform write operations.

package db

import (
	"encoding/hex"
	"strings"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
)

// AuthorizedNodes Set of authorized nodes that can write to the database
type AuthorizedNodes map[string]bool

// authorizedNodes Interface
var authorizedNodes = AuthorizedNodes{}

// Modifier idiomatic check if a node is authorized to write to the database
func isAuthorized(nodeID string) bool {
	// Check if the node is in the list of authorized nodes
	for id := range authorizedNodes {
		if id == nodeID {
			return true
		}
	}
	return false
}

// Verifier checks if the given host is allowed to access to the database and verifies the signature
func Verifier(h host.Host, data []byte, signature []byte) bool {
	// Load configuration instance
	cfg := config.GetInstance()

	// Get allowed peer ID and public key from the configuration
	allowedPeerID := cfg.AllowedPeerId
	allowedPeerPubKeyString := cfg.AllowedPeerPublicKey

	if allowedPeerID == "" || allowedPeerPubKeyString == "" {
		logrus.Warn("Allowed peer ID or public key not found in configuration")
		return false
	}

	// Decode the public key
	allowedPeerPubKeyBytes, err := hex.DecodeString(allowedPeerPubKeyString)
	if err != nil {
		logrus.WithError(err).Error("Failed to decode allowed peer public key")
		return false
	}

	// Unmarshal the public key
	allowedPeerPubKey, err := libp2pCrypto.UnmarshalPublicKey(allowedPeerPubKeyBytes)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal allowed peer public key")
		return false
	}

	// Check if the host ID matches the allowed peer ID
	if h.ID().String() != allowedPeerID {
		logrus.WithFields(logrus.Fields{
			"hostID":        h.ID().String(),
			"allowedPeerID": allowedPeerID,
		}).Warn("Host ID does not match allowed peer ID")
		return false
	}

	// Verify the signature
	isValid, err := consensus.VerifySignature(allowedPeerPubKey, data, hex.EncodeToString(signature))
	if err != nil || !isValid {
		logrus.WithFields(logrus.Fields{
			"hostID":        h.ID().String(),
			"allowedPeerID": allowedPeerID,
			"error":         err,
		}).Warn("Failed to verify signature or signature is invalid")
		return false
	}

	if strings.ToLower(cfg.Validator) == "true" {

		logrus.WithFields(logrus.Fields{
			"hostID":        h.ID().String(),
			"allowedPeerID": allowedPeerID,
		}).Info("Host is allowed to write to the database")

		authorizedNodes = map[string]bool{
			allowedPeerID: true,
		}
	}

	return true
}
