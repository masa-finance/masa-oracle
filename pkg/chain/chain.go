package chain

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Chain struct {
	LastHash     []byte
	storage      *Persistance
	CurrentBlock uint64
}

// Init initializes the blockchain.
//
// This function performs the following tasks:
// 1. Creates a data directory for storing blocks if it doesn't exist.
// 2. Initializes the storage for the blockchain.
// 3. Creates and stores the genesis block if the blockchain is empty.
//
// Returns:
//   - error: An error if any step in the initialization process fails, nil otherwise.
func (c *Chain) Init(path string) error {
	dataDir := filepath.Join(path, "./blocks")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			logrus.Fatal("[-] Failed to create directory: ", err)
		}
	}
	logrus.WithFields(logrus.Fields{"block": Difficulty}).Info("[+] Initializing blockchain...")
	c.storage = &Persistance{}
	c.storage.Init(dataDir, func() (Serializable, []byte) {
		genesisBlock := makeGenesisBlock()
		return genesisBlock, genesisBlock.Hash
	})
	return nil
}

// makeGenesisBlock creates and returns the genesis block for the blockchain.
//
// This function:
//  1. Logs the creation of the genesis block.
//  2. Initializes a new Block.
//  3. Builds the block with "Genesis" as data, an empty link (as it's the first block),
//     and a nonce of 1.
//
// Returns:
//   - *Block: A pointer to the newly created genesis block.
func makeGenesisBlock() *Block {
	logrus.Info("[+] Generating genesis block...")
	newBlock := &Block{}
	emptyLink := []byte{}
	newBlock.Build([]byte("Genesis"), emptyLink, big.NewInt(1), 0)
	return newBlock
}

// UpdateLastHash updates the LastHash field of the Chain struct with the most recent hash from storage.
//
// This function:
// 1. Retrieves the last hash from the storage.
// 2. Updates the LastHash field of the Chain struct with the retrieved hash.
//
// Returns:
//   - error: An error if retrieving the last hash from storage fails, nil otherwise.
func (c *Chain) UpdateLastHash() error {
	lastHash, err := c.storage.GetLastHash()
	if err != nil {
		logrus.Error("[-] Failed to get last hash from the storage: ", err)
		return err
	}
	c.LastHash = lastHash
	return nil
}

// AddBlock adds a new block to the blockchain with the given data.
//
// This function:
// 1. Updates the last hash of the chain.
// 2. Creates a new block with the provided data, the last hash, and a nonce of 1.
// 3. Validates the new block using Proof of Stake (PoS).
// 4. Saves the new block to storage if valid.
// 5. Updates the LastHash of the chain to the hash of the new block.
//
// Parameters:
//   - data: The data to be included in the new block.
//
// Returns:
//   - error: An error if any step fails, nil otherwise.
func (c *Chain) AddBlock(data []byte) error {
	logrus.Info("[+] Adding block...")
	if err := c.UpdateLastHash(); err != nil {
		return err
	}
	newBlock := &Block{}
	nextBlockNumber := c.getNextBlockNumber()
	newBlock.Build(data, c.LastHash, big.NewInt(1), nextBlockNumber)

	if !IsValidPoS(newBlock, big.NewInt(1)) {
		logrus.Error("[-] Invalid PoS block")
		return fmt.Errorf("invalid PoS block")
	}

	err := c.storage.SaveBlock(newBlock.Hash, newBlock)
	if err != nil {
		logrus.Error("[-] Failed to save block into the storage: ", newBlock, err)
		return err
	}
	c.LastHash = newBlock.Hash
	c.CurrentBlock = nextBlockNumber
	return nil
}

func (c *Chain) getNextBlockNumber() uint64 {
	return c.CurrentBlock + 1
}

// IterateLink iterates through the blockchain, executing provided functions at specific points.
//
// This function:
// 1. Updates the last hash of the chain.
// 2. Executes a pre-iteration function.
// 3. Iterates through each block in the chain, starting from the last block.
// 4. For each block, it executes a provided function.
// 5. After iteration, executes a post-iteration function.
//
// Parameters:
//   - each: A function to be executed for each block in the chain.
//   - pre: A function to be executed before the iteration begins.
//   - post: A function to be executed after the iteration completes.
//
// Returns:
//   - error: An error if any operation fails during iteration, nil otherwise.
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

// GetLastBlock retrieves the most recent block in the blockchain.
//
// This function:
// 1. Updates the last hash of the chain to ensure it's current.
// 2. Retrieves the block corresponding to the last hash.
//
// Returns:
//   - *Block: A pointer to the last block in the chain.
//   - error: An error if the retrieval fails, nil otherwise.
func (c *Chain) GetLastBlock() (*Block, error) {
	c.UpdateLastHash()
	return c.GetBlock(c.LastHash)
}

func (c *Chain) GetBlock(hash []byte) (*Block, error) {
	logrus.Infof("[+] transaction %x", hash)
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

// GetBlockByHash retrieves a specific block from the blockchain using its hash.
//
// Parameters:
//   - c: A pointer to the Chain instance.
//   - hash: A byte slice representing the hash of the block to retrieve.
//
// Returns:
//   - *Block: A pointer to the retrieved Block if found.
//   - error: An error if the block retrieval fails, nil otherwise.
//
// This function attempts to fetch the block data from storage using the provided hash,
// then deserializes the data into a Block struct. If any step fails, an error is returned.
func GetBlockByHash(c *Chain, hash []byte) (*Block, error) {
	logrus.Infof("[+] transaction %x", hash)
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

// GetBlockchain retrieves all blocks from the blockchain.
//
// Parameters:
//   - c: A pointer to the Chain instance.
//
// Returns:
//   - []*Block: A slice of pointers to Block, representing all blocks in the chain.
//
// This function iterates through the entire blockchain, collecting all blocks
// into a slice. If an error occurs during iteration, it logs the error and
// returns nil. The blocks are returned in the order they are stored in the chain.
func GetBlockchain(c *Chain) []*Block {
	var blockchain []*Block

	each := func(b *Block) {
		blockchain = append(blockchain, b)
	}

	err := c.IterateLink(each, func() {}, func() {})
	if err != nil {
		logrus.Errorf("[-] Error iterating through blockchain: %v", err)
		return nil
	}

	return blockchain
}
