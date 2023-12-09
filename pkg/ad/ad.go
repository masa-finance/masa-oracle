package ad

import (
	"encoding/json"

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
}

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	var ad Ad
	err := json.Unmarshal(message.Data, &ad)
	if err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return
	}
	handler.Ads = append(handler.Ads, ad) // Add the ad to the list

	// Handle the ad here
	logrus.Infof("received ad: %v", ad)
}
