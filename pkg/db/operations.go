package db

import (
	"context"
	"fmt"
	"time"

	masa "github.com/masa-finance/masa-oracle/pkg"

	"github.com/sirupsen/logrus"
)

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var err error
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var err error
	var val []byte

	if key != node.Host.ID().String() {
		val, err = node.DHT.GetValue(ctx, "/db/"+key)
		cached, er := GetCache(ctx, key)
		if er != nil {
			logrus.Errorf("%v", er)
		}
		logrus.Info("cached", string(cached))
	} else {
		val, err = node.DHT.GetValue(ctx, "/db/"+node.Host.ID().String())
		cached, er := GetCache(ctx, node.Host.ID().String())
		if er != nil {
			logrus.Errorf("%v", er)
		}
		logrus.Info("cached", string(cached))
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to read from the database")
		return nil
	}

	return val
}
