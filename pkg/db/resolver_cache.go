package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/nodestatus"

	masa "github.com/masa-finance/masa-oracle/pkg"

	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

var cache ds.Datastore
var nodeStatusCh = make(chan []byte)

type Record struct {
	Key   string
	Value []byte
}

func InitResolverCache(node *masa.OracleNode, keyManager *masacrypto.KeyManager) {
	var err error
	cachePath := config.GetInstance().CachePath
	cache, err = leveldb.NewDatastore(cachePath, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println("ResolverCache initialized")

	data := []byte(node.Host.ID().String())
	signature, err := consensus.SignData(keyManager.Libp2pPrivKey, data)
	if err != nil {
		logrus.Errorf("%v", err)
	}
	_ = Verifier(node.Host, data, signature)

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

func PutCache(ctx context.Context, keyStr string, value []byte) (any, error) {
	// key, _ := stringToCid(keyStr)
	err := cache.Put(ctx, ds.NewKey(keyStr), value)
	if err != nil {
		return nil, err
	}
	return keyStr, nil
}

func GetCache(ctx context.Context, keyStr string) ([]byte, error) {
	value, err := cache.Get(ctx, ds.NewKey(keyStr))
	if err != nil {
		return nil, err
	}
	return value, nil
}

func DelCache(ctx context.Context, keyStr string) bool {
	var err error
	key := ds.NewKey(keyStr)
	err = cache.Delete(ctx, key)
	if err != nil {
		return false
	}
	return true
}

func UpdateCache(ctx context.Context, keyStr string, newValue []byte) (bool, error) {
	// Check if the key exists
	key := ds.NewKey(keyStr)
	res, err := cache.Has(ctx, key)
	if err != nil {
		return false, fmt.Errorf("error checking key existence: %w", err)
	}

	if !res {
		return false, fmt.Errorf("key does not exist, adding new key-value pair, %+v", err.Error())
	}

	// Put the new value for the key in the datastore
	if err := cache.Put(ctx, key, newValue); err != nil {
		return false, fmt.Errorf("error updating data: %w", err)
	}

	return true, nil
}

func sync(ctx context.Context, node *masa.OracleNode, interval time.Duration) {
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

func iterateAndPublish(ctx context.Context, node *masa.OracleNode) {
	records, err := QueryAll(ctx)
	if err != nil {
		logrus.Errorf("%+v", err)
	}
	for _, record := range records {
		logrus.Printf("syncing record %s %s", record.Key, record.Value)
		// ok := DelCache(ctx, record.Key)
		//if ok {
		//	logrus.Println("deleted")
		//}
		_, _ = WriteData(node, record.Key, record.Value)
	}
}

func QueryAll(ctx context.Context) ([]Record, error) {
	results, err := cache.Query(ctx, query.Query{})
	if err != nil {
		logrus.Errorf("Failed to query the resolver cache: %v", err)
		return nil, err
	}
	defer results.Close()

	var records []Record

	for result := range results.Next() {
		if result.Error != nil {
			logrus.Errorf("Error iterating query results: %v", result.Error)
			return nil, result.Error
		}
		// Append the record to the slice
		records = append(records, Record{Key: result.Entry.Key, Value: result.Entry.Value})
	}

	return records, nil
}

func monitorNodeData(ctx context.Context, node *masa.OracleNode) {
	syncInterval := time.Second * 60
	nodeStatusHandler := &nodestatus.SubscriptionHandler{NodeStatusCh: nodeStatusCh}
	err := node.PubSubManager.Subscribe(config.TopicWithVersion(config.NodeStatusTopic), nodeStatusHandler)
	if err != nil {
		logrus.Println(err)
	}

	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
			jsonData, _ := json.Marshal(nodeData)
			e := node.PubSubManager.Publish(config.TopicWithVersion(config.NodeStatusTopic), jsonData)
			if e != nil {
				logrus.Printf("%v", e)
			}

		case <-nodeStatusCh:
			nodes := nodeStatusHandler.NodeStatus
			for _, n := range nodes {
				jsonData, _ := json.Marshal(n)
				_, _ = WriteData(node, "/db/"+n.PeerID, jsonData)
			}
		case <-ctx.Done():
			return
		}
	}
}
