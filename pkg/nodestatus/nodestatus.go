package nodestatus

import (
	"encoding/json"
	"sync"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type NodeStatus struct {
	PeerID                    string        `json:"peerId"`
	IsActive                  bool          `json:"isActive"`
	IsStaked                  bool          `json:"isStaked"`
	IsWriterNode              bool          `json:"isWriterNode"`
	AccumulatedUptime         time.Duration `json:"accumulatedUptime"`
	CurrentUptime             time.Duration `json:"currentUptime"`
	ReadableAccumulatedUptime string        `json:"readableAccumulatedUptime"`
	FirstJoined               time.Time     `json:"firstJoined"`
	LastJoined                time.Time     `json:"lastJoined"`
}

// SubscriptionHandler handles storing node status updates and publishing
// them to the node status topic.
type SubscriptionHandler struct {
	NodeStatus      []NodeStatus
	NodeStatusTopic *pubsub.Topic
	mu              sync.Mutex
	NodeStatusCh    chan []byte
}

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	nodeStatus := NodeStatus{}
	err := json.Unmarshal(message.Data, &nodeStatus)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}

	handler.mu.Lock()
	handler.NodeStatus = append(handler.NodeStatus, nodeStatus)
	handler.mu.Unlock()

	jsonData, _ := json.Marshal(nodeStatus)
	handler.NodeStatusCh <- jsonData
}
