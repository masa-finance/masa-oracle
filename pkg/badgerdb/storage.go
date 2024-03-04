package badgerdb

import (
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

// InitializeDB initializes and returns a BadgerDB instance
func InitializeDB() (*badger.DB, error) {
	setupDatabasePath()
	db, err := InitializeStorage(config.GetInstance().DbPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize BadgerDB in masa-node")
		return nil, err
	}
	logrus.Info("BadgerDB initialized successfully in masa-node")
	return db, nil
}

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

// setupDatabasePath checks and prepares the database path
func setupDatabasePath() {
	cfg := config.GetInstance()
	if _, err := os.Stat(cfg.DbPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.DbPath, 0755); err != nil {
			logrus.Fatalf("Failed to create database directory: %v", err)
		} else {
			logrus.Infof("Database directory created at: %v", cfg.DbPath)
		}
	}
}
