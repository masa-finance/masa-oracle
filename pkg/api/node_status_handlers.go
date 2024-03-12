package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/nodestatus"
	"github.com/sirupsen/logrus"
	"net/http"
)

// WIP
// *** Store NodeStatus ***
//up := node.NodeTracker.GetNodeData(node.Host.ID().String())
//if up != nil {
//	totalUpTime := up.GetAccumulatedUptime()
//	status := db.NodeStatus{
//		PeerID:        node.Host.ID().String(),
//		IsStaked:      isStaked,
//		TotalUpTime:   totalUpTime,
//		FirstLaunched: time.Now().Add(-totalUpTime),
//		LastLaunched:  time.Now(),
//	}
//	jsonData, _ := json.Marshal(status)
//	logrus.Printf("jsonData %s", jsonData)
//	// str := fmt.Sprintf("%s", jsonData)
//	if err := node.PubSubManager.Publish(config.TopicWithVersion("nodeStatus"), jsonData); err != nil {
//		logrus.Errorf("PublishMessage %+v", err)
//	}
//
//	keyStr := node.Host.ID().String() // user ID for this nodes status key
//	time.Sleep(time.Second * 1)       // delay needed to wait for node to finish starting
//	success, er := db.WriteData(node, "/db/"+keyStr, jsonData)
//	if er != nil {
//		logrus.Errorf("Store NodeStatus err %+v", er)
//	}
//	logrus.Infof("Store NodeStatus %+v", success)
//}
// *** Store NodeStatus ***

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
