// Developer Note: This code is part of an ongoing development effort. It is crucial for managing access control within a BadgerDB database environment.
// Specifically, it involves verifying if a peer, identified through its host ID, has the authorization to write data to the database. This process is
// accomplished by retrieving an allowed peer ID and its corresponding public key from the application's configuration settings. Subsequently, a signature
// provided by the peer is verified to ensure it matches the allowed credentials. This mechanism is vital for maintaining the database's security by
// ensuring that only authorized peers can perform write operations.

package badgerdb

import (
	"encoding/hex"
	"fmt"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
)

// AuthorizedNodes Set of authorized nodes that can write to the database
type AuthorizedNodes map[peer.ID]bool

// Tmp
var authorizedNodes = AuthorizedNodes{
	"node1": true,
	"node2": true,
}

// checks if a node is authorized to write to the database
func isAuthorized(nodeID peer.ID) bool {
	// authorizedNodes := viper.GetStringMapStringSlice("ALLOWED_PEER_IDS")
	// Check if the node is in the list of authorized nodes
	for id := range authorizedNodes {
		if id == nodeID {
			return true
		}
	}
	return false
}

// Load the allowedPeerID and its public key from Viper configuration
func getAllowedPeerIDAndKey() (string, libp2pCrypto.PubKey, error) {
	allowedPeerID := config.GetInstance().AllowedPeerId
	// Assuming the public key is stored in a configuration key "ALLOWED_PEER_PUBKEY"
	allowedPeerPubKeyString := config.GetInstance().AllowedPeerPublicKey

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

	// Convert the signature byte slice into a hexadecimal string
	signatureHex := hex.EncodeToString(signature)

	// Verify the signature using the hexadecimal string
	isValid, err := consensus.VerifySignature(allowedPeerPubKey, data, signatureHex)
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
