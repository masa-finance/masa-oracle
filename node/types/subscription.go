package types

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// SubscriptionHandler defines the interface for handling pubsub messages.
// Implementations should subscribe to topics and handle incoming messages.
type SubscriptionHandler interface {
	HandleMessage(msg *pubsub.Message)
}
