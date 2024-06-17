package blockchain

import (
	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Room represents a subscription to a single PubSub topic. Messages
// can be published to the topic with Room.Publish, and received
// messages are pushed to the Messages channel.
type room struct {
	ctx   context.Context
	ps    *pubsub.PubSub
	Topic *pubsub.Topic
	sub   *pubsub.Subscription

	roomName string
	self     peer.ID
}

// connect tries to subscribe to the PubSub topic for the room name, returning
// a Room on success.
func connect(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, roomName string, messageChan chan *Message) (*room, error) {
	// join the pubsub topic
	topic, err := ps.Join(roomName)
	if err != nil {
		return nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	cr := &room{
		ctx:      ctx,
		ps:       ps,
		Topic:    topic,
		sub:      sub,
		self:     selfID,
		roomName: roomName,
	}

	// start reading messages from the subscription in a loop
	go cr.readLoop(messageChan)
	return cr, nil
}

// publishMessage sends a message to the pubsub topic.
func (cr *room) publishMessage(m *Message) error {
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cr.Topic.Publish(cr.ctx, msgBytes)
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (cr *room) readLoop(messageChan chan *Message) {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == cr.self {
			continue
		}
		cm := new(Message)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}

		cm.SenderID = msg.ReceivedFrom.String()

		// send valid messages onto the Messages channel
		messageChan <- cm
	}
}
