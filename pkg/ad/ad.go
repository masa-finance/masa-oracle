package ad

import (
	"encoding/json"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

// Ad represents an advertisement with content and metadata.
type Ad struct {
	Content  string
	Metadata map[string]string
}

// SubscriptionHandler handles storing advertisements and publishing
// them to the advertisement topic.
type SubscriptionHandler struct {
	Ads     []Ad
	AdTopic *pubsub.Topic
	mu      sync.Mutex
}

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	logrus.Info("Received a message")
	var ad Ad
	err := json.Unmarshal(message.Data, &ad)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}

	handler.mu.Lock()
	handler.Ads = append(handler.Ads, ad)
	handler.mu.Unlock()

	logrus.Infof("Ad added: %+v", ad)
}
