package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	"github.com/masa-finance/masa-oracle/pkg/twitter"
)

// SearchTweetsRequest remains unchanged
type SearchTweetsRequest struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

// SearchTweetsAndAnalyzeSentiment method adjusted to match the pattern
func (api *API) SearchTweetsAndAnalyzeSentiment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			Query string `json:"query"`
			Count int    `json:"count"`
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

		tweets, err := twitter.ScrapeTweetsByQuery(reqBody.Query, reqBody.Count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tweets"})
			return
		}

		_, sentimentSummary, err := llmbridge.AnalyzeSentiment(tweets)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze tweets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sentiment": sentimentSummary})
	}
}
