package network

import (
	"bufio"
	"context"
	"fmt"
	"os"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sirupsen/logrus"
)

func NewPubSub(ctx context.Context, host host.Host, topicName string) (*pubsub.Topic, error) {
	gossipSub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	// Subscribe to a topic
	topic, err := gossipSub.Join(topicName)
	if err != nil {
		return nil, err
	}
	go streamConsoleTo(ctx, topic)

	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	// Read messages from the subscription
	go func() {
		for {
			m, err := sub.Next(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
		}
	}()
	return topic, nil
}

func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			logrus.Error(err)
		}
		if err := topic.Publish(ctx, []byte(s)); err != nil {
			fmt.Println("### Publish error:", err)
		}
	}
}
