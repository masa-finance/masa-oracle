package db

import (
	"context"
	"fmt"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"log"
	"time"

	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/masa-finance/masa-oracle/pkg/config"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

var cache ds.Datastore

type Record struct {
	Key   string
	Value []byte
}

func InitResolverCache(node *masa.OracleNode) {
	var err error
	cachePath := config.GetInstance().CachePath
	cache, err = leveldb.NewDatastore(cachePath, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ResolverCache initialized")

	syncInterval := 10 * time.Second // Change as needed
	sync(context.Background(), node, syncInterval)
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
	for i, record := range records {
		logrus.Printf("%d, %s %s ", i, record.Key, record.Value)
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

// stringToCid option to use
func stringToCid(str string) (string, error) {
	// Create a multihash from the string
	mhHash, err := mh.Sum([]byte(str), mh.SHA2_256, -1)
	if err != nil {
		return "", err
	}

	// Create a CID from the multihash
	cidKey := cid.NewCidV1(cid.Raw, mhHash).String()

	return cidKey, nil
}
