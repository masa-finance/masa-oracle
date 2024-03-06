package badgerdb

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
)

func TestInitializeStorage(t *testing.T) {

	// Test with valid path
	storage, err := InitializeStorage("valid/path")
	if err != nil {
		t.Fatal("Unexpected error initializing storage with valid path:", err)
	}
	if storage == nil {
		t.Fatal("Expected non-nil storage instance")
	}

	// Test with invalid path
	storage, err = InitializeStorage("invalid/path")
	if err == nil {
		t.Fatal("Expected error initializing storage with invalid path")
	}
	if storage != nil {
		t.Fatal("Expected nil storage instance")
	}

	// Test storage options
	opts := badger.DefaultOptions("/tmp/badger").WithLoggingLevel(badger.WARNING)
	storage, err = badger.Open(opts)
	if err != nil {
		t.Fatal("Unexpected error creating storage with options:", err)
	}
	if storage.Opts().Logger == nil {
		t.Fatal("Expected non-nil logger in storage options")
	}
}

func TestSetupDatabasePath(t *testing.T) {
	setupDatabasePath()
}
