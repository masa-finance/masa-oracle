package network

import (
	"bufio"
	"context"
	"os"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
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
				logrus.Errorf("sub.Next: %s", err.Error())
			}
			// Skip messages from the same node
			if m.ReceivedFrom == host.ID() {
				continue
			}
			// Get the peer's IP address
			var addrs multiaddr.Multiaddr
			connectedness := host.Network().Connectedness(m.ReceivedFrom)
			if connectedness == network.Connected {
				peerInfo := host.Peerstore().PeerInfo(m.ReceivedFrom)
				if len(peerInfo.Addrs) == 0 {
					continue
				}
				addrs = peerInfo.Addrs[0]
				logrus.Infof("%s : %s : %s", m.ReceivedFrom, string(m.Message.Data), addrs.String())
			} else {
				logrus.Info(m.ReceivedFrom, ": ", string(m.Message.Data))
			}
		}
	}()
	return topic, nil
}

func streamConsoleTo(ctx context.Context, topic *pubsub.Topic) {
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
