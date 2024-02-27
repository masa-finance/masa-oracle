package consensus

import (
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

// PublicKeyMessage represents the structure of the public key messages.
type PublicKeyMessage struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Data      string `json:"data"`
}

// PublicKeySubscriptionHandler handles incoming messages on public key topics.
type PublicKeySubscriptionHandler struct {
}

// HandleMessage processes messages received on the public key topic.
func (h *PublicKeySubscriptionHandler) HandleMessage(m *pubsub.Message) {
	var message PublicKeyMessage
	if err := json.Unmarshal(m.Data, &message); err != nil {
		logrus.WithError(err).Error("Failed to unmarshal public key message")
		return
	}
	logrus.Infof("Received public key: %s", message.PublicKey)

	// Future implementation will include verification and other logic

	// Removed return nil as the function does not return any value
}
