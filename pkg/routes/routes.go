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

	return router
}
