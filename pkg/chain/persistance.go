package chain

import (
	"github.com/dgraph-io/badger"
	"github.com/sirupsen/logrus"
)

const (
	KeyLastHash = "last_hash"
)

type Persistence struct {
	db *badger.DB
}

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
}

func (p *Persistence) Init(path string, genesisFn func() (Serializable, []byte)) ([]byte, error) {
	dbOptions := badger.DefaultOptions(path)
	dbOptions.Logger = nil
	db, err := badger.Open(dbOptions)
	if err != nil {
		logrus.Error("[-] Failed to initialize datastore: ", db, err)
		return nil, err
	}
	p.db = db

	var lastHash []byte

	err = p.db.Update(func(transaction *badger.Txn) error {
		lastHashKey := KeyLastHash
		item, err := transaction.Get([]byte(lastHashKey))

		if err != nil {
			if err == badger.ErrKeyNotFound {
				logrus.Warn("Creating genesis transaction...")
				genesisBlock, genesisHash := genesisFn()
				serialData, err := genesisBlock.Serialize()
				if err != nil {
					return err
				}
				if err := transaction.Set(genesisHash, serialData); err != nil {
					return err
				}
				if err := transaction.Set([]byte(KeyLastHash), genesisHash); err != nil {
					return err
				}
				lastHash = genesisHash
				return nil
			}
			return err
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			logrus.Error("[-] Error occured while getting item value by key: ", KeyLastHash, item, err)
			return err
		}
		lastHash = value
		return nil

	})
	if err != nil {
		logrus.Error("[-] Failed to run Init transaction in the datastore: ", err)
		return nil, err
	}

	return lastHash, nil
}

func (p *Persistence) Get(key []byte) ([]byte, error) {
	var value []byte
	err := p.db.View(func(transaction *badger.Txn) error {
		item, err := transaction.Get(key)
		if err != nil {
			logrus.Error("[-] Error occured while getting item by key: ", key, item, err)
			return err
		}
		value, err = item.ValueCopy(nil)
		if err != nil {
			logrus.Error("[-] Error occured while getting item value by key: ", key, item, err)
			return err
		}
		return err
	})
	if err != nil {
		logrus.Error("[-] Failed to run Get transaction in the datastore: ", err)
		return nil, err
	}
	return value, nil
}

func (p *Persistence) GetLastHash() ([]byte, error) {
	return p.Get([]byte(KeyLastHash))
}

func (p *Persistence) SaveBlock(hash []byte, block Serializable) error {
	err := p.db.Update(func(transaction *badger.Txn) error {
		serialData, err := block.Serialize()
		if err != nil {
			return err
		}
		if err := transaction.Set(hash, serialData); err != nil {
			return err
		}
		if err := transaction.Set([]byte(KeyLastHash), hash); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logrus.Error("[-] Failed to run SaveBlock transaction in the datastore: ", err)
		return err
	}
	return nil
}

func (p *Persistence) Iterate(prefix []byte, block Serializable, callback func(value []byte) error) error {
	err := p.db.View(func(transaction *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		iterator := transaction.NewIterator(opts)
		defer iterator.Close()

		for iterator.Seek(prefix); iterator.ValidForPrefix(prefix); iterator.Next() {
			item := iterator.Item()
			err := item.Value(callback)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
