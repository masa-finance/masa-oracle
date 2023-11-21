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

// Handler for the /peerAddresses endpoint
func GetPeerAddresses(node *masa.OracleNode) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the list of peers from the network
		peers := node.Host.Network().Peers()

		// Create a map to hold the peer addresses
		peerAddresses := make(map[string][]string)

		// Iterate over the peers and get their addresses
		for _, peer := range peers {
			// Get the network connection to the peer
			conns := node.Host.Network().ConnsToPeer(peer)

			// Iterate over the connections and get the remote multiaddresses
			for _, conn := range conns {
				// Get the remote multiaddress of the connection
				addr := conn.RemoteMultiaddr()

				// Add the address to the map
				peerAddresses[peer.Pretty()] = append(peerAddresses[peer.Pretty()], addr.String())
			}
		}

		c.JSON(http.StatusOK, gin.H{"peerAddresses": peerAddresses})
	}
}
