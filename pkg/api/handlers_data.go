package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/masa-finance/masa-oracle/pkg/workers"
	"github.com/sirupsen/logrus"

	"strings"

	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

type LLMChat struct {
	Model    string `json:"model,omitempty"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages,omitempty"`
	Stream bool `json:"stream"`
}

// publishWorkRequest sends a work request to the PubSubManager for processing by a worker.
// It marshals the request details into JSON and publishes it to the configured topic.
//
// Parameters:
// - api: The API instance containing the Node and PubSubManager.
// - requestID: A unique identifier for the request.
// - request: The type of work to be performed by the worker.
// - bodyBytes: The request body in byte slice format.
//
// Returns:
// - error: An error object if the request could not be published, otherwise nil.
func publishWorkRequest(api *API, requestID string, request workers.WorkerType, bodyBytes []byte) error {
	workRequest := map[string]string{
		"request":    string(request),
		"request_id": requestID,
		"body":       string(bodyBytes),
	}
	jsn, err := json.Marshal(workRequest)
	if err != nil {
		return err
	}
	return api.Node.PubSubManager.Publish(config.TopicWithVersion(config.WorkerTopic), jsn)
}

// handleWorkResponse processes the response from a worker and sends it back to the client.
// It listens on the provided response channel for a response or a timeout signal.
// If a response is received within the timeout period, it unmarshals the JSON response and sends it back to the client.
// If no response is received within the timeout period, it sends a timeout error to the client.
//
// Parameters:
// - c: The gin.Context object, which provides the context for the HTTP request.
// - responseCh: A channel that receives the worker's response as a byte slice.
func handleWorkResponse(c *gin.Context, responseCh chan []byte) {
	for {
		select {
		case response := <-responseCh:
			var result map[string]interface{}
			if err := json.Unmarshal(response, &result); err != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, result)
			return
		case <-time.After(60 * time.Second):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out"})
			return
		case <-c.Done():
			return
		}
	}
}

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
			string(config.Models.CloudflareOpenchat350106),
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

		// worker handler implementation
		bodyBytes, wErr := json.Marshal(reqBody)
		if wErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": wErr.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		wErr = publishWorkRequest(api, requestID, workers.WORKER.TwitterSentiment, bodyBytes)
		if wErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": wErr.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
	}
}

// SearchWebAndAnalyzeSentiment returns a gin.HandlerFunc that processes web search requests and performs sentiment analysis.
// It first validates the request body for required fields such as URL, Depth, and Model. If the Model is set to "all",
// it iterates through all available models to perform sentiment analysis on the web content fetched from the specified URL.
// The function responds with the sentiment analysis results in JSON format.// Models Supported:
// Models Supported:
//
//	chose a model or use "all"
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

		// worker handler implementation
		bodyBytes, wErr := json.Marshal(reqBody)
		if wErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": wErr.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		wErr = publishWorkRequest(api, requestID, workers.WORKER.WebSentiment, bodyBytes)
		if wErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": wErr.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
	}
}

// SearchTweetsProfile returns a gin.HandlerFunc that processes a request to search for tweets from a specific user profile.
// It expects a URL parameter "username" representing the Twitter username to search for.
// The handler validates the username, ensuring it is provided.
// If the request is valid, it attempts to scrape the user's profile and tweets.
// On success, it returns the scraped profile information in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchTweetsProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			Username string `json:"username"`
		}
		if c.Param("username") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be provided and valid"})
			return
		}
		reqBody.Username = c.Param("username")

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		err = publishWorkRequest(api, requestID, workers.WORKER.TwitterProfile, bodyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
	}
}

// SearchDiscordProfile returns a gin.HandlerFunc that processes a request to search for a Discord user profile.
// It expects a URL parameter "userID" representing the Discord user ID to search for.
// The handler validates the userID, ensuring it is provided.
// If the request is valid, it attempts to fetch the user's profile.
// On success, it returns the fetched profile information in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchDiscordProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			UserID   string `json:"userID"`
			BotToken string `json:"botToken"`
		}

		reqBody.UserID = c.Param("userID")
		reqBody.BotToken = os.Getenv("DISCORD_BOT_TOKEN")

		if reqBody.UserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID must be provided and valid"})
			return
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		err = publishWorkRequest(api, requestID, workers.WORKER.Discord, bodyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
	}
}

// SearchTwitterFollowers returns a gin.HandlerFunc that retrieves the followers of a given Twitter user.
func (api *API) SearchTwitterFollowers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			Username string `json:"username"`
			Count    int    `json:"count"`
		}

		username := c.Param("username") // Assuming you're using a URL parameter for the username
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is missing"})
			return
		}
		reqBody.Username = username
		if reqBody.Count == 0 {
			reqBody.Count = 20
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		err = publishWorkRequest(api, requestID, workers.WORKER.TwitterFollowers, bodyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation

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

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		err = publishWorkRequest(api, requestID, workers.WORKER.Twitter, bodyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
	}
}

// SearchTweetsTrends returns a gin.HandlerFunc that processes a request to search for trending tweets.
// It does not expect any request parameters.
// The handler attempts to scrape trending tweets using the ScrapeTweetsByTrends function.
// On success, it returns the scraped tweets in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchTweetsTrends() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Info(c)
		// worker handler implementation
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		err := publishWorkRequest(api, requestID, workers.WORKER.TwitterTrends, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
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

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		err = publishWorkRequest(api, requestID, workers.WORKER.Web, bodyBytes)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
	}
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
		body := c.Request.Body
		var reqBody LLMChat
		if err := json.NewDecoder(body).Decode(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		reqBody.Model = strings.TrimPrefix(reqBody.Model, "ollama/")
		reqBody.Stream = false

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		err = publishWorkRequest(api, requestID, workers.WORKER.LLMChat, bodyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation
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
			logrus.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var payload map[string]interface{}
		err = json.Unmarshal(respBody, &payload)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, payload)
	}
}

func (api *API) Test() gin.HandlerFunc {
	return func(c *gin.Context) {

		var reqBody struct {
			Count int `json:"count"`
		}

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			reqBody.Count = rand.Intn(100)
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		requestID := uuid.New().String()
		responseCh := pubsub2.GetResponseChannelMap().CreateChannel(requestID)
		defer pubsub2.GetResponseChannelMap().Delete(requestID)
		err = publishWorkRequest(api, requestID, workers.WORKER.Test, bodyBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		handleWorkResponse(c, responseCh)
		// worker handler implementation

	}
}
