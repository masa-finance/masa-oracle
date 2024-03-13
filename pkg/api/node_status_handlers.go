package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/nodestatus"
	"github.com/sirupsen/logrus"
	"net/http"
)

// PostNodeStatusHandler allows posting a message to the NodeStatus Topic
func (api *API) PostNodeStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// WIP
		var nodeStatus nodestatus.NodeStatus

		if err := c.BindJSON(&nodeStatus); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Node or PubSubManager is not initialized"})
			return
		}

		jsonData, _ := json.Marshal(nodeStatus)
		logrus.Printf("jsonData %s", jsonData)

		// Publish the message to the specified topic.
		if err := api.Node.PubSubManager.Publish(config.TopicWithVersion(config.NodeStatusTopic), jsonData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Message posted to topic successfully"})
	}
}
