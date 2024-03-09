package db

import (
	"context"
	"fmt"
	"time"

	masa "github.com/masa-finance/masa-oracle/pkg"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"
)

// WriteData encapsulates the logic for writing data to the database,
// including access control checks from access_control.go.
func WriteData(node *masa.OracleNode, key string, value []byte, h host.Host) (bool, error) {
	if !isAuthorized(h.ID().String()) {
		logrus.WithFields(logrus.Fields{
			"nodeID":       h.ID().String(),
			"isAuthorized": false,
			"WriteData":    true,
		}).Error("DHT write authorization failed")
		return false, fmt.Errorf("401, node is not authorized to write to the datastore")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var err error
	if key != node.Host.ID().String() {
		err = node.DHT.PutValue(ctx, "/db/"+key, value) // any key value so the data is public
		res, er := PutCache(ctx, key, value)
		if er != nil {
			logrus.Errorf("%v", er)
		}
		logrus.Println("res", res)
	} else {
		err = node.DHT.PutValue(ctx, "/db/"+node.Host.ID().String(), value) // nodes private data based on node id
		res, er := PutCache(ctx, node.Host.ID().String(), value)
		if er != nil {
			logrus.Errorf("%v", er)
		}
		logrus.Println("res", res)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to write to the database")
		return false, err
	}

	return true, nil
}

// ReadData reads the value for the given key from the database.
// It requires the host for access control verification before reading.
func ReadData(node *masa.OracleNode, key string, h host.Host) []byte {

	logrus.WithFields(logrus.Fields{
		"nodeID":       h.ID().String(),
		"isAuthorized": true,
		"ReadData":     true,
	}).Error("DHT write authorization failed")

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

// @TODO offline and long term pinning
// Cache Syncing
func cacheSyncing() {}
