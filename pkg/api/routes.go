package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	masa "github.com/masa-finance/masa-oracle/pkg"
)

func SetupRoutes(node *masa.OracleNode) *gin.Engine {
	router := gin.Default()
	// @TODO need to add a Authorization Bearer methodology for api security
	// add cors middleware

	API := NewAPI(node)

	router.GET("/peers", API.GetPeersHandler())
	router.GET("/peerAddresses", API.GetPeerAddresses())

	router.POST("/ads", API.PostAd())
	router.GET("/ads", API.GetAds())
	router.POST("/subscribeToAds", API.SubscribeToAds())

	router.GET("/nodeData", API.GetNodeDataHandler())
	router.GET("/nodeData/:peerID", API.GetNodeHandler())

	router.GET("/publicKeys", API.GetPublicKeysHandler())
	router.POST("/publishPublicKey", API.PublishPublicKeyHandler())

	router.POST("/createTopic", API.CreateNewTopicHandler())
	router.POST("/postToTopic", API.PostToTopicHandler())

	router.GET("/dht", API.GetFromDHT())
	router.POST("/dht", API.PostToDHT())

	router.POST("/nodestatus", API.PostNodeStatusHandler())

	router.LoadHTMLGlob("pkg/api/templates/*.html")

	router.GET("/status", func(c *gin.Context) {
		nodeData := API.Node.NodeTracker.GetNodeData(node.Host.ID().String())
		totalUpTimeInNs := nodeData.GetAccumulatedUptime()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Name":        "Masa Status Page",
			"PeerID":      nodeData.PeerId.String(),
			"IsStaked":    nodeData.IsStaked,
			"FirstJoined": time.Now().Add(-totalUpTimeInNs).Format("2006-01-02 15:04:05"),
			"LastJoined":  nodeData.LastJoined.Format("2006-01-02 15:04:05"),
		})
	})

	return router
}
