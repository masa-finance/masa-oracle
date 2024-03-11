/*
Package for Oracle Node Subscription Management in the Masa Oracle system. It enables OracleNodes to subscribe to network topics, ensuring they receive and process relevant data. The SubscribeToTopics function is central to this, facilitating node participation by keeping them updated, thus supporting the network's consensus mechanism. Ideal for developers needing to manage node subscriptions within the network.
*/
package masa

import (
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/ad"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/nodestatus"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
)

// SubscribeToTopics handles the subscription to various topics for an OracleNode.
// It subscribes the node to the NodeGossipTopic, AdTopic, and PublicKeyTopic.
// Each subscription is managed through the node's PubSubManager, which orchestrates the message passing for these topics.
// Errors during subscription are logged and returned, halting the process to ensure the node's correct setup before operation.
func SubscribeToTopics(node *OracleNode) error {
	// Subscribe to NodeGossipTopic to participate in the network's gossip protocol.
	if err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.NodeGossipTopic), node.NodeTracker); err != nil {
		return err
	}

	// Initialize and subscribe to AdTopic for receiving advertisement-related messages.
	node.AdSubscriptionHandler = &ad.SubscriptionHandler{}
	if err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.AdTopic), node.AdSubscriptionHandler); err != nil {
		logrus.Errorf("Failed to subscribe to ad topic: %v", err)
		return err
	}

	// Initialize and subscribe to NodeStatusTopic for receiving advertisement-related messages.
	node.NodeStatusSubscriptionsHandler = &nodestatus.SubscriptionHandler{}
	if err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.NodeStatusTopic), node.NodeStatusSubscriptionsHandler); err != nil {
		logrus.Errorf("Failed to subscribe to nodeStatus topic: %v", err)
		return err
	}

	// Subscribe to PublicKeyTopic to manage and verify public keys within the network.
	if err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.PublicKeyTopic), &pubsub2.PublicKeySubscriptionHandler{}); err != nil {
		return err
	}

	return nil
}
