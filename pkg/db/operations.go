// TODO Rename this to something else, this is NOT a database (it just stores data in the DHT and in a cache)
package db

import (
	"context"
	"encoding/json"
	"fmt"
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
	err = node.DHT.PutValue(ctx, "/db/"+key, value) // any key value so the data is public

	_, er := PutCache(ctx, key, value)
	if er != nil {
		logrus.Errorf("[-] Error putting cache: %v", er)
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

	val, err = GetCache(ctx, key)
	if val == nil || err != nil {
		val, err = node.DHT.GetValue(ctx, "/db/"+key)
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
