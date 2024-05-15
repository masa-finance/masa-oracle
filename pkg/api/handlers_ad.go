package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/masa-finance/masa-oracle/pkg/ad"
	"github.com/masa-finance/masa-oracle/pkg/config"
)

// PostAd handles posting a new ad. It requires the node to be staked,
// publishes the ad to the ad topic, and returns a success response.

func (api *API) PostAd() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		// Must be staked to publish to ad topic
		if !api.Node.IsStaked {
			c.JSON(http.StatusPreconditionRequired, gin.H{"error": "node must be staked to be an ad publisher"})
			return
		}

		if err := api.Node.PubSubManager.Publish(config.TopicWithVersion(config.AdTopic), bodyBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Ad published"})
	}
}

// SubscribeToAds handles subscribing to the ad topic. It adds a subscription
// using the AdSubscriptionHandler and returns a success response.

func (api *API) SubscribeToAds() gin.HandlerFunc {
	handler := &ad.SubscriptionHandler{}
	return func(c *gin.Context) {
		err := api.Node.PubSubManager.AddSubscription(config.TopicWithVersion(config.AdTopic), handler, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Subscribed to get ads"})
	}
}

// GetAds handles getting the collected ads. It responds with the ads collected
// by the AdSubscriptionHandler if any, otherwise responds with a message that
// there are no ads.

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
