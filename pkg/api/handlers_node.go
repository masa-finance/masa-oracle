package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/sirupsen/logrus"

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
		peerID := c.Param("peerid") // Get the peer ID from the URL parameters
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
		peers := api.Node.NodeTracker.GetAllNodeData()

		// Create a slice to hold the data
		data := make([]map[string]interface{}, len(peers))

		// Populate the data slice
		for i, peer := range peers {
			data[i] = map[string]interface{}{
				"peerId": peer.PeerId.String(),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"data":       data,
			"totalCount": len(peers),
		})
	}
}

// GetPeerAddresses handles GET requests to retrieve the list of peer
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

		peers := api.Node.NodeTracker.GetAllNodeData()

		// Create a slice to hold the data
		data := make([]map[string]interface{}, len(peers))

		for i, peer := range peers {
			data[i] = map[string]interface{}{
				"peerId":      peer.PeerId.String(),
				"peerAddress": peer.Multiaddrs[0].String(),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"data":       data,
			"totalCount": len(peers),
		})
	}
}

// PublishPublicKeyHandler handles the /publickey endpoint. It retrieves the node's
// public key, signs the public key with the private key, creates a public key
// message with the key info, signs it, and publishes it to the public key topic.
// This allows other nodes to obtain this node's public key.
func (api *API) PublishPublicKeyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Node or PubSubManager is not initialized"})
			return
		}

		keyManager := masacrypto.KeyManagerInstance()

		// Set the data to be signed as the signer's Peer ID
		data := []byte(api.Node.Host.ID().String())

		// Sign the data using the private key
		signature, err := consensus.SignData(keyManager.Libp2pPrivKey, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to sign data: %v", err)})
			return
		}

		// Serialize the public key message
		msg := pubsub.PublicKeyMessage{
			PublicKey: keyManager.HexPubKey,
			Signature: hex.EncodeToString(signature),
			Data:      string(data),
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal public key message"})
			return
		}

		// Publish the public key using its string representation, data, and signature
		publicKeyTopic := config.TopicWithVersion(config.PublicKeyTopic)
		if err := api.Node.PubSubManager.Publish(publicKeyTopic, msgBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Public key published successfully"})
	}
}

// GetPublicKeysHandler handles the endpoint to retrieve all known public keys.
// It gets the public key subscription handler from the PubSub manager,
// extracts the public keys, and returns them in the response.
func (api *API) GetPublicKeysHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "An unexpected error occurred.",
			})
			return
		}

		// Use the PublicKeyTopic constant from the masa package
		handler, err := api.Node.PubSubManager.GetHandler(string(config.ProtocolWithVersion(config.PublicKeyTopic)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		publicKeyHandler, ok := handler.(*pubsub.PublicKeySubscriptionHandler)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "handler is not of type PublicKeySubscriptionHandler"})
			return
		}

		publicKeys := publicKeyHandler.GetPublicKeys()
		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"publicKeys": publicKeys,
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
		nv := db.ReadData(api.Node, keyStr)
		err := json.Unmarshal(nv, &sharedData)
		if err != nil {
			if IsBase64(string(nv)) {
				decodedString, _ := base64.StdEncoding.DecodeString(string(nv))
				_ = json.Unmarshal(decodedString, &sharedData)
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": sharedData,
				})
				return
			} else {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"message": string(nv),
				})
				return
			}
		}
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
		err = db.WriteData(api.Node, keyStr, jsonData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": keyStr,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": keyStr,
		})
	}
}

// PostNodeStatusHandler allows posting a message to the Topic
func (api *API) PostNodeStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var nodeData pubsub.NodeData
		if err := c.BindJSON(&nodeData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Node or PubSubManager is not initialized"})
			return
		}

		jsonData, _ := json.Marshal(nodeData)
		logrus.Printf("jsonData %s", jsonData)

		// Publish the message to the specified topic.
		if err := api.Node.PubSubManager.Publish(config.TopicWithVersion(config.NodeGossipTopic), jsonData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Message posted to topic successfully"})
	}
}

// ChatPageHandler returns a gin.HandlerFunc that renders the chat page.
// It responds to HTTP GET requests by serving the "chat.html" template.
// The handler sets the HTTP status to 200 (OK) and provides an empty gin.H map
// to the template, which can be used to pass data if needed in the future.
func (api *API) ChatPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	}
}

// NodeStatusPageHandler handles HTTP requests to show the node status page.
// It retrieves the node data from the node tracker, formats it, and renders
// an HTML page displaying the node's status and uptime info.
func (api *API) NodeStatusPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		peers := api.Node.Host.Network().Peers()
		nodeData := api.Node.NodeTracker.GetNodeData(api.Node.Host.ID().String())
		if nodeData == nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"TotalPeers":       0,
				"Name":             "Masa Status Page",
				"PeerID":           api.Node.Host.ID().String(),
				"IsStaked":         false,
				"IsTwitterScraper": false,
				"IsDiscordScraper": false,
				"IsWebScraper":     false,
				"FirstJoined":      api.Node.FromUnixTime(time.Now().Unix()),
				"LastJoined":       api.Node.FromUnixTime(time.Now().Unix()),
				"CurrentUptime":    "0",
				"Rewards":          "Coming Soon!",
				"BytesScraped":     0,
			})
			return
		} else {
			nd := *nodeData
			nd.CurrentUptime = nodeData.GetCurrentUptime()
			nd.AccumulatedUptime = nodeData.GetAccumulatedUptime()
			nd.CurrentUptimeStr = pubsub.PrettyDuration(nd.CurrentUptime)
			nd.AccumulatedUptimeStr = pubsub.PrettyDuration(nd.AccumulatedUptime)

			sharedData := db.SharedData{}
			nv := db.ReadData(api.Node, api.Node.Host.ID().String())
			_ = json.Unmarshal(nv, &sharedData)
			bytesScraped, _ := strconv.Atoi(fmt.Sprintf("%v", sharedData["bytesScraped"]))
			c.HTML(http.StatusOK, "index.html", gin.H{
				"TotalPeers":       len(peers),
				"Name":             "Masa Status Page",
				"PeerID":           api.Node.Host.ID().String(),
				"IsStaked":         nd.IsStaked,
				"IsTwitterScraper": nd.IsTwitterScraper,
				"IsDiscordScraper": nd.IsDiscordScraper,
				"IsWebScraper":     nd.IsWebScraper,
				"FirstJoined":      api.Node.FromUnixTime(nd.FirstJoinedUnix),
				"LastJoined":       api.Node.FromUnixTime(nd.LastJoinedUnix),
				"CurrentUptime":    nd.CurrentUptimeStr,
				"TotalUptime":      nd.AccumulatedUptimeStr,
				"BytesScraped":     fmt.Sprintf("%.4f MB", float64(bytesScraped)/(1024*1024)),
			})
		}
	}
}

// GetNodeApiKey returns a gin.HandlerFunc that generates and returns a JWT token for the node.
// The JWT token is signed using the node's host ID as the secret key.
// On success, it returns the generated JWT token in a JSON response.
// On failure, it returns an appropriate error message and HTTP status code.
func (api *API) GetNodeApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := consensus.GenerateJWTToken(api.Node.Host.ID().String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": jwtToken,
		})
	}
}
