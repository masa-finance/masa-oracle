package pubsub

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sirupsen/logrus"
)

// topicPublicKeyMap maps topics to their associated public keys.
var topicPublicKeyMap = make(map[string]string)

// PublicKeySubscriptionHandler handles incoming messages on public key topics.
type PublicKeySubscriptionHandler struct {
	PublicKeys  []PublicKeyMessage
	PubKeyTopic *pubsub.Topic
}

// PublicKeyMessage represents the structure of the public key messages.
type PublicKeyMessage struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Data      string `json:"data"`
}

// PublicKeyPublisher uses the existing Manager to publish public keys.
type PublicKeyPublisher struct {
	pubSubManager *Manager
	pubKey        libp2pCrypto.PubKey
}

// NewPublicKeyPublisher creates a new instance of PublicKeyPublisher.
func NewPublicKeyPublisher(manager *Manager, pubKey libp2pCrypto.PubKey) *PublicKeyPublisher {
	logrus.Info("Creating new PublicKeyPublisher")
	return &PublicKeyPublisher{
		pubSubManager: manager,
		pubKey:        pubKey,
	}
}

// PublishNodePublicKey publishes the node's public key to the designated topic.
func (p *PublicKeyPublisher) PublishNodePublicKey(publicKey string, data, signature []byte) error {
	topicName := "bootNodePublicKey"
	logrus.Infof("[PublicKeyPublisher] Publishing node's public key to topic: %s", topicName)

	// Ensure the topic exists or create it
	_, err := p.ensureTopic(topicName)
	if err != nil {
		logrus.WithError(err).Errorf("[PublicKeyPublisher] Failed to ensure topic '%s' exists", topicName)
		return err
	}

	// Check if a public key has already been published to the topic
	existingPubKey, exists := topicPublicKeyMap[topicName]

	if exists {
		logrus.Infof("[PublicKeyPublisher] Public key already published for topic: %s. Verifying signature.", topicName)
		// If a public key exists, verify the signature against the existing public key
		pubKeyBytes, err := hex.DecodeString(existingPubKey)
		if err != nil {
			logrus.WithError(err).Error("[PublicKeyPublisher] Failed to decode existing public key for verification")
			return err
		}
		pubKey, err := libp2pCrypto.UnmarshalPublicKey(pubKeyBytes)
		if err != nil {
			logrus.WithError(err).Error("[PublicKeyPublisher] Failed to unmarshal existing public key for verification")
			return err
		}
		isValid, err := pubKey.Verify(data, signature)
		if err != nil || !isValid {
			logrus.WithError(err).Error("[PublicKeyPublisher] Unauthorized: Failed signature verification or signature is invalid")
			return errors.New("unauthorized: only the owner of the public key can publish changes")
		}
		logrus.Info("[PublicKeyPublisher] Signature verified successfully for topic: ", topicName)
	} else {
		logrus.Infof("[PublicKeyPublisher] No existing public key for topic: %s. Proceeding with initial publication.", topicName)
		// If no public key is associated with the topic, this is the initial publication
		topicPublicKeyMap[topicName] = publicKey
	}

	// Serialize the public key message
	msg := PublicKeyMessage{
		PublicKey: publicKey,
		Signature: hex.EncodeToString(signature),
		Data:      string(data),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logrus.WithError(err).Error("[PublicKeyPublisher] Failed to marshal public key message")
		return errors.New("failed to marshal message")
	}

	// Use the existing Manager to publish the message
	logrus.Infof("[PublicKeyPublisher] Publishing serialized message to topic: %s", topicName)
	if err := p.pubSubManager.Publish(topicName, msgBytes); err != nil {
		return err
	}

	// Print the published data in the console
	logrus.Infof("[PublicKeyPublisher] Published data: PublicKey: %s, Signature: %s, Data: %s", msg.PublicKey, msg.Signature, msg.Data)
	return nil
}

// ensureTopic checks if a topic exists and creates it if not.
func (p *PublicKeyPublisher) ensureTopic(topicName string) (*pubsub.Topic, error) {
	// Check if the topic already exists
	if topic, exists := p.pubSubManager.topics[topicName]; exists {
		return topic, nil
	}

	// If the topic does not exist, attempt to create it
	topic, err := p.pubSubManager.createTopic(topicName)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// HandleMessage handles incoming public key messages, with verification and update logic.
func (handler *PublicKeySubscriptionHandler) HandleMessage(m *pubsub.Message) {
	var incomingMsg PublicKeyMessage
	if err := json.Unmarshal(m.Data, &incomingMsg); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal public key message")
		return
	}

	// Attempt to find an existing message with the same public key
	for i, existingMsg := range handler.PublicKeys {
		if existingMsg.PublicKey == incomingMsg.PublicKey {
			// Found an existing message, verify the signature
			pubKeyBytes, err := hex.DecodeString(existingMsg.PublicKey)
			if err != nil {
				logrus.WithError(err).Error("Failed to decode public key for verification")
				return
			}

			pubKey, err := libp2pCrypto.UnmarshalPublicKey(pubKeyBytes)
			if err != nil {
				logrus.WithError(err).Error("Failed to unmarshal public key for verification")
				return
			}

			isValid, err := pubKey.Verify([]byte(incomingMsg.Data), []byte(incomingMsg.Signature))
			if err != nil || !isValid {
				logrus.WithError(err).Error("Failed signature verification or signature is invalid")
				return // Do not update or add if the signature is invalid
			}

			// Signature is valid, update the existing message's data
			handler.PublicKeys[i].Data = incomingMsg.Data
			logrus.Infof("Updated public key message: %s", incomingMsg.PublicKey)
			return
		}
	}

	// No existing message with the same public key, add the new message
	handler.PublicKeys = append(handler.PublicKeys, incomingMsg)
	logrus.Infof("Added new public key message: %s", incomingMsg.PublicKey)
}