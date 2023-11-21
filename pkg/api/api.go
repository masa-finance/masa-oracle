package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	masa "github.com/masa-finance/masa-oracle/pkg"
)

// Handler for the /peers endpoint
func GetPeersHandler(node *masa.OracleNode) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the routing table
		routingTable := node.DHT.RoutingTable()

		// Get the list of peers from the routing table
		peers := routingTable.ListPeers()

		c.JSON(http.StatusOK, gin.H{"peers": peers})
	}
}
