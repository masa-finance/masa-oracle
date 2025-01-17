package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/chain"
	"github.com/masa-finance/masa-oracle/pkg/config"
	pubsub2 "github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/telegram"
	"github.com/masa-finance/masa-oracle/pkg/workers"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type LLMChat struct {
	Model    string `json:"model,omitempty"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages,omitempty"`
	Stream bool `json:"stream"`
}

// sendWorkRequest sends a work request to a worker for processing.
// It marshals the request details into JSON and sends it over a libp2p stream.
// It is currently re-using the response channel map for this; however, it could be a simple synchronous call
// in which case the worker handlers would be responseible for preparing the data to be sent back to the client
//
// Parameters:
// - api: The API instance containing the Node and PubSubManager.
// - requestID: A unique identifier for the request.
// - workType: The type of work to be performed by the worker.
// - bodyBytes: The request body in byte slice format.
//
// Returns:
// - error: An error object if the request could not be sent or processed, otherwise nil.
func (api *API) sendWorkRequest(requestID string, workType data_types.WorkerType, bodyBytes []byte, wg *sync.WaitGroup) error {
	request := data_types.WorkRequest{
		WorkType:  workType,
		RequestId: requestID,
		Data:      bodyBytes,
	}
	response := api.WorkManager.DistributeWork(api.Node, request)
	responseChannel, exists := workers.GetResponseChannelMap().Get(requestID)
	if !exists {
		return fmt.Errorf("response channel not found")
	}
	select {
	case responseChannel <- response:
		wg.Add(1)
		// Successfully sent JSON response to the response channel
	default:
		// Log an error if the channel is blocking for debugging purposes
		logrus.Errorf("response channel is blocking for request ID: %s", requestID)
	}
	return nil
}

// handleWorkResponse processes the response from a worker and sends it back to the client.
// It listens on the provided response channel for a response or a timeout signal.
// If a response is received within the timeout period, it unmarshals the JSON response and sends it back to the client.
// If no response is received within the timeout period, it sends a timeout error to the client.
//
// Parameters:
// - c: The gin.Context object, which provides the context for the HTTP request.
// - responseCh: A channel that receives the worker's response as a byte slice.
func handleWorkResponse(c *gin.Context, responseCh <-chan data_types.WorkResponse, wg *sync.WaitGroup) {
	cfg, err := LoadConfig()
	if err != nil {
		handleError(c, "Failed to load API cfg", err)
		return
	}

	select {
	case response := <-responseCh:
		handleResponse(c, response, wg)
	case <-time.After(cfg.WorkerResponseTimeout):
		handleTimeout(c)
	case <-c.Done():
		// Context cancelled, no action needed
	}
}

func handleResponse(c *gin.Context, response data_types.WorkResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	if response.Error != "" {
		handleErrorResponse(c, response)
		return
	}

	if response.Data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":        "No data returned",
			"workerPeerId": response.WorkerPeerId,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func handleErrorResponse(c *gin.Context, response data_types.WorkResponse) {
	logrus.Errorf("[+] Work error: %s", response.Error)

	errorResponse := func(status int, message string) {
		c.JSON(status, gin.H{
			"error":        message,
			"details":      response.Error,
			"workerPeerId": response.WorkerPeerId,
		})
	}

	switch {
	case strings.Contains(response.Error, "Twitter API rate limit exceeded (429 error)"):
		errorResponse(http.StatusTooManyRequests, "Twitter API rate limit exceeded")
	case strings.Contains(response.Error, "no workers could process"):
		errorResponse(http.StatusServiceUnavailable, "No available workers to process the request")
	default:
		errorResponse(http.StatusInternalServerError, "An error occurred while processing the request")
	}
}

func handleError(c *gin.Context, message string, err error) {
	logrus.Errorf("%s: %v", message, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

func handleTimeout(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{"error": "Request timed out in API layer"})
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

		api.sendTrackingEvent(data_types.TwitterProfile, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.TwitterProfile, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
	}
}

// SearchTweetById returns a gin.HandlerFunc that processes a request to search for a tweet by its ID.
// It expects a JSON body with a field "id" (string), representing the tweet ID to search for.
// The handler validates the request body, ensuring the ID is not empty.
// If the request is valid, it attempts to scrape the tweet using the specified ID.
// On success, it returns the scraped tweet in a JSON response. On failure, it returns an appropriate error message and HTTP status code.
func (api *API) SearchTweetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			ID string `json:"id"`
		}

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if reqBody.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be provided and valid"})
			return
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		api.sendTrackingEvent(data_types.TwitterTweet, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.TwitterTweet, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
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

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		api.sendTrackingEvent(data_types.Twitter, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.Twitter, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
	}
}

// SearchTwitterFollowers returns a gin.HandlerFunc that retrieves the followers of a given Twitter user.
//
// Dev Notes:
// - This function uses URL parameters to get the username.
// - The default count is set to 20 if not provided.
func (api *API) SearchTwitterFollowers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			Username string `json:"username"`
			Count    int    `json:"count"`
		}

		username := c.Param("username")
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

		api.sendTrackingEvent(data_types.TwitterFollowers, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.TwitterFollowers, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
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
			UserID string `json:"userID"`
		}

		reqBody.UserID = c.Param("userID")

		if reqBody.UserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID must be provided and valid"})
			return
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		api.sendTrackingEvent(data_types.DiscordProfile, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.DiscordProfile, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
	}
}

