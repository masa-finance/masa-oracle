package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/twitter"
)

// SearchTweetsRequest remains unchanged
type SearchTweetsRequest struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

// SearchTweetsAndAnalyzeSentiment method adjusted to match the pattern
// Models Supported:
//
//	"claude-3-opus-20240229"
//	"claude-3-sonnet-20240229"
//	"claude-3-haiku-20240307"
//	"gpt-4"
//	"gpt-4-turbo-preview"
//	"gpt-3.5-turbo"
func (api *API) SearchTweetsAndAnalyzeSentiment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			Query string `json:"query"`
			Count int    `json:"count"`
			Model string `json:"model"`
		}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if reqBody.Query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is missing"})
			return
		}
		if reqBody.Count <= 0 {
			reqBody.Count = 50 // Default count
		}
		if reqBody.Model == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Model parameter is missing. Available models are claude-3-opus-20240229, claude-3-sonnet-20240229, claude-3-haiku-20240307, gpt-4, gpt-4-turbo-preview, gpt-3.5-turbo"})
			return
		}

		// Testing scrape using actor engine
		sentimentSummary, err := twitter.ScrapeTweetsUsingActors(reqBody.Query, reqBody.Count, reqBody.Model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tweets and analyze sentiment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sentiment": sentimentSummary})
	}
}
