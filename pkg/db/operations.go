package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/sirupsen/logrus"
)

type WorkEvent struct {
	CID       string          `json:"cid"`
	PeerId    string          `json:"peer_id"`
	Payload   json.RawMessage `json:"payload"`
	Duration  float64         `json:"duration"`
	Timestamp int64           `json:"timestamp"`
}

// WriteData encapsulates the logic for writing data to the database,
// including access control checks from access_control.go.
func WriteData(node *node.OracleNode, key string, value []byte) error {
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
			logrus.Errorf("[-] Error putting cache: %v", er)
		}
	} else {
		err = node.DHT.PutValue(ctx, "/db/"+node.Host.ID().String(), value) // nodes private data based on node id

		_, er := PutCache(ctx, node.Host.ID().String(), value)
		if er != nil {
			logrus.Errorf("[-] Error putting cache: %v", er)
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
func ReadData(node *node.OracleNode, key string) ([]byte, error) {
	logrus.WithFields(logrus.Fields{
		"nodeID":       node.Host.ID().String(),
		"isAuthorized": true,
		"ReadData":     true,
	}).Info("[+] Attempting to read data")

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
			// we don't need to check for err since !exists gives an err also - we only need to know if the record exists or not in this context
		}).Debug("[-] Failed to read from the database")
		return nil, err
	}

	return val, nil
}

// SendToS3 sends a payload to an S3-compatible API.
//
// Parameters:
//   - uid: The unique identifier for the payload.
//   - payload: The payload to be sent, represented as a map of key-value pairs.
//
// Returns:
//   - error: Returns an error if the operation fails, otherwise returns nil.
func SendToS3(uid string, payload map[string]string) error {

	apiURL := os.Getenv("API_URL")
	authToken := "your-secret-token"

	// Creating the JSON payload
	// payload := map[string]string{
	// 	"key1": "value1",
	// 	"key2": "value2",
	// }

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[-] Failed to marshal JSON payload: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("[-] Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("[-] Failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("[-] Received non-OK response: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}
