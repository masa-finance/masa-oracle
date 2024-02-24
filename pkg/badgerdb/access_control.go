package badgerdb

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Load the allowedPeerID from Viper configuration
func getAllowedPeerID() string {
	// Assuming the configuration key for the allowed peer ID is "ALLOWED_PEER_ID"
	// with a default value of "defaultPeerID" if not specified in the configuration.
	allowedPeerID := viper.GetString("ALLOWED_PEER_ID")
	if allowedPeerID == "" {
		allowedPeerID = "defaultPeerID"
		logrus.WithFields(logrus.Fields{
			"key":      "ALLOWED_PEER_ID",
			"fallback": allowedPeerID,
		}).Warn("Allowed peer ID not found in configuration, using fallback")
	} else {
		logrus.WithFields(logrus.Fields{
			"key":   "ALLOWED_PEER_ID",
			"value": allowedPeerID,
		}).Info("Allowed peer ID loaded from configuration")
	}
	return allowedPeerID
}

// CanWrite checks if the given host is allowed to write to the database
func CanWrite(h host.Host) bool {
	allowedPeerID := getAllowedPeerID()
	isAllowed := h.ID().String() == allowedPeerID
	if !isAllowed {
		logrus.WithFields(logrus.Fields{
			"hostID":        h.ID().String(),
			"allowedPeerID": allowedPeerID,
		}).Warn("Host is not allowed to write to the database")
	} else {
		logrus.WithFields(logrus.Fields{
			"hostID":        h.ID().String(),
			"allowedPeerID": allowedPeerID,
		}).Info("Host is allowed to write to the database")
	}
	return isAllowed
}
