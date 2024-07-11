package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	masa "github.com/masa-finance/masa-oracle/pkg"
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
func ReadData(node *masa.OracleNode, key string) ([]byte, error) {
	logrus.WithFields(logrus.Fields{
		"nodeID":       node.Host.ID().String(),
		"isAuthorized": true,
		"ReadData":     true,
	}).Info("Attempting to read data")

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
		}).Error("Failed to read from the database")
		return nil, err
	}

	return val, nil
}
