package workerstatus

import (
	"encoding/json"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type WorkerStatus struct {
	PeerID string `json:"peerId"`
	Data   []byte
}

// SubscriptionHandler is a struct that handles subscriptions for worker status updates.
// It contains the following fields:
// - WorkerStatus: A slice of WorkerStatus structs representing the status of workers.
// - Data: A byte slice containing the raw data received from subscriptions.
// - mu: A sync.Mutex used for synchronizing access to the handler's fields.
// - WorkerCh: A channel for sending worker status updates as byte slices.
type SubscriptionHandler struct {
	WorkerStatus       []WorkerStatus
	CompletedWorkTopic *pubsub.Topic
	mu                 sync.Mutex
	WorkerStatusCh     chan []byte
}

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	workerStatus := WorkerStatus{}
	err := json.Unmarshal(message.Data, &workerStatus)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}

	handler.mu.Lock()
	handler.WorkerStatus = append(handler.WorkerStatus, workerStatus)
	handler.mu.Unlock()

	jsonData, _ := json.Marshal(workerStatus)
	handler.WorkerStatusCh <- jsonData
}
