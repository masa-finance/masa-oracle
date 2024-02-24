package badgerdb

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"
)

// DataEntry represents the structure of the data to be stored in the database.
type DataEntry struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Nonce     int64  `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
}

// WriteDataEntry writes a DataEntry to the BadgerDB if the host has the right to write.
// It now requires the host, data, and signature for access control verification.
func WriteDataEntry(db *badger.DB, h host.Host, data []byte, signature []byte, entry DataEntry) error {
	// Perform access control check
	if !CanWrite(h, data, signature) {
		logrus.WithFields(logrus.Fields{
			"hostID": h.ID().String(),
			"key":    entry.Key,
		}).Warn("Access denied for writing data entry to the database")
		return fmt.Errorf("access denied for host %s", h.ID().String())
	}

	entryKey := []byte(entry.Key)
	entryValue := []byte(entry.Value)

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(entryKey, entryValue)
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"key":   entry.Key,
			"error": err,
		}).Error("Failed to write data entry to the database")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"key":       entry.Key,
		"timestamp": entry.Timestamp,
	}).Info("Successfully wrote data entry to the database")
	return nil
}
