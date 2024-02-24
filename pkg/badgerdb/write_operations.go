package badgerdb

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/sirupsen/logrus"
)

// WriteData encapsulates the logic for writing data to the database,
// including access control checks from access_control.go.
// It now requires the data (key + value) and signature for verification.
func WriteData(db *badger.DB, key, value, signature []byte, h host.Host) error {
	// Combine key and value as the data to be verified
	data := append(key, value...)

	if !CanWrite(h, data, signature) {
		logrus.WithFields(logrus.Fields{
			"nodeID": h.ID().String(),
		}).Error("Node is not authorized to write to the database")
		return errors.New("unauthorized write attempt")
	}

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to write to the database")
		return err
	}

	return nil
}
