package chain

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/sirupsen/logrus"
)

type Chain struct {
	LastHash []byte
	storage  *Persistance
}

func (c *Chain) Init() error {
	dataDir := filepath.Join(config.GetInstance().MasaDir, "./blocks")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			logrus.Fatal("Failed to create directory: ", err)
		}
	}
	logrus.WithFields(logrus.Fields{"block": Difficulty}).Info("Initializing blockchain...")
	c.storage = &Persistance{}
	c.storage.Init(dataDir, func() (Serializable, []byte) {
		genesisBlock := makeGenesisBlock()
		return genesisBlock, genesisBlock.Hash
	})
	return nil
}

func makeGenesisBlock() *Block {
	logrus.Info("Generating genesis transaction...")
	newBlock := &Block{}
	emptyLink := []byte{}
	newBlock.Build([]byte("Genesis"), emptyLink, big.NewInt(1))
	return newBlock
}

func (c *Chain) UpdateLastHash() error {
	logrus.Info("Fetching last transaction...")
	lastHash, err := c.storage.GetLastHash()
	if err != nil {
		logrus.Error("Failed to get last hash from the storage: ", err)
		return err
	}
	c.LastHash = lastHash
	return nil
}

func (c *Chain) AddBlock(data []byte) error {
	logrus.Info("Adding transaction...")
	if err := c.UpdateLastHash(); err != nil {
		return err
	}
	newBlock := &Block{}
	newBlock.Build(data, c.LastHash, big.NewInt(1))

	if !IsValidPoS(newBlock, big.NewInt(1)) {
		logrus.Error("Invalid PoS block")
		return fmt.Errorf("invalid PoS block")
	}

	err := c.storage.SaveBlock(newBlock.Hash, newBlock)
	if err != nil {
		logrus.Error("Failed to save block into the storage: ", newBlock, err)
		return err
	}
	c.LastHash = newBlock.Hash
	return nil
}

func (c *Chain) IterateLink(each func(b *Block), pre, post func()) error {
	c.UpdateLastHash()
	currentHash := c.LastHash
	pre()
	for len(currentHash) > 0 {
		data, err := c.storage.Get(currentHash)
		if err != nil {
			return err
		}
		block := &Block{}
		if err = block.Deserialize(data); err != nil {
			return err
		}
		each(block)
		currentHash = block.Link
	}
	post()
	return nil
}

func (c *Chain) GetLastBlock() (*Block, error) {
	c.UpdateLastHash()
	return c.GetBlock(c.LastHash)
}

func (c *Chain) GetBlock(hash []byte) (*Block, error) {
	logrus.Infof("transaction %x", hash)
	data, err := c.storage.Get(hash)
	if err != nil {
		return nil, err
	}
	block := &Block{}
	if err = block.Deserialize(data); err != nil {
		return nil, err
	}
	return block, nil
}
