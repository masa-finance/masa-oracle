package pubsub

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

// TopicHandler is responsible for handling messages from subscribed topics.
type TopicHandler struct {
	Subscription *pubsub.Subscription
}

// NewTopicHandler creates a new TopicHandler with necessary initializations.
func NewTopicHandler() *TopicHandler {
	return &TopicHandler{}
}

// StartListening starts listening to messages on the subscribed topic.
func (h *TopicHandler) StartListening() {
	go func() {
		for {
			msg, err := h.Subscription.Next(context.Background()) // Use the context package as needed
			if err != nil {
				logrus.Error("Error while reading message: ", err)
				return
			}
			h.HandleMessage(msg)
		}
	}()
}

// HandleMessage processes messages received on the subscribed topics.
func (h *TopicHandler) HandleMessage(msg *pubsub.Message) {
	logrus.Infof("Received message on topic: %s", string(msg.Data))
}
