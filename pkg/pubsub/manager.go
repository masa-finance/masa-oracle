package pubsub

import (
	"bufio"
	"context"
	"fmt"
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
	handlers      map[string]SubscriptionHandler
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
		handlers:      make(map[string]SubscriptionHandler),
		gossipSub:     gossipSub,
		host:          host,
	}

	return manager, nil
}

// SetUpSubscriptions can be used to set up a default set of subscriptions where the handler can be created separately
func (sm *Manager) SetUpSubscriptions() {
}

func (sm *Manager) createTopic(topicName string) (*pubsub.Topic, error) {
	topic, err := sm.gossipSub.Join(topicName)
	if err != nil {
		return nil, err
	}
	sm.topics[topicName] = topic
	return topic, nil
}

func (sm *Manager) AddSubscription(topicName string, handler SubscriptionHandler) error {
	topic, err := sm.createTopic(topicName)
	if err != nil {
		return err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	sm.subscriptions[topicName] = sub
	sm.handlers[topicName] = handler

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

func (sm *Manager) RemoveSubscription(topic string) error {
	sub, ok := sm.subscriptions[topic]
	if !ok {
		return fmt.Errorf("no subscription for topic %s", topic)
	}
	// Close the subscription
	sub.Cancel()

	// Remove the subscription and handler
	delete(sm.subscriptions, topic)
	delete(sm.handlers, topic)
	return nil
}
func (sm *Manager) GetSubscription(topic string) (*pubsub.Subscription, error) {
	sub, ok := sm.subscriptions[topic]
	if !ok {
		return nil, fmt.Errorf("no subscription for topic %s", topic)
	}
	return sub, nil
}

func (sm *Manager) Publish(topic string, data []byte) error {
	t, ok := sm.topics[topic]
	if !ok {
		return fmt.Errorf("no topic named %s", topic)
	}
	return t.Publish(sm.ctx, data)
}

func (sm *Manager) GetHandler(topic string) (SubscriptionHandler, error) {
	handler, ok := sm.handlers[topic]
	if !ok {
		return nil, fmt.Errorf("no handler for topic %s", topic)
	}
	return handler, nil
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

// GetTopicNames returns a slice of the names of all topics currently managed.
func (sm *Manager) GetTopicNames() []string {
	var topicNames []string
	for name := range sm.topics {
		topicNames = append(topicNames, name)
	}
	return topicNames
}
