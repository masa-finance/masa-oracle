package network

import (
	"bufio"
	"context"
	"os"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
)

type SubscriptionHandler interface {
	HandleMessage(msg *pubsub.Message)
}

func SubscribeToTopic(ctx context.Context, host host.Host, topic *pubsub.Topic, handler SubscriptionHandler) error {
	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}

	go func() {
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				logrus.Errorf("Error reading from topic: %v", err)
				continue
			}
			// Skip messages from the same node
			if msg.ReceivedFrom == host.ID() {
				continue
			}

			// Use the handler to process the message
			handler.HandleMessage(msg)
		}
	}()
	return nil
}

// WithPubSub TODO: this code is no longer used and should be removed as soon as the PubSubManager is successfully tested
func WithPubSub(ctx context.Context, host host.Host, topicName string, peerChan chan PeerEvent) (*pubsub.Topic, error) {
	gossipSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	// Subscribe to a topic
	topic, err := gossipSub.Join(topicName)
	if err != nil {
		return nil, err
	}
	go StreamConsoleTo(ctx, topic)

	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	// Read messages from the subscription
	go func() {
		for {
			m, err := sub.Next(ctx)
			if err != nil {
				logrus.Errorf("sub.Next: %s", err.Error())
			}
			// Skip messages from the same node
			if m.ReceivedFrom == host.ID() {
				continue
			}
			// Get the peer's IP address
			//var addrs multiaddr.Multiaddr
			connectedness := host.Network().Connectedness(m.ReceivedFrom)
			if connectedness == network.Connected {
				peerInfo := host.Peerstore().PeerInfo(m.ReceivedFrom)
				if len(peerInfo.Addrs) == 0 {
					continue
				}
				pe := PeerEvent{
					AddrInfo: peer.AddrInfo{ID: peerInfo.ID},
					Action:   PeerAdded,
					Source:   "topic",
				}
				peerChan <- pe

				//addrs = peerInfo.Addrs[0]
				//logrus.Infof("%s : %s : %s", m.ReceivedFrom, string(m.Message.Data), addrs.String())
			} else {
				logrus.Info(m.ReceivedFrom, ": ", string(m.Message.Data))
			}
		}
	}()
	return topic, nil
}

type PubSubManager struct {
	ctx           context.Context
	topics        map[string]*pubsub.Topic
	subscriptions map[string]*pubsub.Subscription
	gossipSub     *pubsub.PubSub
	host          host.Host
}

func NewPubSubManager(ctx context.Context, host host.Host) (*PubSubManager, error) {
	gossipSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}
	manager := &PubSubManager{
		ctx:           ctx,
		subscriptions: make(map[string]*pubsub.Subscription),
		topics:        make(map[string]*pubsub.Topic),
		gossipSub:     gossipSub,
		host:          host,
	}
	return manager, nil
}

func (sm *PubSubManager) AddSubscription(topicName string, handler SubscriptionHandler) error {
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

func (sm *PubSubManager) RemoveSubscription(topic string) {
	delete(sm.subscriptions, topic)
}

func (sm *PubSubManager) GetSubscription(topic string) *pubsub.Subscription {
	return sm.subscriptions[topic]
}

func (sm *PubSubManager) Publish(topic string, data []byte) error {
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
