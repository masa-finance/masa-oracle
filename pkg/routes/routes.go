package routes

import (
	"github.com/gin-gonic/gin"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
)

func SetupRoutes(node *masa.OracleNode) *gin.Engine {
	router := gin.Default()

	api := api.NewAPI(node)

	router.GET("/peers", api.GetPeersHandler())
	router.GET("/peerAddresses", api.GetPeerAddresses())

	router.POST("/ads", api.PostAd())
	router.GET("/ads", api.GetAds())
	router.POST("/subscribeToAds", api.SubscribeToAds())

	router.GET("/nodeData", api.GetNodeDataHandler())
	router.GET("/nodes/:peerID", api.GetNodeHandler())

	return router
}
