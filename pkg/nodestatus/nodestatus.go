package nodestatus

import (
	"encoding/json"
	"sync"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type NodeStatus struct {
	PeerID       string        `json:"peerId"`
	IsStaked     bool          `json:"isStaked"`
	IsWriterNode bool          `json:"isWriterNode"`
	TotalUpTime  time.Duration `json:"totalUpTime"`
	FirstJoined  time.Time     `json:"firstJoined"`
	LastJoined   time.Time     `json:"lastJoined"`
}

// SubscriptionHandler handles storing node status updates and publishing
// them to the node status topic.
type SubscriptionHandler struct {
	NodeStatus      []NodeStatus
	NodeStatusTopic *pubsub.Topic
	mu              sync.Mutex
}

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	logrus.Infof("Received a message %s", message.Data)
	var nodeStatus NodeStatus
	err := json.Unmarshal(message.Data, &nodeStatus)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}

	handler.mu.Lock()
	handler.NodeStatus = append(handler.NodeStatus, nodeStatus)
	handler.mu.Unlock()

	logrus.Infof("NodeStatus received: %+v", nodeStatus)
}
