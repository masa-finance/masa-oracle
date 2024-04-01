package pubsub

import (
	"context"
	"encoding/json"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

// NodeLifecycleEvent represents a join or leave event
type NodeLifecycleEvent struct {
	EventType string `json:"eventType"` // "join" or "leave"
	NodeID    string `json:"nodeID"`
	Nonce     int64  `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
}

// BroadcastEvent marshals a NodeLifecycleEvent into JSON, publishes it
// to the given PubSub topic, and logs the operation. Returns any error.
// This allows broadcasting node join/leave events to other nodes.
func BroadcastEvent(ctx context.Context, ps *pubsub.PubSub, topicName string, event NodeLifecycleEvent) error {
	event.Timestamp = time.Now().Unix()

	eventBytes, err := json.Marshal(event)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal node lifecycle event")
		return err
	}

	topic, err := ps.Join(topicName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topicName": topicName,
			"error":     err,
		}).Error("Failed to join topic")
		return err
	}

	if err := topic.Publish(ctx, eventBytes); err != nil {
		logrus.WithFields(logrus.Fields{
			"topicName": topicName,
			"error":     err,
		}).Error("Failed to publish node lifecycle event")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"eventType": event.EventType,
		"nodeID":    event.NodeID,
		"topicName": topicName,
	}).Info("Successfully broadcasted node lifecycle event")
	return nil
}
