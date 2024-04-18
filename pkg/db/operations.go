package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sort"
	"time"

	masa "github.com/masa-finance/masa-oracle/pkg"

	"github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var migrations embed.FS

// ConnectToPostgres Function to connect to Postgres
func ConnectToPostgres(migrations bool) (*sql.DB, error) {
	database, err := sql.Open("postgres", os.Getenv("PG_URL"))
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	if migrations {
		err = applyMigrations(database)
		if err != nil {
			return nil, err
		}
	}
	return database, nil
}

// applyMigrations Function to execute SQL
func applyMigrations(database *sql.DB) error {
	// Read directory entries from embedded FS
	entries, err := fs.ReadDir(migrations, "migrations")
	if err != nil {
		return err
	}

	// Sort files by name to ensure correct order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if !entry.IsDir() {
			content, err := fs.ReadFile(migrations, "migrations/"+entry.Name())
			if err != nil {
				return err
			}

			// Execute SQL file content
			log.Printf("Applying migration: %s", entry.Name())
			if _, err := database.Exec(string(content)); err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteData encapsulates the logic for writing data to the database,
// including access control checks from access_control.go.
func WriteData(node *masa.OracleNode, key string, value []byte) (bool, error) {
	if !isAuthorized(node.Host.ID().String()) {
		logrus.WithFields(logrus.Fields{
			"nodeID":       node.Host.ID().String(),
			"isAuthorized": false,
			"WriteData":    true,
		})
		return false, fmt.Errorf("401, node is not authorized to write to the datastore")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	var err error
	node.DHT.ForceRefresh()
	if key != node.Host.ID().String() {
		err = node.DHT.PutValue(ctx, "/db/"+key, value) // any key value so the data is public

		_, er := PutCache(ctx, key, value)
		if er != nil {
			logrus.Errorf("%v", er)
		}
	} else {
		err = node.DHT.PutValue(ctx, "/db/"+node.Host.ID().String(), value) // nodes private data based on node id

		_, er := PutCache(ctx, node.Host.ID().String(), value)
		if er != nil {
			logrus.Errorf("%v", er)
		}
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		})
		return false, err
	}

	return true, nil
}

// ReadData reads the value for the given key from the database.
// It requires the host for access control verification before reading.
func ReadData(node *masa.OracleNode, key string) []byte {
	logrus.WithFields(logrus.Fields{
		"nodeID":       node.Host.ID().String(),
		"isAuthorized": true,
		"ReadData":     true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	var err error
	var val []byte

	if key != node.Host.ID().String() {
		val, err = node.DHT.GetValue(ctx, "/db/"+key)
	} else {
		val, err = node.DHT.GetValue(ctx, "/db/"+node.Host.ID().String())
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("Failed to read from the database")
		return nil
	}

	return val
}
