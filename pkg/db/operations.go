package db

import (
	"context"
	"fmt"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"
)

// WriteData encapsulates the logic for writing data to the database,
// including access control checks from access_control.go.
func WriteData(node *masa.OracleNode, key string, value []byte, h host.Host) (bool, error) {
	if !isAuthorized(h.ID().String()) {
		logrus.WithFields(logrus.Fields{
			"nodeID": h.ID().String(),
		}).Error("node is not authorized to write to the datastore")
		return false, fmt.Errorf("401, node is not authorized to write to the datastore")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var err error
	if key != node.Host.ID().String() {
		err = node.DHT.PutValue(ctx, "/db/"+string(key), value)
	} else {
		err = node.DHT.PutValue(ctx, "/db/"+node.Host.ID().String(), value)
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

	// @TODO verify we want any node to read if not we can implement here
	//if isAuthorized(h.ID().String()) {
	//	logrus.WithFields(logrus.Fields{
	//		"nodeID": h.ID().String(),
	//	}).Error("node is not authorized to read the datastore")
	//	return []byte("401, node is not authorized to read from the datastore")
	//}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
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
		}).Error("Failed to read from the database")
		return nil
	}

	return val
}
