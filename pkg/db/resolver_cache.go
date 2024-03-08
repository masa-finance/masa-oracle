package db

import (
	"context"
	"fmt"
	"log"

	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	leveldb "github.com/ipfs/go-ds-leveldb"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

var cache ds.Datastore

type Record struct {
	Key   string
	Value []byte
}

func InitCacheStore() {
	var err error
	cachePath := "./CACHE"
	cache, err = leveldb.NewDatastore(cachePath, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cache initialized")
}

func putData(ctx context.Context, keyStr string, value []byte) (any, error) {
	key, _ := stringToCid(keyStr)
	err := cache.Put(ctx, ds.NewKey("cache/"+key), value)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func getData(ctx context.Context, key string) []byte {
	value, err := cache.Get(ctx, ds.NewKey("cache/"+key))
	if err != nil {
		log.Fatalf("Failed to get data: %v", err)
	}
	return value
}

func delData(ctx context.Context, keyStr string) bool {
	var err error
	key := ds.NewKey("cache/" + keyStr)
	err = cache.Delete(ctx, key)
	if err != nil {
		return false
	}
	return true
}

func updateData(ctx context.Context, keyStr string, newValue []byte) (bool, error) {
	// Check if the key exists
	key := ds.NewKey("cache/" + keyStr)
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

func queryAllData(ctx context.Context) ([]Record, error) {
	results, err := cache.Query(ctx, query.Query{})
	if err != nil {
		logrus.Errorf("Failed to query the datastore: %v", err)
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
