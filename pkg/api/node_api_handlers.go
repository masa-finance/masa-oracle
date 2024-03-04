package api

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

func (api *API) GetNodeDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageNbr, err := GetPathInt(c, "pageNbr")
		if err != nil {
			pageNbr = 0
		}
		pageSize, err := GetPathInt(c, "pageSize")
		if err != nil {
			pageSize = masa.PageSize
		}

		if api.Node == nil || api.Node.DHT == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "An unexpected error occurred.",
			})
			return
		}
		allNodeData := api.Node.NodeTracker.GetAllNodeData()
		totalRecords := len(allNodeData)
		totalPages := int(math.Ceil(float64(totalRecords) / masa.PageSize))

		startIndex := pageNbr * pageSize
		endIndex := startIndex + pageSize
		if endIndex > totalRecords {
			endIndex = totalRecords
		}
		nodeDataPage := masa.NodeDataPage{
			Data:         allNodeData[startIndex:endIndex],
			PageNumber:   pageNbr,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
		}
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"data":         nodeDataPage.Data,
			"pageNbr":      nodeDataPage.PageNumber,
			"total":        nodeDataPage.TotalRecords,
			"totalRecords": nodeDataPage.TotalRecords,
		})
	}
}

func (api *API) GetNodeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		peerID := c.Param("peerID") // Get the peer ID from the URL parameters
		if api.Node == nil || api.Node.NodeTracker == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "An unexpected error occurred.",
			})
			return
		}
		nodeData := api.Node.NodeTracker.GetNodeData(peerID)
		if nodeData == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Node not found",
			})
			return
		}
		nd := *nodeData
		nd.CurrentUptime = nodeData.GetCurrentUptime()
		nd.AccumulatedUptime = nodeData.GetAccumulatedUptime()
		nd.CurrentUptimeStr = pubsub.PrettyDuration(nd.CurrentUptime)
		nd.AccumulatedUptimeStr = pubsub.PrettyDuration(nd.AccumulatedUptime)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    nd,
		})
	}
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
