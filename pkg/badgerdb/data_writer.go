package badgerdb

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

// DataEntry represents the structure of the data to be stored in the database.
type DataEntry struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Nonce     int64  `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
}

// WriteDataEntry writes a DataEntry to the BadgerDB.
func WriteDataEntry(db *badger.DB, entry DataEntry) error {
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
