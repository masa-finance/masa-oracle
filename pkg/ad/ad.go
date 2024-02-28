package ad

import (
	"encoding/json"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type Ad struct {
	Content  string
	Metadata map[string]string
}

type SubscriptionHandler struct {
	Ads     []Ad
	AdTopic *pubsub.Topic
	mu      sync.Mutex
}

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	var ad Ad
	err := json.Unmarshal(message.Data, &ad)
	if err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return
	}

	handler.mu.Lock()
	handler.Ads = append(handler.Ads, ad)
	handler.mu.Unlock()

	// Handle the ad here
	logrus.Infof("received ad: %v", ad)
}
