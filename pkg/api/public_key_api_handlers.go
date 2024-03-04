package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/masa-finance/masa-oracle/pkg/crypto"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

func (api *API) PublishPublicKeyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Node or PubSubManager is not initialized"})
			return
		}

		keyManager := crypto.KeyManagerInstance()

		// Set the data to be signed as the signer's Peer ID
		data := []byte(api.Node.Host.ID().String())

		// Sign the data using the private key
		signature, err := consensus.SignData(keyManager.Libp2pPrivKey, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to sign data: %v", err)})
			return
		}

		// Serialize the public key message
		msg := pubsub.PublicKeyMessage{
			PublicKey: keyManager.HexPubKey,
			Signature: hex.EncodeToString(signature),
			Data:      string(data),
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal public key message"})
			return
		}

		// Publish the public key using its string representation, data, and signature
		publicKeyTopic := config.TopicWithVersion(config.PublicKeyTopic)
		if err := api.Node.PubSubManager.Publish(publicKeyTopic, msgBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Public key published successfully"})
	}
}

func (api *API) GetPublicKeysHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "An unexpected error occurred.",
			})
			return
		}

		// Use the PublicKeyTopic constant from the masa package
		handler, err := api.Node.PubSubManager.GetHandler(string(config.ProtocolWithVersion(config.PublicKeyTopic)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		publicKeyHandler, ok := handler.(*pubsub.PublicKeySubscriptionHandler)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "handler is not of type PublicKeySubscriptionHandler"})
			return
		}

		publicKeys := publicKeyHandler.GetPublicKeys()
		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"publicKeys": publicKeys,
		})
	}
}
