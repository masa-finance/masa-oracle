package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/chain"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/workers"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

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

	err := response.UnsealDataIfNeeded()
	if err != nil {
		return fmt.Errorf("failed to get response data: %v", err)
	}

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

	if response.Data == "" {
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

// GetTweetByID returns a gin.HandlerFunc that processes a request to get a specific tweet by ID.
// The tweet ID is expected as a URL parameter.
// On success, it returns the tweet data in a JSON response.
// On failure, it returns an appropriate error message and HTTP status code.
func (api *API) GetTweetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		tweetID := c.Param("tweet_id")
		if tweetID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tweet ID must be provided"})
			return
		}

		// Create request body
		bodyBytes, err := json.Marshal(map[string]interface{}{
			"tweet_id": tweetID,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		api.sendTrackingEvent(data_types.TwitterTweetByID, bodyBytes)
		requestID := uuid.New().String()
		responseCh := workers.GetResponseChannelMap().CreateChannel(requestID)
		wg := &sync.WaitGroup{}
		defer workers.GetResponseChannelMap().Delete(requestID)
		go handleWorkResponse(c, responseCh, wg)

		err = api.sendWorkRequest(requestID, data_types.TwitterTweetByID, bodyBytes, wg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		wg.Wait()
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
