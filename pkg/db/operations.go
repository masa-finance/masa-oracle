package db

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
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

type WorkEvent struct {
	ID      int64           `json:"id"`
	WorkId  string          `json:"work_id"`
	Payload json.RawMessage `json:"payload"`
}
type Work struct {
	ID       int64           `json:"id"`
	Uuid     string          `json:"uuid"`
	Payload  json.RawMessage `json:"payload"`
	Response json.RawMessage `json:"response"`
}

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

func PostData(uid string, payload []byte, response []byte) error {
	database, err := ConnectToPostgres(false)
	if err != nil {
		return nil
	}
	defer database.Close()
	insertQuery := `INSERT INTO "public"."work" ("uuid", "payload", "response") VALUES ($1, $2, $3)`
	payloadJSON := json.RawMessage(payload)
	responseJSON := json.RawMessage(response)
	_, err = database.Exec(insertQuery, uid, payloadJSON, responseJSON)
	if err != nil {
		return err
	}
	return nil
}

func FireEvent(uid string, value []byte) error {
	database, err := ConnectToPostgres(false)
	if err != nil {
		return err
	}
	defer database.Close()

	insertEventQuery := `INSERT INTO "public"."event" ("work_id", "payload") VALUES ($1, $2)`
	payloadEventJSON := json.RawMessage(`{"event":"Actor Started"}`)
	_, err = database.Exec(insertEventQuery, uid, payloadEventJSON)
	if err != nil {
		return err
	}
	return nil
}

func GetData(uid string) ([]Work, error) {
	database, err := ConnectToPostgres(false)
	if err != nil {
		logrus.Error(err)
	}
	defer database.Close()
	data := []Work{}
	query := `SELECT "id", "payload", "response" FROM "public"."work"`
	rows, err := database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id       int64
		payload  json.RawMessage
		response json.RawMessage
	)

	for rows.Next() {
		if err = rows.Scan(&id, &payload, &response); err != nil {
			return nil, err
		}
		data = append(data, Work{id, uid, payload, response})
	}
	return data, nil
}

// WriteData encapsulates the logic for writing data to the database,
// including access control checks from access_control.go.
func WriteData(node *masa.OracleNode, key string, value []byte) error {
	if !isAuthorized(node.Host.ID().String()) {
		logrus.WithFields(logrus.Fields{
			"nodeID":       node.Host.ID().String(),
			"isAuthorized": false,
			"WriteData":    true,
		})
		return fmt.Errorf("401, node is not authorized to write to the datastore")
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
		return err
	}

	return nil
}

// ReadData reads the value for the given key from the database.
// It requires the host for access control verification before reading.
func ReadData(node *masa.OracleNode, key string) []byte {
	logrus.WithFields(logrus.Fields{
		"nodeID":       node.Host.ID().String(),
		"isAuthorized": true,
		"ReadData":     true,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	var err error
	var val []byte

	if key != node.Host.ID().String() {
		val, err = GetCache(ctx, key)
		if val == nil || err != nil {
			val, err = node.DHT.GetValue(ctx, "/db/"+key)
		}
	} else {
		val, err = GetCache(ctx, node.Host.ID().String())
		if val == nil || err != nil {
			val, err = node.DHT.GetValue(ctx, "/db/"+node.Host.ID().String())
		}
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Debug("Failed to read from the database")
		return val
	}

	return val
}
