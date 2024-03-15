package api

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/masa-finance/masa-oracle/pkg/db"

	"github.com/gin-gonic/gin"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

// GetNodeDataHandler handles GET requests to retrieve paginated node data from the node tracker.
// It parses the page number and page size from the request path, retrieves all node data from the
// node tracker, calculates pagination details like total pages based on page size, and returns a
// page of node data in the response.
func (api *API) GetNodeDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageNbr, err := GetPathInt(c, "pageNbr")
		if err != nil {
			pageNbr = 0
		}
		pageSize, err := GetPathInt(c, "pageSize")
		if err != nil {
			pageSize = config.PageSize
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
		totalPages := int(math.Ceil(float64(totalRecords) / config.PageSize))

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

// GetNodeHandler handles GET requests to retrieve node data for a specific peer ID.
// It extracts the peer ID from the request URL parameters, retrieves the node data
// from the node tracker, calculates additional uptime info, and returns the node
// data in the response.
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

// GetPeersHandler handles GET requests to retrieve the list of peer IDs
// from the DHT routing table. It retrieves the routing table from the
// node's DHT instance, extracts the peer IDs, and returns them in the
// response.
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

// GetPeerAddressesHandler handles GET requests to retrieve the list of peer
// addresses from the node's libp2p host network. It gets the list of connected
// peers, finds the multiaddrs for connections to each peer, and returns the
// peer IDs mapped to their addresses.
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

// GetFromDHT handles GET requests to retrieve data from the DHT
// given a key. It looks up the key in the DHT, unmarshals the
// value into a SharedData struct, and returns the data in the response.
func (api *API) GetFromDHT() gin.HandlerFunc {
	return func(c *gin.Context) {

		keyStr := c.Query("key")
		if len(keyStr) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "missing key param",
			})
			return
		}
		sharedData := db.SharedData{}
		nodeVal := db.ReadData(api.Node, "/db/"+keyStr)
		_ = json.Unmarshal(nodeVal, &sharedData)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": sharedData,
		})
	}
}

// PostToDHT handles POST requests to write data to the DHT.
// It expects a JSON body with "key" and "value" fields.
// The "key" is used to store the data in the DHT under /db/key.
// The "value" is marshalled to JSON and written to the DHT.
// Returns 200 OK on success with the key in the response.
// Returns 400 Bad Request on invalid request or JSON errors.
func (api *API) PostToDHT() gin.HandlerFunc {
	return func(c *gin.Context) {

		sharedData := db.SharedData{}
		if err := c.BindJSON(&sharedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "invalid request",
			})
			return
		}

		var keyStr = sharedData["key"].(string)
		jsonData, err := json.Marshal(sharedData["value"])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "invalid json",
			})
			return
		}
		success, err := db.WriteData(api.Node, "/db/"+keyStr, jsonData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": success,
				"message": keyStr,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": success,
			"message": keyStr,
		})
	}
}
