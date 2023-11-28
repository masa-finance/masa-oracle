package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/ad"
	"github.com/masa-finance/masa-oracle/pkg/api"
)

func SetupRoutes(node *masa.OracleNode) *gin.Engine {
	router := gin.Default()

	api := api.NewAPI(node)

	router.GET("/peers", api.GetPeersHandler())
	router.GET("/peerAddresses", api.GetPeerAddresses())

	router.POST("/ads", func(c *gin.Context) {
		var newAd ad.Ad
		if err := c.ShouldBindJSON(&newAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := node.PublishAd(newAd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Ad published"})
	})

	router.GET("/ads", func(c *gin.Context) {
		c.JSON(http.StatusOK, node.Ads)
	})

	// New route to subscribe to ad topic
	router.POST("/subscribeToAds", func(c *gin.Context) {
		err := node.SubscribeToAds()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Subscribed to ad topic"})
	})

	return router
}
