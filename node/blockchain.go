package node

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/masa-finance/masa-oracle/pkg/chain"
	"github.com/sirupsen/logrus"
)

// Blockchain Implementation

type BlockData struct {
	Block            uint64      `json:"block"`
	InputData        interface{} `json:"input_data"`
	TransactionHash  string      `json:"transaction_hash"`
	PreviousHash     string      `json:"previous_hash"`
	TransactionNonce int         `json:"nonce"`
}

type Blocks struct {
	BlockData []BlockData `json:"blocks"`
}

type BlockEvents map[string]interface{}

type BlockEventTracker struct {
	BlockEvents []BlockEvents
	mu          sync.Mutex
	blocksCh    chan *pubsub.Message
}

func NewBlockChain() *BlockEventTracker {
	return &BlockEventTracker{
		blocksCh: make(chan *pubsub.Message),
	}
}

// HandleMessage processes incoming pubsub messages containing block events.
// It unmarshals the message data into a slice of BlockEvents and appends them
// to the tracker's BlockEvents slice.
func (b *BlockEventTracker) HandleMessage(m *pubsub.Message) {
	var blockEvents any

	// Try to decode as base64 first
	decodedData, err := base64.StdEncoding.DecodeString(string(m.Data))
	if err == nil {
		m.Data = decodedData
	}

	// Try to unmarshal as JSON
	err = json.Unmarshal(m.Data, &blockEvents)
	if err != nil {
		// If JSON unmarshal fails, try to interpret as string
		blockEvents = string(m.Data)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	switch v := blockEvents.(type) {
	case []BlockEvents:
		b.BlockEvents = append(b.BlockEvents, v...)
	case BlockEvents:
		b.BlockEvents = append(b.BlockEvents, v)
	case map[string]interface{}:
		// Convert map to BlockEvents struct
		newBlockEvent := BlockEvents(v)
		// You might need to add logic here to properly convert the map to BlockEvents
		b.BlockEvents = append(b.BlockEvents, newBlockEvent)
	case []interface{}:
		// Convert each item in the slice to BlockEvents
		for _, item := range v {
			if be, ok := item.(BlockEvents); ok {
				b.BlockEvents = append(b.BlockEvents, be)
			}
		}
	case string:
		// Handle string data
		newBlockEvent := BlockEvents{}
		// You might need to add logic here to properly convert the string to BlockEvents
		b.BlockEvents = append(b.BlockEvents, newBlockEvent)
	default:
		logrus.Warnf("[-] Unexpected data type in message: %v", reflect.TypeOf(v))
	}

	b.blocksCh <- m
}

func updateBlocks(ctx context.Context, node *OracleNode) error {

	var existingBlocks Blocks
	blocks := chain.GetBlockchain(node.Blockchain)

	for _, block := range blocks {
		var inputData interface{}
		err := json.Unmarshal(block.Data, &inputData)
		if err != nil {
			inputData = string(block.Data) // Fallback to string if unmarshal fails
		}

		blockData := BlockData{
			Block:            block.Block,
			InputData:        base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", inputData))),
			TransactionHash:  fmt.Sprintf("%x", block.Hash),
			PreviousHash:     fmt.Sprintf("%x", block.Link),
			TransactionNonce: int(block.Nonce),
		}
		existingBlocks.BlockData = append(existingBlocks.BlockData, blockData)
	}
	jsonData, err := json.Marshal(existingBlocks)
	if err != nil {
		return err
	}

	err = node.DHT.PutValue(ctx, "/db/blocks", jsonData)
	if err != nil {
		logrus.Warningf("[-] Unable to store block on DHT: %v", err)
	}

	if os.Getenv("IPFS_URL") != "" {

		infuraURL := fmt.Sprintf("https://%s:%s@%s", os.Getenv("PID"), os.Getenv("PS"), os.Getenv("IPFS_URL"))
		sh := shell.NewShell(infuraURL)

		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			logrus.Errorf("[-] Error marshalling JSON: %s", err)
		}

		reader := bytes.NewReader(jsonBytes)

		hash, err := sh.AddWithOpts(reader, true, true)
		if err != nil {
			logrus.Errorf("[-] Error persisting to IPFS: %s", err)
		} else {
			logrus.Printf("[+] Ledger persisted with IPFS hash: https://dwn.infura-ipfs.io/ipfs/%s\n", hash)
			_ = node.DHT.PutValue(ctx, "/db/ipfs", []byte(fmt.Sprintf("https://dwn.infura-ipfs.io/ipfs/%s", hash)))

		}
	}

	return nil
}

func (b *BlockEventTracker) Start(path string) func(ctx context.Context, node *OracleNode) {
	return func(ctx context.Context, node *OracleNode) {
		err := node.Blockchain.Init(path)
		if err != nil {
			logrus.Error(err)
		}

		updateTicker := time.NewTicker(time.Second * 60)
		defer updateTicker.Stop()

		for {
			select {
			case block, ok := <-b.blocksCh:
				if !ok {
					logrus.Error("[-] Block channel closed")
					return
				}
				if err := processBlock(node, block); err != nil {
					logrus.Errorf("[-] Error processing block: %v", err)
					// Consider adding a retry mechanism or circuit breaker here
				}

			case <-updateTicker.C:
				logrus.Info("[+] blockchain tick")
				if err := updateBlocks(ctx, node); err != nil {
					logrus.Errorf("[-] Error updating blocks: %v", err)
					// Consider adding a retry mechanism or circuit breaker here
				}

			case <-ctx.Done():
				logrus.Info("[+] Context cancelled, stopping block subscription")
				return
			}
		}
	}
}

func processBlock(node *OracleNode, block *pubsub.Message) error {
	blocks := chain.GetBlockchain(node.Blockchain)
	for _, b := range blocks {
		if string(b.Data) == string(block.Data) {
			return nil // Block already exists
		}
	}

	if err := node.Blockchain.AddBlock(block.Data); err != nil {
		return fmt.Errorf("[-] failed to add block: %w", err)
	}

	if node.Blockchain.LastHash != nil {
		b, err := node.Blockchain.GetBlock(node.Blockchain.LastHash)
		if err != nil {
			return fmt.Errorf("[-] failed to get last block: %w", err)
		}
		b.Print()
	}

	return nil
}
