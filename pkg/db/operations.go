package db

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
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
	Event    json.RawMessage `json:"event"`
}

// ConnectToPostgres Function to connect to Postgres
func ConnectToPostgres(migrations bool) (*sql.DB, error) {
	writer := os.Getenv("WRITER_NODE")
	if writer == "" || writer == "false" {
		return nil, fmt.Errorf("401, node is not authorized to write to the datastore")
	}

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

func FireData(uid string, payload []byte, response []byte) error {
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
	payloadEventJSON := json.RawMessage(value) // `{"event":"Actor Started"}`
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
		return nil, err
	}
	defer database.Close()

	data := []Work{}
	query := `SELECT "work".id, "work".uuid, "work".payload, "work".response, event.payload as "event" FROM "work" INNER JOIN event ON "work".uuid = event.work_id WHERE "work"."uuid" = $1;`
	rows, err := database.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var work Work
		if err = rows.Scan(&work.ID, &work.Uuid, &work.Payload, &work.Response, &work.Event); err != nil {
			return nil, err
		}
		data = append(data, work)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func SendToS3(uid string, payload map[string]string) error {

	// apiUrl := os.Getenv("API_URL")
	apiURL := "https://test.oracle-api.masa.ai/data"
	authToken := "your-secret-token"

	// Create the JSON payload
	// payload := map[string]string{
	// 	"key1": "value1",
	// 	"key2": "value2",
	// }
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
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
