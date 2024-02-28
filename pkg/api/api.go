package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/ad"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type API struct {
	Node *masa.OracleNode
}

func NewAPI(node *masa.OracleNode) *API {
	return &API{Node: node}
}

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

func (api *API) PostAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if !api.Node.IsStaked {
			c.JSON(http.StatusPreconditionRequired, gin.H{"error": "node must be staked to be an ad publisher"})
			return
		}

		if err := api.Node.PubSubManager.Publish(masa.TopicWithVersion(masa.AdTopic), bodyBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Ad published"})
	}
}

func (api *API) GetAds() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Directly access the AdSubscriptionHandler from the OracleNode
		if api.Node.AdSubscriptionHandler == nil || len(api.Node.AdSubscriptionHandler.Ads) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No ads"})
			return
		}
		// Respond with the ads collected by the AdSubscriptionHandler
		c.JSON(http.StatusOK, api.Node.AdSubscriptionHandler.Ads)
	}
}

func (api *API) SubscribeToAds() gin.HandlerFunc {
	handler := &ad.SubscriptionHandler{}
	return func(c *gin.Context) {
		err := api.Node.PubSubManager.AddSubscription(masa.TopicWithVersion(masa.AdTopic), handler)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Subscribed to get ads"})
	}
}
func GetPathInt(ctx *gin.Context, name string) (int, error) {
	val, ok := ctx.GetQuery(name)
	if !ok {
		return 0, errors.New(fmt.Sprintf("the value for path parameter %s empty or not specified", name))
	}
	return strconv.Atoi(val)
}

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
		handler, err := api.Node.PubSubManager.GetHandler(string(masa.ProtocolWithVersion(masa.PublicKeyTopic)))
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

func (api *API) PublishPublicKeyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Node or PubSubManager is not initialized"})
			return
		}

		keyManager := masa.KeyManager{}

		// Load the private key directly using the KeyManager from keys package
		privKey, err := keyManager.LoadPrivKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load private key: %v", err)})
			return
		}

		// Load the public key directly using the KeyManager from keys package
		pubKey, err := keyManager.LoadPubKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to load public key: %v", err)})
			return
		}

		// Convert the public key to a string representation
		pubKeyString, err := masa.PubKeyToString(pubKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to convert public key to string: %v", err)})
			return
		}

		// Set the data to be signed as the signer's Peer ID
		data := []byte(api.Node.Host.ID().String())

		// Sign the data using the private key
		signature, err := consensus.SignData(privKey, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to sign data: %v", err)})
			return
		}

		// Serialize the public key message
		msg := pubsub.PublicKeyMessage{
			PublicKey: pubKeyString,
			Signature: hex.EncodeToString(signature),
			Data:      string(data),
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal public key message"})
			return
		}

		// Publish the public key using its string representation, data, and signature
		publicKeyTopic := masa.TopicWithVersion(masa.PublicKeyTopic)
		if err := api.Node.PubSubManager.Publish(publicKeyTopic, msgBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Public key published successfully"})
	}
}

func (api *API) GetPublishedPublicKeysHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.Node == nil || api.Node.PubSubManager == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Node or PubSubManager is not initialized",
			})
			return
		}

		// Assuming you have a way to access your PublicKeyPublisher instance here
		publisher := api.Node.PubSubManager.PublicKeyPublisher

		messages := publisher.GetPublishedMessages()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    messages,
		})
	}
}
