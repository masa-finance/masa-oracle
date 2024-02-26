// storage.go
package badgerdb

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"
)

// InitializeStorage creates and returns a new BadgerDB storage instance
func InitializeStorage(storagePath string) (*badger.DB, error) {
	options := badger.DefaultOptions(storagePath).WithLoggingLevel(badger.WARNING)
	storage, err := badger.Open(options)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize BadgerDB storage")
		return nil, err
	}
	return storage, nil
}
