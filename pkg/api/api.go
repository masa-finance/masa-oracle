package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	masa "github.com/masa-finance/masa-oracle/pkg"
)

// Handler for the /peers endpoint
func GetPeersHandler(node *masa.OracleNode) gin.HandlerFunc {
	return func(c *gin.Context) {
		peers := node.DHT.RoutingTable().Peers()
		c.JSON(http.StatusOK, gin.H{"peers": peers})
	}
}
