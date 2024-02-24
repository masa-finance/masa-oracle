package main

import (
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
	"github.com/masa-finance/masa-oracle/pkg/badgerdb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SetupDatabasePath checks and prepares the database path
func SetupDatabasePath() string {
	dbPath := viper.GetString("DB_PATH") // Make sure you have this configured
	if dbPath == "" {
		dbPath = filepath.Join(viper.GetString("MASA_DIR"), "masa-node-db") // Adjust the viper key as necessary
		// Check if the directory exists, create it if not
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			if err := os.MkdirAll(dbPath, 0755); err != nil {
				logrus.Fatalf("Failed to create database directory: %v", err)
			} else {
				logrus.Infof("Database directory created at: %v", dbPath)
			}
		}
	}
	return dbPath
}

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
