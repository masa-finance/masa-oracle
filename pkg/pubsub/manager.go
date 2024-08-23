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

// SubscriptionHandler defines the interface for handling pubsub messages.
// Implementations should subscribe to topics and handle incoming messages.
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

// NewPubSubManager creates a new PubSubManager instance.
// It initializes a new GossipSub and associates it with the given host.
// It also initializes data structures to track topics, subscriptions and handlers.
// The PublicKeyPublisher is initialized to enable publishing public keys over pubsub.
// The manager instance is returned, along with any error from initializing GossipSub.
func NewPubSubManager(ctx context.Context, host host.Host) (*Manager, error) { // Modify this line to accept pubKey
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

// SetUpSubscriptions sets up default subscriptions for the PubSub manager
// based on predefined topics and handlers. This allows initializing subscriptions
// separately from creating the handlers.
func (sm *Manager) SetUpSubscriptions() {
}

// createTopic joins a PubSub topic with the given topic name,
// adds it to the manager's topic map, and returns the topic
// instance along with any error from joining.
func (sm *Manager) createTopic(topicName string) (*pubsub.Topic, error) {
	topic, err := sm.gossipSub.Join(topicName)
	if err != nil {
		return nil, err
	}
	sm.topics[topicName] = topic
	return topic, nil
}

// AddSubscription subscribes to the PubSub topic with the given topicName.
// It creates the topic if needed, subscribes to it, and adds the subscription
// and handler to the manager's maps. It launches a goroutine to handle incoming
// messages, skipping messages from self, and calling the handler on each message.
func (sm *Manager) AddSubscription(topicName string, handler SubscriptionHandler, includeSelf bool) error {
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
				logrus.Errorf("[-] Error reading from topic: %v", err)
				continue
			}
			if !includeSelf {
				// if !includeSelf && msg.ReceivedFrom == sm.host.ID() {
				// if msg.ReceivedFrom == sm.host.ID() {
				// Skip messages from the same node
				continue
			}
			// Use the handler to process the message
			handler.HandleMessage(msg)
		}
	}()

	return nil
}

// RemoveSubscription unsubscribes from the PubSub topic with the given
// topic name. It closes the existing subscription, removes it from the
// manager's subscription map, and removes the associated handler. Returns
// an error if no subscription exists for the given topic.
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

// GetSubscription returns the Subscription for the given topic name.
// It returns an error if no Subscription exists for the given topic.
func (sm *Manager) GetSubscription(topic string) (*pubsub.Subscription, error) {
	sub, ok := sm.subscriptions[topic]
	if !ok {
		return nil, fmt.Errorf("no subscription for topic %s", topic)
	}
	return sub, nil
}

// Publish publishes a message to the PubSub topic with the given topic name.
// It returns an error if no topic with the given name exists.
func (sm *Manager) Publish(topic string, data []byte) error {
	t, ok := sm.topics[topic]
	if !ok {
		return fmt.Errorf("no topic named %s", topic)
	}
	return t.Publish(sm.ctx, data)
}

// GetHandler returns the SubscriptionHandler for the given topic name.
// It returns an error if no handler exists for the given topic.
func (sm *Manager) GetHandler(topic string) (SubscriptionHandler, error) {
	handler, ok := sm.handlers[topic]
	if !ok {
		return nil, fmt.Errorf("no handler for topic %s", topic)
	}
	return handler, nil
}

// StreamConsoleTo streams data read from stdin to the given PubSub topic.
// It launches a goroutine that continuously reads from stdin using a bufio.Reader.
// Each line that is read is published to the topic. Any errors are logged.
// The goroutine runs until ctx is canceled.
func StreamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			// Add check for EOF error and continue
			if err.Error() == "EOF" {
				continue
			}
			logrus.Errorf("[-] streamConsoleTo: %s", err.Error())
		}
		if err := topic.Publish(ctx, []byte(s)); err != nil {
			logrus.Errorf("[-] Publish error: %s", err)
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

// PublishMessage publishes a message to the PubSub topic with the given topicName.
// It converts the message to a byte slice, checks if the topic exists,
// optionally creates the topic if it doesn't exist, and publishes using the
// existing Publish method.
// Returns an error if the topic does not exist and cannot be created.
func (sm *Manager) PublishMessage(topicName, message string) error {
	// Convert the message string to a byte slice
	data := []byte(message)

	// Check if the topic exists
	t, ok := sm.topics[topicName]
	if !ok {
		// Optionally, create the topic if it doesn't exist
		var err error
		t, err = sm.createTopic(topicName)
		if err != nil {
			return fmt.Errorf("[-] Failed to create topic %s: %w", topicName, err)
		}
	}

	// Use the existing Publish method to publish the message
	return t.Publish(sm.ctx, data)
}

// Subscribe registers a subscription handler to receive messages for the
// given topic name. It gets the existing subscription, saves it and the
// handler, and starts a goroutine to call the handler for each new message.
// Returns an error if unable to get the subscription.
func (sm *Manager) Subscribe(topicName string, handler SubscriptionHandler) error {
	sub, err := sm.GetSubscription(topicName)
	if err != nil {
		return err
	}
	sm.subscriptions[topicName] = sub
	sm.handlers[topicName] = handler

	go func() {
		for {
			msg, err := sub.Next(sm.ctx)
			if err != nil {
				logrus.Errorf("[-] Error reading from topic: %v", err)
				continue
			}
			// Skip messages from the same node
			//if msg.ReceivedFrom == sm.host.ID() {
			//	continue
			//}
			// Use the handler to process the message
			handler.HandleMessage(msg)
		}
	}()
	return nil
}
