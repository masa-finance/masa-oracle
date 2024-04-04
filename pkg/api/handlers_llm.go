package api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/llmbridge"

	"github.com/masa-finance/masa-oracle/pkg/scraper"

	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/twitter"
)

// SearchTweetsRequest remains unchanged
type SearchTweetsRequest struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

// SearchTweets returns a gin.HandlerFunc that processes a request to search for tweets based on a query and count.
// It expects a JSON body with fields "query" (string) and "count" (int), representing the search query and the number of tweets to return, respectively.
// The handler validates the request body, ensuring the query is not empty and the count is positive.
// If the request is valid, it attempts to scrape tweets using the specified query and count.
// On success, it returns the scraped tweets in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchTweets() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody SearchTweetsRequest
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if reqBody.Query == "" || reqBody.Count <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query and count must be provided and valid"})
			return
		}

		tweets, err := twitter.ScrapeTweetsByQuery(reqBody.Query, reqBody.Count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scrape tweets", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tweets": tweets})
	}
}

// SearchTweetsAndAnalyzeSentiment method adjusted to match the pattern
// Models Supported:
//
//	"all"
//	"claude-3-opus-20240229"
//	"claude-3-sonnet-20240229"
//	"claude-3-haiku-20240307"
//	"gpt-4"
//	"gpt-4-turbo-preview"
//	"gpt-3.5-turbo"
func (api *API) SearchTweetsAndAnalyzeSentiment() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !api.Node.IsStaked {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Node has not staked and cannot participate"})
			return
		}
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

		var sentimentSummary string
		var err error

		if reqBody.Model == "all" {
			models := config.Models
			val := reflect.ValueOf(models)

			type ModelResult struct {
				Model     string `json:"model"`
				Sentiment string `json:"sentiment"`
				Duration  string `json:"duration"`
			}
			var results []ModelResult

			for i := 0; i < val.NumField(); i++ {
				model := val.Field(i).Interface().(config.ModelType)
				startTime := time.Now() // Start time measurement

				sentimentSummary, err := twitter.ScrapeTweetsUsingActors(reqBody.Query, reqBody.Count, string(model))
				duration := time.Since(startTime) // Calculate duration

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tweets and analyze sentiment for model " + string(model)})
					return
				}

				results = append(results, ModelResult{
					Model:     string(model),
					Sentiment: sentimentSummary,
					Duration:  duration.String(),
				})
			}

			// Return the results as JSON
			c.JSON(http.StatusOK, gin.H{"results": results})
			return
		} else {
			sentimentSummary, err = twitter.ScrapeTweetsUsingActors(reqBody.Query, reqBody.Count, reqBody.Model)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tweets and analyze sentiment"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"sentiment": sentimentSummary})
	}
}

func (api *API) SearchWebAndAnalyzeSentiment() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !api.Node.IsStaked {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Node has not staked and cannot participate"})
			return
		}
		var reqBody struct {
			Url   string `json:"url"`
			Depth int    `json:"depth"`
			Model string `json:"model"`
		}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if reqBody.Url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is missing"})
			return
		}
		if reqBody.Depth <= 0 {
			reqBody.Depth = 10 // Default count
		}
		if reqBody.Model == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Model parameter is missing. Available models are claude-3-opus-20240229, claude-3-sonnet-20240229, claude-3-haiku-20240307, gpt-4, gpt-4-turbo-preview, gpt-3.5-turbo"})
			return
		}

		var sentimentRequest, sentimentSummary string
		var err error

		if reqBody.Model == "all" {
			models := config.Models
			val := reflect.ValueOf(models)

			type ModelResult struct {
				Model     string `json:"model"`
				Sentiment string `json:"sentiment"`
				Duration  string `json:"duration"`
			}
			var results []ModelResult

			for i := 0; i < val.NumField(); i++ {
				model := val.Field(i).Interface().(config.ModelType)
				startTime := time.Now() // Start time measurement

				sentimentRequest, sentimentSummary, err = scraper.Collect([]string{reqBody.Url}, reqBody.Depth, reqBody.Model)
				j, _ := json.Marshal(sentimentSummary)
				sentimentSummary = string(j)

				duration := time.Since(startTime) // Calculate duration

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch web data and analyze sentiment for model " + string(model)})
					return
				}

				results = append(results, ModelResult{
					Model:     string(model),
					Sentiment: sentimentSummary,
					Duration:  duration.String(),
				})
			}

			// Return the results as JSON
			c.JSON(http.StatusOK, gin.H{"data": sentimentRequest, "sentiment": results})
			return
		} else {

			sentimentRequest, sentimentSummary, err = scraper.Collect([]string{reqBody.Url}, reqBody.Depth, reqBody.Model)
			j, _ := json.Marshal(sentimentSummary)
			sentimentSummary = llmbridge.SanitizeResponse(string(j))

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch web data and analyze sentiment"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"data": sentimentRequest, "sentiment": sentimentSummary})
	}
}