// SearchChannelMessages returns a gin.HandlerFunc that processes a request to search for messages in a Discord channel.
func (api *API) SearchChannelMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqParams struct {
			ChannelID string `json:"channelID"`
			Limit     string `json:"limit"`
			Before    string `json:"before"`
		}

		reqParams.ChannelID = c.Param("channelID")
		if reqParams.ChannelID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ChannelID must be provided and valid"})
			return
		}

		reqParams.Limit = c.Query("limit")
		reqParams.Before = c.Query("before")

		if reqParams.Limit != "" {
			if _, err := strconv.Atoi(reqParams.Limit); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
				return
			}
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		api.sendTrackingEvent(data_types.DiscordChannelMessages, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.DiscordChannelMessages, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		wg.Wait()
	}
}

// SearchGuildChannels returns a gin.HandlerFunc that processes a request to search for channels in a Discord guild.
func (api *API) SearchGuildChannels() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			GuildID string `json:"guildID"`
		}

		reqBody.GuildID = c.Param("guildID")

		if reqBody.GuildID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserID must be provided and valid"})
			return
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		api.sendTrackingEvent(data_types.DiscordGuildChannels, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.DiscordGuildChannels, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
	}
}

// SearchUserGuilds returns a gin.HandlerFunc that processes a request to search for guilds associated with a Discord user.
func (api *API) SearchUserGuilds() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct{}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		api.sendTrackingEvent(data_types.DiscordUserGuilds, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.DiscordUserGuilds, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
	}
}

