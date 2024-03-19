package api

import (
	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	"github.com/masa-finance/masa-oracle/pkg/twitter"
	"net/http"
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

		tweets, err := twitter.Scrape(reqBody.Query, reqBody.Count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tweets"})
			return
		}

		sentimentRequest, sentimentSummary, err := llmbridge.AnalyzeSentiment(tweets)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze tweets"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"sentimentRequest": sentimentRequest, "sentiment": sentimentSummary})
	}
}
