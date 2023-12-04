package pubsub

import (
	"bufio"
	"context"
	"os"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"
)

type SubscriptionHandler interface {
	HandleMessage(msg *pubsub.Message)
}

type Manager struct {
	ctx           context.Context
	topics        map[string]*pubsub.Topic
	subscriptions map[string]*pubsub.Subscription
	gossipSub     *pubsub.PubSub
	host          host.Host
}

func NewPubSubManager(ctx context.Context, host host.Host) (*Manager, error) {
	gossipSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}
	manager := &Manager{
		ctx:           ctx,
		subscriptions: make(map[string]*pubsub.Subscription),
		topics:        make(map[string]*pubsub.Topic),
		gossipSub:     gossipSub,
		host:          host,
	}
	return manager, nil
}

func SetUpSubscriptions() {

}

func (sm *Manager) AddSubscription(topicName string, handler SubscriptionHandler) error {
	// Subscribe to a topic
	topic, err := sm.gossipSub.Join(topicName)
	if err != nil {
		return err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	sm.topics[topicName] = topic
	sm.subscriptions[topicName] = sub

	go func() {
		for {
			msg, err := sub.Next(sm.ctx)
			if err != nil {
				logrus.Errorf("Error reading from topic: %v", err)
				continue
			}
			// Skip messages from the same node
			if msg.ReceivedFrom == sm.host.ID() {
				continue
			}
			// Use the handler to process the message
			handler.HandleMessage(msg)
		}
	}()

	return nil
}

func (sm *Manager) RemoveSubscription(topic string) {
	delete(sm.subscriptions, topic)
}

func (sm *Manager) GetSubscription(topic string) *pubsub.Subscription {
	return sm.subscriptions[topic]
}

func (sm *Manager) Publish(topic string, data []byte) error {
	t, ok := sm.topics[topic]
	if !ok {
		return nil
	}
	return t.Publish(sm.ctx, data)
}

func StreamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			// Add check for EOF error and continue
			if err.Error() == "EOF" {
				continue
			}
			logrus.Errorf("streamConsoleTo: %s", err.Error())
		}
		if err := topic.Publish(ctx, []byte(s)); err != nil {
			logrus.Errorf("### Publish error: %s", err)
		}
	}
}
