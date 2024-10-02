package node

import (
	"context"
	"fmt"

	"github.com/masa-finance/masa-oracle/node/types"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/masa-finance/masa-oracle/pkg/config"
)

const (
	masaPrefix = "/masa"
)

// ProtocolWithVersion returns a libp2p protocol ID string
// with the configured version and environment suffix.
func (node *OracleNode) protocolWithVersion(protocolName string) protocol.ID {
	if node.Options.Environment == "" {
		return protocol.ID(fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, node.Options.Version))
	}
	return protocol.ID(fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, node.Options.Version, node.Options.Environment))
}

// TopicWithVersion returns a topic string with the configured version
// and environment suffix.
func (node *OracleNode) topicWithVersion(protocolName string) string {
	if node.Options.Environment == "" {
		return fmt.Sprintf("%s/%s/%s", masaPrefix, protocolName, node.Options.Version)
	}
	return fmt.Sprintf("%s/%s/%s-%s", masaPrefix, protocolName, node.Options.Version, node.Options.Environment)
}

func (node *OracleNode) ProtocolStream(ctx context.Context, peerID peer.ID, protocolName string) (network.Stream, error) {
	return node.Host.NewStream(ctx, peerID, node.protocolWithVersion(protocolName))
}

// SubscribeToTopics handles the subscription to various topics for an OracleNode.
// It subscribes the node to the NodeGossipTopic, AdTopic, and PublicKeyTopic.
// Each subscription is managed through the node's PubSubManager, which orchestrates the message passing for these topics.
// Errors during subscription are logged and returned, halting the process to ensure the node's correct setup before operation.
func (node *OracleNode) subscribeToTopics() error {
	for _, handler := range node.Options.PubSubHandles {
		if err := node.SubscribeTopic(handler.ProtocolName, handler.Handler, handler.IncludeSelf); err != nil {
			return err
		}
	}

	// Subscribe to NodeGossipTopic to participate in the network's gossip protocol.
	if err := node.SubscribeTopic(config.NodeGossipTopic, node.NodeTracker, false); err != nil {
		return err
	}

	return nil
}

func (node *OracleNode) PublishTopic(protocolName string, data []byte) error {
	return node.PubSubManager.Publish(node.topicWithVersion(protocolName), data)
}

func (node *OracleNode) PublishTopicMessage(protocolName string, data string) error {
	return node.PubSubManager.PublishMessage(node.topicWithVersion(protocolName), data)
}

func (node *OracleNode) SubscribeTopic(protocolName string, handler types.SubscriptionHandler, includeSelf bool) error {
	return node.PubSubManager.AddSubscription(node.topicWithVersion(protocolName), handler, includeSelf)
}

func (node *OracleNode) Subscribe(protocolName string, handler types.SubscriptionHandler) error {
	return node.PubSubManager.Subscribe(node.topicWithVersion(protocolName), handler)
}
