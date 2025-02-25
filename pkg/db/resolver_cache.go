package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"

	"github.com/masa-finance/masa-oracle/node"

	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/sirupsen/logrus"
)

var cache ds.Datastore
var nodeDataChan = make(chan *pubsub.NodeData)

type Record struct {
	Key   string
	Value []byte
}

// InitResolverCache initializes the resolver cache for the Masa Oracle node.
//
// Parameters:
//   - node: A pointer to the Masa Oracle node (masa.OracleNode) that the resolver cache will be associated with.
//   - keyManager: A pointer to the key manager (masacrypto.KeyManager) used for cryptographic operations.
//
// This function sets up the resolver cache for the Masa Oracle node. The resolver cache is responsible for storing and managing resolved data within the node.
//
// The function takes two parameters:
//  1. `node`: It represents the Masa Oracle node instance to which the resolver cache will be attached. The node provides the necessary context and dependencies for the resolver cache to operate.
//  2. `keyManager`: It is an instance of the key manager that handles cryptographic operations. The key manager is used by the resolver cache for any required cryptographic tasks, such as signing or verifying data.
//
// The purpose of this function is to initialize the resolver cache and perform any necessary setup or configuration. It associates the resolver cache with the provided Masa Oracle node and key manager.
//
// Note: The specific implementation details of the `InitResolverCache` function are not provided in the given code snippet. The function signature suggests that it initializes the resolver cache, but the actual initialization logic would be present in the function body.
func InitResolverCache(node *node.OracleNode, keyManager *masacrypto.KeyManager, allowedPeerID string, allowedPeerPubKeyString string, isValidator bool) {
	var err error

	cache, err = leveldb.NewDatastore(node.Options.CachePath, nil)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("[+] ResolverCache initialized")

	data := []byte(node.Host.ID().String())
	signature, err := consensus.SignData(keyManager.Libp2pPrivKey, data)
	if err != nil {
		logrus.Errorf("[-] Error signing data: %v", err)
	}
	_ = Verifier(node.Host, data, signature, allowedPeerID, allowedPeerPubKeyString, isValidator)

	go monitorNodeData(context.Background(), node)

	if !isAuthorized(node.Host.ID().String()) {
		logrus.WithFields(logrus.Fields{
			"nodeID":       node.Host.ID().String(),
			"isAuthorized": false,
			"Sync":         true,
		})
		return
	} else {
		syncInterval := time.Second * 60 // Change as needed
		go sync(context.Background(), node, syncInterval)
	}
}

// PutCache puts a key-value pair into the resolver cache.
//
// It takes a context, a key as a string, and a value as a byte slice.
// It converts the key into a datastore key and puts the key-value pair
// into the cache.
//
// It returns the original key string and a possible error.
func PutCache(ctx context.Context, keyStr string, value []byte) (any, error) {
	err := cache.Put(ctx, ds.NewKey(keyStr), value)
	if err != nil {
		return nil, err
	}
	return keyStr, nil
}

