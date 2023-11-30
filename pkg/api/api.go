package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/ad"
)

type API struct {
	Node *masa.OracleNode
}

func NewAPI(node *masa.OracleNode) *API {
	return &API{Node: node}
}

func (api *API) GetPeersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.DHT == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "An unexpected error occurred.",
			})
			return
		}

		routingTable := api.Node.DHT.RoutingTable()
		peers := routingTable.ListPeers()

		// Create a slice to hold the data
		data := make([]map[string]interface{}, len(peers))

		// Populate the data slice
		for i, peer := range peers {
			data[i] = map[string]interface{}{
				"peer": peer.String(),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"data":       data,
			"totalCount": len(peers),
		})
	}
}

func (api *API) GetPeerAddresses() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.Host == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "An unexpected error occurred.",
			})
			return
		}

		peers := api.Node.Host.Network().Peers()
		peerAddresses := make(map[string][]string)

		// Create a slice to hold the data
		data := make([]map[string]interface{}, len(peers))

		for i, peer := range peers {
			conns := api.Node.Host.Network().ConnsToPeer(peer)
			for _, conn := range conns {
				addr := conn.RemoteMultiaddr()
				peerAddresses[peer.String()] = append(peerAddresses[peer.String()], addr.String())
			}

			data[i] = map[string]interface{}{
				"peer":        peer.String(),
				"peerAddress": peerAddresses[peer.String()],
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"data":       data,
			"totalCount": len(peers),
		})
	}
}

// api.go

// ... (existing code)

func (api *API) PostAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newAd ad.Ad
		if err := c.ShouldBindJSON(&newAd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := api.Node.PublishAd(newAd); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Ad published"})
	}
}

func (api *API) GetAds() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Node.Ads)
	}
}

func (api *API) SubscribeToAds() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := api.Node.SubscribeToAds()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Subscribed to ad topic"})
	}
}