// SearchAllGuilds returns a gin.HandlerFunc that queries each node for the Discord guilds they are part of.
func (api *API) SearchAllGuilds() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get all nodes from the tracker
		peers := api.Node.NodeTracker.GetAllNodeData()

		// Prepare a wait group to wait for all go routines to finish
		var wg sync.WaitGroup
		wg.Add(len(peers))

		// This will store the combined list of guilds from all nodes
		allGuilds := make([]discord.Guild, 0)

		// Mutex to synchronize access to allGuilds slice
		var mutex sync.Mutex

		// Channel to collect errors
		errCh := make(chan error, len(peers))

		for _, p := range peers {
			go func(peer pubsub2.NodeData) {
				defer wg.Done()

				// Construct the URL for the GetUserGuilds endpoint
				var ipAddr string
				var err error
				for _, addr := range peer.Multiaddrs {
					ipAddr, err = addr.ValueForProtocol(multiaddr.P_IP4)
					if err == nil {
						break
					}
				}
				if ipAddr == "" {
					errCh <- fmt.Errorf("no IP4 address found for peer %s", peer.PeerId)
					return
				}

				url := fmt.Sprintf("http://%s:%s/api/v1/data/discord/user/guilds", ipAddr, os.Getenv("PORT"))

				// Make the HTTP request
				resp, err := http.Get(url)
				if err != nil {
					errCh <- fmt.Errorf("[-] Failed to make HTTP request: %v", err)
					return
				}

				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						logrus.Error("[-] Error closing response body: ", err)
					}
				}(resp.Body)
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					errCh <- fmt.Errorf("[-] Failed to read response body: %v", err)
					return
				}

				// Read and decode the response
				var result map[string]interface{}
				if err := json.Unmarshal(respBody, &result); err != nil {
					errCh <- fmt.Errorf("[-] Failed to unmarshal response body: %v", err)
					return
				}

				// Extract guilds from the result
				guildsData, ok := result["data"]
				if !ok {
					errCh <- fmt.Errorf("[-] Data field not found in response")
					return
				}

				guildsBytes, err := json.Marshal(guildsData)
				if err != nil {
					errCh <- fmt.Errorf("failed to marshal guilds data: %v", err)
					return
				}

				var guilds []discord.Guild
				if err := json.Unmarshal(guildsBytes, &guilds); err != nil {
					errCh <- fmt.Errorf("[-] Failed to unmarshal guilds: %v", err)
					return
				}

				// Safely append the guilds to the allGuilds slice
				mutex.Lock()
				allGuilds = append(allGuilds, guilds...)
				mutex.Unlock()
			}(p)
		}

		// Wait for all requests to finish
		wg.Wait()
		close(errCh)

		if len(allGuilds) > 0 {
			// Return the combined list of guilds
			c.JSON(http.StatusOK, gin.H{"guilds": allGuilds})
			return
		} else if len(errCh) > 0 {
			// Check if there were any errors
			for err := range errCh {
				logrus.Error("[-] Error fetching guilds: ", err)
			}
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "429 too many requests"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"guilds": allGuilds})
		}
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
		if !api.Node.Options.IsStaked {
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
			reqBody.Depth = 1 // Default count
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		api.sendTrackingEvent(data_types.Web, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.Web, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
	}
}

// StartAuth starts the authentication process with Telegram.
func (api *API) StartAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			PhoneNumber string `json:"phone_number"`
		}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		phoneCodeHash, err := telegram.StartAuthentication(context.Background(), reqBody.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start authentication"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Code sent to Telegram app", "phone_code_hash": phoneCodeHash})
	}
}

// CompleteAuth completes the authentication process with Telegram.
func (api *API) CompleteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			PhoneNumber   string `json:"phone_number"`
			Code          string `json:"code"`
			PhoneCodeHash string `json:"phone_code_hash"`
			Password      string `json:"password"`
		}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		auth, err := telegram.CompleteAuthentication(context.Background(), reqBody.PhoneNumber, reqBody.Code, reqBody.PhoneCodeHash, reqBody.Password)
		if err != nil {
			// Check if 2FA is required
			if err.Error() == "2FA required" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Two-factor authentication is required"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete authentication", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Authentication successful", "auth": auth})
	}
}

func (api *API) GetChannelMessagesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			Username string `json:"username"` // Telegram usernames are used instead of channel IDs
		}

		// Bind the JSON body to the struct
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if reqBody.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is missing"})
			return
		}

		// worker handler implementation
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		api.sendTrackingEvent(data_types.TelegramChannelMessages, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.TelegramChannelMessages, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
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

		api.sendTrackingEvent(data_types.LLMChat, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.LLMChat, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
	}
}

// TODO: review if we are still planning on doing the DfLlmChat and if so, make it conform to how we are doing other work

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
		api.sendTrackingEvent(data_types.LLMChat, bodyBytes)

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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logrus.Error("[-] Error closing response body: ", err)
			}
		}(resp.Body)
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Error("[-] Error reading response body: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var payload map[string]interface{}
		err = json.Unmarshal(respBody, &payload)
		if err != nil {
			logrus.Error("[-] Error unmarshalling response body: ", err)
			c.JSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, payload)
	}
}

