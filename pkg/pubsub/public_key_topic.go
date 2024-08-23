package pubsub

import (
	"encoding/hex"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/sirupsen/logrus"
)

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

// HandleMessage handles incoming public key messages, with verification and update logic.
func (handler *PublicKeySubscriptionHandler) HandleMessage(m *pubsub.Message) {
	logrus.Info("[+] Handling incoming public key message")
	var incomingMsg PublicKeyMessage
	if err := json.Unmarshal(m.Data, &incomingMsg); err != nil {
		logrus.WithError(err).Error("[-] Failed to unmarshal public key message")
		return
	}

	logrus.Infof("[+] Received public key message: %s", incomingMsg.PublicKey)

	// Proceed with verification and update logic as before
	for i, existingMsg := range handler.PublicKeys {
		if existingMsg.PublicKey == incomingMsg.PublicKey {
			logrus.Infof("[+] Verifying signature for public key: %s", incomingMsg.PublicKey)
			// Decode the public key from hexadecimal to bytes
			pubKeyBytes, err := hex.DecodeString(existingMsg.PublicKey)
			if err != nil {
				logrus.WithError(err).Error("[-] Failed to decode public key for verification")
				return
			}

			// Unmarshal the public key bytes into a libp2pCrypto.PubKey
			pubKey, err := libp2pCrypto.UnmarshalPublicKey(pubKeyBytes)
			if err != nil {
				logrus.WithError(err).Error("[-] Failed to unmarshal public key for verification")
				return
			}

			// Use the VerifySignature function from the consensus package
			isValid, err := consensus.VerifySignature(pubKey, []byte(incomingMsg.Data), incomingMsg.Signature)
			if err != nil || !isValid {
				logrus.WithError(err).Error("[-] Failed signature verification or signature is invalid")
				return // Do not update or add if the signature is invalid
			}

			// Signature is valid, update the existing message's data
			handler.PublicKeys[i].Data = incomingMsg.Data
			logrus.Infof("[+] Updated public key message: %s", incomingMsg.PublicKey)
			logrus.Info("[+] Data stored in the slice successfully.")
			return
		}
	}

	// If no public key is stored yet, add the new message
	if len(handler.PublicKeys) == 0 {
		handler.PublicKeys = append(handler.PublicKeys, incomingMsg)
		logrus.Infof("[+] Added new public key message: %s", incomingMsg.PublicKey)
		logrus.Info("[+] Data stored in the slice successfully.")
	}
}

// GetPublicKeys returns the list of PublicKeyMessages.
func (handler *PublicKeySubscriptionHandler) GetPublicKeys() []PublicKeyMessage {
	logrus.Info("[+] Retrieving stored public keys")
	return handler.PublicKeys
}
