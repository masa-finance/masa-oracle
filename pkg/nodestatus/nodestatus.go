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
	FirstJoined               time.Time     `json:"firstJoined"`
	LastJoined                time.Time     `json:"lastJoined"`
	LastLeft                  time.Time     `json:"lastLeft"`
	LastUpdated               time.Time     `json:"lastUpdated"`
	CurrentUptime             time.Duration `json:"currentUptime"`
	ReadableCurrentUptime     string        `json:"readableCurrentUptime"`
	AccumulatedUptime         time.Duration `json:"accumulatedUptime"`
	ReadableAccumulatedUptime string        `json:"readableAccumulatedUptime"`
	IsActive                  bool          `json:"isActive"`
	IsStaked                  bool          `json:"isStaked"`
	IsWriterNode              bool          `json:"isWriterNode"`
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
