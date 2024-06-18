package pubsub

import (
	"encoding/json"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type Blocks struct {
}

// BlockEventTracker is a struct that handles subscriptions for worker status updates.
type BlockEventTracker struct {
	Blocks     []Blocks
	BlockTopic *pubsub.Topic
	mu         sync.Mutex
	BlocksCh   chan *pubsub.Message
}

// HandleMessage implements subscription BlockEventTracker handler
func (b *BlockEventTracker) HandleMessage(m *pubsub.Message) {
	logrus.Infof("chain -> Received block from: %s", m.ReceivedFrom)
	var blocks Blocks
	err := json.Unmarshal(m.Data, &blocks)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}
	b.mu.Lock()
	b.Blocks = append(b.Blocks, blocks)
	b.mu.Unlock()
	b.BlocksCh <- m
}
