package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

// CreateNewTopicHandler creates a new topic with a given name and subscribes a handler to it.
func (api *API) CreateNewTopicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			TopicName string `json:"topicName"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Node or PubSubManager is not initialized"})
			return
		}

		// Initialize a TopicHandler for managing messages from the new topic.
		topicHandler := pubsub.NewTopicHandler()

		// Use the AddSubscription method to create the new topic and subscribe the TopicHandler to it.
		if err := api.Node.PubSubManager.AddSubscription(masa.TopicWithVersion(request.TopicName), topicHandler); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "New topic created and subscribed successfully"})
	}
}

// PostMessageToTopicHandler allows posting a message to a specified topic.
func (api *API) PostToTopicHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			TopicName string `json:"topicName"`
			Message   string `json:"message"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Node or PubSubManager is not initialized"})
			return
		}

		// Publish the message to the specified topic.
		if err := api.Node.PubSubManager.PublishMessage(masa.TopicWithVersion(request.TopicName), request.Message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Message posted to topic successfully"})
	}
}
