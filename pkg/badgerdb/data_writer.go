package badgerdb

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"
)

// NodeLifecycleDataEntry represents the structure of the data to be stored in the database.
type NodeLifecycleDataEntry struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Nonce     int64  `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
}

// WriteDataEntry writes a NodeLifecycleDataEntry to the BadgerDB if the host has the right to write.
// It now requires the host, data, and signature for access control verification.
func WriteDataEntry(db *badger.DB, h host.Host, signature []byte, entry NodeLifecycleDataEntry) error {
	// Use the key and value from the NodeLifecycleDataEntry
	key := []byte(entry.Key)
	value := []byte(entry.Value)

	// Call WriteData from write_operations.go to perform the write operation
	if err := WriteData(db, key, value, signature, h); err != nil {
		logrus.WithFields(logrus.Fields{
			"key":   entry.Key,
			"error": err,
		}).Error("Failed to write node lifecycle data entry to the database")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"key":       entry.Key,
		"timestamp": entry.Timestamp,
	}).Info("Successfully wrote node lifecycle data entry to the database")
	return nil
}
