package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
)

// GetLLMModelsHandler returns a gin.HandlerFunc that retrieves the available LLM models.
// It does not expect any request parameters.
// The handler returns a JSON response containing an array of supported LLM model names.
func (api *API) GetLLMModelsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		models := []string{
			string(config.Models.ClaudeOpus),
			string(config.Models.ClaudeSonnet),
			string(config.Models.ClaudeHaiku),
			string(config.Models.GPT4),
			string(config.Models.GPT4Turbo),
			string(config.Models.GPT35Turbo),
			string(config.Models.LLama2),
			string(config.Models.LLama3),
			string(config.Models.Mistral),
			string(config.Models.Gemma),
			string(config.Models.Mixtral),
			string(config.Models.OpenChat),
			string(config.Models.NeuralChat),
		}
		c.JSON(http.StatusOK, gin.H{"models": models})

	}
}

// SearchTweetsAndAnalyzeSentiment method adjusted to match the pattern
// Models Supported:
//
//	"all"
//
// claude-3-opus-20240229
// claude-3-sonnet-20240229
// claude-3-haiku-20240307
// gpt-4
// gpt-4-turbo-preview
// gpt-3.5-turbo
// llama2
// llama3
// mistral
// gemma
// mixtral
// openchat
// neural-chat
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

				_, sentimentSummary, err = twitter.ScrapeTweetsForSentiment(reqBody.Query, reqBody.Count, string(model))
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
			_, sentimentSummary, err = twitter.ScrapeTweetsForSentiment(reqBody.Query, reqBody.Count, reqBody.Model)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tweets and analyze sentiment"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"sentiment": sentimentSummary})
	}
}

// SearchWebAndAnalyzeSentiment returns a gin.HandlerFunc that processes web search requests and performs sentiment analysis.
// It first validates the request body for required fields such as URL, Depth, and Model. If the Model is set to "all",
// it iterates through all available models to perform sentiment analysis on the web content fetched from the specified URL.
// The function responds with the sentiment analysis results in JSON format.
// Models Supported:
//
//	"all"
//	"claude-3-opus-20240229"
//	"claude-3-sonnet-20240229"
//	"claude-3-haiku-20240307"
//	"gpt-4"
//	"gpt-4-turbo-preview"
//	"gpt-3.5-turbo"
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

				_, sentimentSummary, err = web.ScrapeWebDataForSentiment([]string{reqBody.Url}, reqBody.Depth, reqBody.Model)
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
			c.JSON(http.StatusOK, gin.H{"sentiment": results})
			return
		} else {

			_, sentimentSummary, err = web.ScrapeWebDataForSentiment([]string{reqBody.Url}, reqBody.Depth, reqBody.Model)
			j, _ := json.Marshal(sentimentSummary)
			sentimentSummary = llmbridge.SanitizeResponse(string(j))

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch web data and analyze sentiment"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"sentiment": sentimentSummary})
	}
}

// SearchTweetsProfile returns a gin.HandlerFunc that processes a request to search for tweets from a specific user profile.
// It expects a URL parameter "username" representing the Twitter username to search for.
// The handler validates the username, ensuring it is provided.
// If the request is valid, it attempts to scrape the user's profile and tweets.
// On success, it returns the scraped profile information in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchTweetsProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be provided and valid"})
			return
		}

		profile, err := twitter.ScrapeTweetsProfile(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get twitter profile", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tweets": profile})
	}
}

// SearchTweetsRecent returns a gin.HandlerFunc that processes a request to search for tweets based on a query and count.
// It expects a JSON body with fields "query" (string) and "count" (int), representing the search query and the number of tweets to return, respectively.
// The handler validates the request body, ensuring the query is not empty and the count is positive.
// If the request is valid, it attempts to scrape tweets using the specified query and count.
// On success, it returns the scraped tweets in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchTweetsRecent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			Query string `json:"query"`
			Count int    `json:"count"`
		}

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

// SearchTweetsTrends returns a gin.HandlerFunc that processes a request to search for trending tweets.
// It does not expect any request parameters.
// The handler attempts to scrape trending tweets using the ScrapeTweetsByTrends function.
// On success, it returns the scraped tweets in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchTweetsTrends() gin.HandlerFunc {
	return func(c *gin.Context) {

		tweets, err := twitter.ScrapeTweetsByTrends()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scrape tweets", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tweets": tweets})
	}
}

