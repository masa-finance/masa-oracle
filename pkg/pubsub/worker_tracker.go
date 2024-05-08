package pubsub

import (
	"encoding/json"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type Workers struct {
}

// WorkerEventTracker is a struct that handles subscriptions for worker status updates.
// It contains the following fields:
// - WorkerTracker: A slice of WorkerTracker structs representing the status of workers.
// - Data: A byte slice containing the raw data received from subscriptions.
// - mu: A sync.Mutex used for synchronizing access to the handler's fields.
// - WorkerCh: A channel for sending worker status updates as byte slices.
type WorkerEventTracker struct {
	Workers        []Workers
	WorkerTopic    *pubsub.Topic
	mu             sync.Mutex
	WorkerStatusCh chan *pubsub.Message
}

// HandleMessage implements subscription WorkerEventTracker handler
func (h *WorkerEventTracker) HandleMessage(m *pubsub.Message) {
	logrus.Infof("workerStream -> Received work from: %s", m.ReceivedFrom)
	var workers Workers
	err := json.Unmarshal(m.Data, &workers)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}
	h.mu.Lock()
	h.Workers = append(h.Workers, workers)
	h.mu.Unlock()
	h.WorkerStatusCh <- m
}
