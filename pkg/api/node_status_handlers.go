package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/nodestatus"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
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

func (api *API) NodeStatusPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodeData := api.Node.NodeTracker.GetNodeData(api.Node.Host.ID().String())
		if nodeData == nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Name":        "Masa Status Page",
				"PeerID":      api.Node.Host.ID().String(),
				"IsStaked":    false,
				"FirstJoined": time.Now().Format("2006-01-02 15:04:05"),
				"LastJoined":  time.Now().Format("2006-01-02 15:04:05"),
			})
			return
		}
		totalUpTimeInNs := nodeData.GetAccumulatedUptime()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Name":        "Masa Status Page",
			"PeerID":      nodeData.PeerId.String(),
			"IsStaked":    nodeData.IsStaked,
			"FirstJoined": time.Now().Add(-totalUpTimeInNs).Format("2006-01-02 15:04:05"),
			"LastJoined":  nodeData.LastJoined.Format("2006-01-02 15:04:05"),
		})
	}
}
