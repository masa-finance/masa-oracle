package main

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/masa-finance/masa-oracle/pkg/badgerdb" // Ensure this import path matches your project structure
	"github.com/sirupsen/logrus"
)

// InitializeDB initializes and returns a BadgerDB instance
func InitializeDB(dbPath string) (*badger.DB, error) {
	db, err := badgerdb.InitializeStorage(dbPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize BadgerDB in masa-node")
		return nil, err
	}
	logrus.Info("BadgerDB initialized successfully in masa-node")
	return db, nil
}