// GetCache retrieves a value from the resolver cache for the given key.
// It takes a context and a key string, converts the key into a datastore key,
// gets the value from the cache, and returns the value byte slice and a possible error.
func GetCache(ctx context.Context, keyStr string) ([]byte, error) {
	value, err := cache.Get(ctx, ds.NewKey(keyStr))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// DelCache deletes a key-value pair from the resolver cache.
// It takes a context and a key string, converts the key into a datastore key,
// and deletes the key-value pair from the cache.
// It returns a bool indicating if the deletion succeeded.
func DelCache(ctx context.Context, keyStr string) bool {
	key := ds.NewKey(keyStr)
	err := cache.Delete(ctx, key)
	if err != nil {
		return false
	} else {
		return true
	}
}

// UpdateCache updates the value for the given key in the resolver cache.
// It first checks if the key already exists using Has().
// If it doesn't exist, it returns an error.
// If the key does exist, it puts the new value into the cache using Put().
// It returns a bool indicating if the update succeeded, and a possible error.
func UpdateCache(ctx context.Context, keyStr string, newValue []byte) (bool, error) {
	// Check if the key exists
	key := ds.NewKey(keyStr)
	res, err := cache.Has(ctx, key)
	if err != nil {
		return false, fmt.Errorf("error checking key existence: %w", err)
	}

	if !res {
		return false, fmt.Errorf("key does not exist, adding new key-value pair, %+v", err)
	}

	// Put the new value for the key in the datastore
	if err := cache.Put(ctx, key, newValue); err != nil {
		return false, fmt.Errorf("error updating data: %w", err)
	}

	return true, nil
}

// QueryAll queries the resolver cache for all records and returns them as a slice of Record structs.
// It executes a query.Query{} to get all results, closes the results when done, iterates through
// the results, appending each record to a slice, and returns the slice.
func QueryAll(ctx context.Context) ([]Record, error) {
	results, err := cache.Query(ctx, query.Query{})
	if err != nil {
		logrus.Errorf("[-] Failed to query the resolver cache: %v", err)
		return nil, err
	}
	defer results.Close()

	var records []Record

	for result := range results.Next() {
		if result.Error != nil {
			logrus.Errorf("[-] Error iterating query results: %v", result.Error)
			return nil, result.Error
		}
		// Append the record to the slice
		records = append(records, Record{Key: result.Entry.Key, Value: result.Entry.Value})
	}

	return records, nil
}

// sync periodically calls iterateAndPublish to synchronize the node's state with
// the dht on the provided interval. It runs this in a loop, exiting
// when the context is canceled.
func sync(ctx context.Context, node *node.OracleNode, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			iterateAndPublish(ctx, node)
		case <-ctx.Done():
			return
		}
	}
}

// iterateAndPublish synchronizes the node's local cache with the dht
// by querying all records, and publishing each one to the dht. It
// logs any errors encountered. This allows periodically syncing the node's
// cached data with the latest dht state.
func iterateAndPublish(ctx context.Context, node *node.OracleNode) {
	records, err := QueryAll(ctx)
	if err != nil {
		logrus.Errorf("[-] Error querying all records: %+v", err)
		return
	}

	if node.DHT == nil {
		logrus.Errorf("DHT instance is nil. Skipping iterateAndPublish operation.")
		return
	}

	for _, record := range records {
		key := record.Key
		if len(key) > 0 && key[0] == '/' {
			key = key[1:]
		}
		logrus.Printf("syncing %s", key)
		_ = WriteData(node, key, record.Value)

		// sync blocks
		blocks, err := node.DHT.GetValue(ctx, "/db/blocks")
		if err != nil {
			logrus.Debugf("[-] Error getting block data: %v", err)
			continue
		}

		if err := WriteData(node, "blocks", blocks); err != nil {
			logrus.Debugf("[-] Error writing block data: %v", err)
		}

		// sync ipfs
		ipfs, e := node.DHT.GetValue(ctx, "/db/ipfs")
		if e != nil {
			logrus.Debugf("[-] Error unmarshalling IPFS data: %v", e)
		} else {
			_ = WriteData(node, "ipfs", ipfs)
		}
	}
}

// monitorNodeData periodically publishes the local node's status to the
// dht, and syncs node status data published by other nodes.
// It runs a ticker to call iterateAndPublish on the provided interval.
func monitorNodeData(ctx context.Context, node *node.OracleNode) {
	syncInterval := time.Second * 60

	err := node.Subscribe(config.NodeGossipTopic, node.NodeTracker)
	if err != nil {
		logrus.Errorf("[-] Error subscribing to node gossip topic: %v", err)
	}

	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// if dht does not have data check cache and update if exists
			var nodeData *pubsub.NodeData
			nodeDataBytes, err := GetCache(context.Background(), node.Host.ID().String())
			if err != nil {
				nodeData = node.NodeTracker.GetNodeData(node.Host.ID().String())
			} else {
				err = json.Unmarshal(nodeDataBytes, &nodeData)
				if err != nil {
					logrus.Errorf("[-] Error unmarshalling node data: %v", err)
				}
			}

			jsonData, _ := json.Marshal(nodeData)
			e := node.PublishTopic(config.NodeGossipTopic, jsonData)
			if e != nil {
				logrus.Errorf("[-] Error publishing node data: %v", e)
			}
		case nodeData := <-nodeDataChan:
			jsonData, _ := json.Marshal(nodeData)
			_ = WriteData(node, nodeData.PeerId.String(), jsonData)
		case <-ctx.Done():
			return
		}
	}
}
