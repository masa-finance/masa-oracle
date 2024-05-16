package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"

	"strings"

	"strconv"

	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
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
			string(config.Models.GPT4o),
			string(config.Models.GPT4Turbo),
			string(config.Models.GPT35Turbo),
			string(config.Models.LLama2),
			string(config.Models.LLama3),
			string(config.Models.Mistral),
			string(config.Models.Gemma),
			string(config.Models.Mixtral),
			string(config.Models.OpenChat),
			string(config.Models.NeuralChat),
			string(config.Models.CloudflareQwen15Chat),
			string(config.Models.CloudflareLlama27bChatFp16),
			string(config.Models.CloudflareLlama38bInstruct),
			string(config.Models.CloudflareMistral7bInstruct),
			string(config.Models.CloudflareMistral7bInstructV01),
			string(config.Models.CloudflareOpenchat35_0106),
			string(config.Models.CloudflareMicrosoftPhi2),
			string(config.Models.HuggingFaceGoogleGemma7bIt),
			string(config.Models.HuggingFaceNousresearchHermes2ProMistral7b),
			string(config.Models.HuggingFaceTheblokeLlama213bChatAwq),
			string(config.Models.HuggingFaceTheblokeNeuralChat7bV31Awq),
		}
		c.JSON(http.StatusOK, gin.H{"models": models})
	}
}

// SearchTweetsAndAnalyzeSentiment method adjusted to match the pattern
// Models Supported:
//
//	chose a model or use "all"
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

// SearchDiscordProfile returns a gin.HandlerFunc that processes a request to search for a Discord user profile.
// It expects a URL parameter "userID" representing the Discord user ID to search for.
// The handler validates the userID, ensuring it is provided.
// If the request is valid, it attempts to fetch the user's profile.
// On success, it returns the fetched profile information in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchDiscordProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID must be provided and valid"})
			return
		}

		// Assuming you have a way to access your bot token here. It might be stored in an environment variable or a config file.
		botToken := os.Getenv("DISCORD_BOT_TOKEN")

		profile, err := discord.GetUserProfile(userID, botToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Discord profile", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, profile)
	}
}

// SearchDiscordGuildMemberships returns a gin.HandlerFunc that processes a request to search for guild memberships of a Discord user.
// It expects a URL parameter "userID" representing the Discord user ID to search for.
// The handler validates the userID, ensuring it is provided.
// If the request is valid, it attempts to fetch the user's guild memberships.
// On success, it returns the fetched guild membership information in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchDiscordGuildMemberships() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID must be provided and valid"})
			return
		}

		botToken := os.Getenv("DISCORD_BOT_TOKEN")
		if botToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Bot token is not configured"})
			return
		}

		guildMemberships, err := discord.ListGuildMemberships(userID, botToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Discord guild memberships", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"guild_memberships": guildMemberships})
	}
}

// GetTwitterFollowersHandler returns a gin.HandlerFunc that retrieves the followers of a given Twitter user.
func (api *API) GetTwitterFollowersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username") // Assuming you're using a URL parameter for the username
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is missing"})
			return
		}

		// Extracting maxUsersNbr from query parameters, with a default value if not specified
		maxUsersNbrStr := c.DefaultQuery("maxUsersNbr", "20") // Default to 20 if not specified
		maxUsersNbr, err := strconv.Atoi(maxUsersNbrStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxUsersNbr parameter"})
			return
		}

		followers, err := twitter.ScrapeFollowersForProfile(username, maxUsersNbr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"followers": followers})
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

type LLMChat struct {
	Model    string `json:"model,omitempty"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages,omitempty"`
	Stream bool `json:"stream"`
}

// LocalLlmChat handles requests for chatting with AI models hosted by ollama.
// It expects a JSON request body with a structure formatted for the model. For example for Ollama:
//
//	{
//	    "model": "llama3",
//	    "messages": [
//	        {
//	            "role": "user",
//	            "content": "why is the sky blue?"
//	        }
//	    ],
//	    "stream": false
//	}
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
func (api *API) LocalLlmChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// we just want to proxy the request JSON directly to the endpoint we are calling.
		body := c.Request.Body
		var reqBody LLMChat
		if err := json.NewDecoder(body).Decode(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		reqBody.Model = strings.TrimPrefix(reqBody.Model, "ollama/")
		reqBody.Stream = false
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)

		llmRequest := make(map[string]string)
		llmRequest["request"] = "llm-chat"
		llmRequest["request_id"] = requestID
		llmRequest["body"] = string(bodyBytes)
		jsn, err := json.Marshal(llmRequest)
		if err != nil {
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
		if err := api.Node.PubSubManager.Publish(config.TopicWithVersion(config.WorkerTopic), jsn); err != nil {
			logrus.Errorf("%v", err)
			c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
			return
		}

		result := make(map[string]interface{})
		// Wait for the response
		select {
		case response := <-responseCh:
			err := json.Unmarshal(response, &result)
			if err != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, result)
		case <-time.After(30 * time.Second):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
		}
	}
}

// CfLlmChat handles the Cloudflare LLM chat requests.
// It reads the request body, appends a system message, and forwards the request to the configured LLM endpoint.
// The response from the LLM endpoint is then returned to the client.
//
//	{
//	    "model": "@cf/meta/llama-3-8b-instruct",
//	    "messages": [
//	        {
//	            "role": "user",
//	            "content": "why is the sky blue?"
//	        }
//	    ]
//	}
//
// Models
//
//	@cf/qwen/qwen1.5-0.5b-chat
//	@cf/meta/llama-2-7b-chat-fp16
//	@cf/meta/llama-3-8b-instruct
//	@cf/mistral/mistral-7b-instruct
//	@cf/mistral/mistral-7b-instruct-v0.1
//	@hf/google/gemma-7b-it
//	@hf/nousresearch/hermes-2-pro-mistral-7b
//	@hf/thebloke/llama-2-13b-chat-awq
//	@hf/thebloke/neural-chat-7b-v3-1-awq
//	@cf/openchat/openchat-3.5-0106
//	@cf/microsoft/phi-2
//
// @return gin.HandlerFunc - the handler function for the LLM chat requests.
func (api *API) CfLlmChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := c.Request.Body
		var reqBody LLMChat
		if err := json.NewDecoder(body).Decode(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cfUrl := config.GetInstance().LLMCfUrl
		if cfUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("missing env LLM_CF_URL")})
			return
		}
		uri := fmt.Sprintf("%s%s", cfUrl, reqBody.Model)
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