// WebData returns a gin.HandlerFunc that processes web scraping requests.
// It expects a JSON body with fields "url" (string) and "depth" (int), representing the URL to scrape and the depth of the scrape, respectively.
// The handler validates the request body, ensuring the URL is not empty and the depth is positive.
// If the node has not staked, it returns an error indicating the node cannot participate.
// On a valid request, it attempts to scrape web data using the specified URL and depth.
// On success, it returns the scraped data in a sanitized JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) WebData() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !api.Node.IsStaked {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Node has not staked and cannot participate"})
			return
		}
		var reqBody struct {
			Url   string `json:"url"`
			Depth int    `json:"depth"`
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

		collectedData, err := web.ScrapeWebData([]string{reqBody.Url}, reqBody.Depth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not scrape web data"})
			return
		}
		sanitizedData := llmbridge.SanitizeResponse(collectedData)
		c.JSON(http.StatusOK, gin.H{"data": sanitizedData})
	}
}

// LlmChat handles requests for chatting with AI models hosted by ollama.
// It expects a JSON request body with a structure formatted for the model. For example for Ollama:
//
//		{
//		    "model": "llama3",
//		    "messages": [
//		        {
//		            "role": "user",
//		            "content": "why is the sky blue?"
//		        }
//		    ],
//		    "stream": false
//		}
//
//		{
//	 	"query": "I just successfully staked my $MASA Tokens. It’s super easy. Choose your lock-up time and earn up to 25% APY in MASA rewards. @getmasafi"
//		}
//
// This function acts as a proxy, forwarding the request to hosted models and returning the proprietary structured response.
// This is intended to be compatible with code that is looking to leverage a common payload for LLMs that is based on
// the model name/type
// So if it is an Ollama request it is the responsibility of the caller to properly format their payload to conform
// to the required structure similar to above.
//
// See:
// https://platform.openai.com/docs/api-reference/authentication
// https://docs.anthropic.com/claude/reference/complete_post
// https://github.com/ollama/ollama/blob/main/docs/api.md
// note: Ollama recently added support for the OpenAI structure which can simplify integrating it.
func (api *API) LlmChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// we just want to proxy the request JSON directly to the endpoint we are calling.
		body := c.Request.Body
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		// Process the message
		uri := config.GetInstance().LLMChatUrl
		if uri == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("missing env LLM_CHAT_URL")})
			return
		}
		resp, err := http.Post(uri, "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logrus.Error(err)
			}
		}(resp.Body)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		}
		var payload map[string]interface{}
		err = json.Unmarshal(respBody, &payload)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		}
		// Return the response
		c.JSON(http.StatusOK, payload)
	}
}

func (api *API) LlmChatCf() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := c.Request.Body
		var reqBody struct {
			Query     string `json:"query,omitempty"`
			MaxTokens int    `json:"max_tokens"`
			Messages  []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(body).Decode(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		reqBody.MaxTokens = 2048
		reqBody.Messages = append(reqBody.Messages, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			Role: "system",
			// use this for sentiment analysis
			// Content: "Please perform a sentiment analysis on the following tweets, using an unbiased approach. Sentiment analysis involves identifying and categorizing opinions expressed in text, particularly to determine whether the writer`s attitude towards a particular topic, product, etc., is positive, negative, or neutral. After analyzing, please provide a summary of the overall sentiment expressed in these tweets, including the proportion of positive, negative, and neutral sentiments if applicable.",
			// use this for standard chat
			Content: os.Getenv("PROMPT"),
		})

		reqBody.Messages = append(reqBody.Messages, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			Role:    "user",
			Content: reqBody.Query,
		})

		reqBody.Query = ""

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		model := "@cf/meta/llama-3-8b-instruct"
		uri := fmt.Sprintf("%s%s", os.Getenv("LLM_CF_URL"), model)
		if uri == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("missing env LLM_CF_URL")})
			return
		}
		req, err := http.NewRequest("POST", uri, bytes.NewReader(bodyBytes))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		bearer := fmt.Sprintf("Bearer %s", os.Getenv("LLM_CF_TOKEN"))
		req.Header.Set("Authorization", bearer)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var payload map[string]interface{}
		err = json.Unmarshal(respBody, &payload)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, payload)
	}
}