// GetBlocks returns a gin.HandlerFunc that handles requests to retrieve all blocks from the blockchain.
//
// This function:
// 1. Checks if the node is a validator.
// 2. Retrieves all blocks from the blockchain.
// 3. Formats each block's data into a more readable structure.
// 4. Encodes the input data of each block in base64.
// 5. Returns the formatted blocks as a JSON response.
//
// The function is only accessible to validator nodes and will return an error for non-validator nodes.
func (api *API) GetBlocks() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !api.Node.Options.IsValidator {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not a validator and cannot access this endpoint"})
			return
		}

		type BlockData struct {
			Block            uint64      `json:"block"`
			InputData        interface{} `json:"input_data"`
			TransactionHash  string      `json:"transaction_hash"`
			PreviousHash     string      `json:"previous_hash"`
			TransactionNonce int         `json:"nonce"`
		}

		type Blocks struct {
			BlockData []BlockData `json:"blocks"`
		}
		var existingBlocks Blocks
		blocks := chain.GetBlockchain(api.Node.Blockchain)

		for _, block := range blocks {
			var inputData interface{}
			err := json.Unmarshal(block.Data, &inputData)
			if err != nil {
				inputData = string(block.Data) // Fallback to string if unmarshal fails
			}

			blockData := BlockData{
				Block:            block.Block,
				InputData:        base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", inputData))),
				TransactionHash:  fmt.Sprintf("%x", block.Hash),
				PreviousHash:     fmt.Sprintf("%x", block.Link),
				TransactionNonce: int(block.Nonce),
			}
			existingBlocks.BlockData = append(existingBlocks.BlockData, blockData)
		}

		jsonData, err := json.Marshal(existingBlocks)
		if err != nil {
			logrus.Error("[-] Error marshalling blocks: ", err)
			return
		}
		var blocksResponse Blocks
		err = json.Unmarshal(jsonData, &blocksResponse)
		if err != nil {
			logrus.Error("[-] Error unmarshalling blocks: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, blocksResponse)
	}
}

// GetBlockByHash returns a gin.HandlerFunc that handles requests to retrieve a specific block from the blockchain by its hash.
//
// This function:
// 1. Extracts the block hash from the request parameters.
// 2. Decodes the hexadecimal block hash.
// 3. Retrieves the block from the blockchain using the decoded hash.
// 4. Unmarshals the block data and formats it for the response.
// 5. Returns the formatted block data as a JSON response.
//
// If any errors occur during this process (e.g., invalid hash, block not found),
// appropriate error responses are sent back to the client.
func (api *API) GetBlockByHash() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !api.Node.Options.IsValidator {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not a validator and cannot access this endpoint"})
			return
		}

		blockHash := c.Param("blockHash")
		blockHashBytes, err := hex.DecodeString(blockHash)
		if err != nil {
			logrus.Errorf("[-]Failed to decode block hash: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid block hash"})
			return
		}
		block, err := api.Node.Blockchain.GetBlock(blockHashBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "hash not found"})
			return
		}

		var blockData map[string]interface{}
		err = json.Unmarshal(block.Data, &blockData)
		var inputData any
		if err != nil {
			inputData = string(block.Data)
		} else {
			inputData = blockData
		}
		responseData := gin.H{
			"block":            block.Block,
			"input_data":       inputData,
			"transaction_hash": blockHash,
			"nonce":            block.Nonce,
		}
		c.JSON(http.StatusOK, responseData)
	}
}

// Test is a temporary function that handles test requests.
// TODO: Remove this function once testing is complete.
func (api *API) Test() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody map[string]interface{}

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if len(reqBody) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body cannot be empty"})
			return
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = api.Node.PublishTopic(config.BlockTopic, bodyBytes)
		if err != nil {
			logrus.Errorf("[-] Error publishing block: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "message sent", "data": reqBody})
	}
}

func (api *API) sendTrackingEvent(workType data_types.WorkerType, jsonBytes []byte) {
	// Track work request event
	if api.EventTracker != nil && api.Node != nil {
		peerID := api.Node.Host.ID().String()
		api.EventTracker.TrackWorkRequest(workType, peerID, string(jsonBytes))
	} else {
		logrus.Warn("EventTracker or Node is nil in API")
	}
}
