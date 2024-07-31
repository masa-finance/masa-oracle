package api

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"github.com/masa-finance/masa-oracle/docs"

	"github.com/gin-contrib/cors"

	masa "github.com/masa-finance/masa-oracle/pkg"
)

//go:embed templates/*.html
var htmlTemplates embed.FS

// SetupRoutes configures the router with all API routes.
// It takes an OracleNode instance and returns a configured gin.Engine.
// Routes are added for peers, ads, subscriptions, node data, public keys,
// topics, the DHT, node status, and serving HTML pages. Middleware is added
// for CORS and templates.
func SetupRoutes(node *masa.OracleNode) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	API := NewAPI(node)

	// Initialize CORS middleware with a configuration that allows all origins and specifies
	// the HTTP methods and headers that can be used in requests.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:     true,                                      // Allow requests from any origin
		AllowMethods:        []string{"GET", "POST", "PUT", "OPTIONS"}, // Specify allowed methods
		AllowHeaders:        []string{"Origin", "Authorization"},       // Specify allowed headers
		AllowPrivateNetwork: true,
	}))

	// Define a list of routes that should not require authentication.
	ignoredRoutes := []string{
		"/status",
	}

	// Middleware to enforce API token authentication, excluding ignored routes.
	router.Use(func(c *gin.Context) {

		if API.Node.IsStaked {
			c.Next() // Proceed to the next middleware or handler as a staked node.
			return
		}

		// Iterate over the ignored routes to determine if the current request should bypass authentication.
		for _, route := range ignoredRoutes {
			if c.Request.URL.Path == route {
				c.Next() // Proceed to the next middleware or handler without authentication.
				return
			}
			if strings.HasPrefix(c.Request.URL.Path, "/auth") {
				c.Next() // Proceed to the next middleware or handler without authentication.
				return
			}
			if strings.HasPrefix(c.Request.URL.Path, "/health") {
				c.Next() // Proceed to the next middleware or handler without authentication.
				return
			}
			if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
				c.Next() // Proceed to the next middleware or handler without authentication.
				return
			}
		}

		// Define the prefix expected in the Authorization header.
		const BearerSchema = "Bearer "
		// Retrieve the Authorization header from the request.
		authHeader := c.GetHeader("Authorization")
		// If the Authorization header is missing, abort the request with an unauthorized status code and message.
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token required"})
			return
		}
		// Extract the token from the Authorization header by removing the Bearer schema prefix.
		token := authHeader[len(BearerSchema):]

		// Validate the token against the expected API key stored in environment variables.
		// logrus.Info(os.Getenv("API_KEY"))
		if os.Getenv("API_KEY") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT Token required"})
			return
		}

		// Validate the JWT token
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			//@todo decode the token get the apiKey, hash and compare it
			//if token != ???? {
			//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API token"})
			//	return
			//} else {
			//	c.Next()
			//	return
			//}
			return []byte(API.Node.Host.ID().String()), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			return
		}
		// Check if the token has expired
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT token has expired"})
				return
			}
		}
		c.Next()
	})

	// Serving html
	templ := template.Must(template.ParseFS(htmlTemplates, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	// Update Swagger info
	docs.SwaggerInfo.Host = "" // Leave this empty for relative URLs
	docs.SwaggerInfo.BasePath = "/api/v1"
	// Remove the Schemes setting from here, as we'll set it dynamically

	setupSwaggerHandler(router)

	v1 := router.Group("/api/v1")
	{

		// @Summary Get list of peers
		// @Description Retrieves a list of peers connected to the node
		// @Tags Peers
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} string "List of peer IDs"
		// @Router /peers [get]
		v1.GET("/peers", API.GetPeersHandler())

		// @Summary Get peer addresses
		// @Description Retrieves a list of peer addresses connected to the node
		// @Tags Peers
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} string "List of peer addresses"
		// @Router /peer/addresses [get]
		v1.GET("/peer/addresses", API.GetPeerAddresses())

		// @Summary Search Followers by Twitter Username
		// @Description Retrieves followers from a specific Twitter profile.
		// @Tags Twitter
		// @Accept  json
		// @Produce  json
		// @Param   username   path    string  true  "Twitter Username"
		// @Param   count   query   int     false  "Maximum number of users to return"  default(20)
		// @Success 200 {array} Profile "Array of profiles a user has as followers"
		// @Failure 400 {object} ErrorResponse "Invalid username or error fetching followers"
		// @Router /data/twitter/followers/{username} [get]
		v1.GET("/data/twitter/followers/:username", API.SearchTwitterFollowers())

		// @Summary Search Twitter Profile
		// @Description Retrieves tweets from a specific Twitter profile
		// @Tags Twitter
		// @Accept  json
		// @Produce  json
		// @Param   username   path    string  true  "Twitter Username"
		// @Success 200 {array} Tweet "List of tweets from the profile"
		// @Failure 400 {object} ErrorResponse "Invalid username or error fetching tweets"
		// @Router /data/twitter/profile/{username} [get]
		v1.GET("/data/twitter/profile/:username", API.SearchTweetsProfile())

		// @Summary Search recent tweets
		// @Description Retrieves recent tweets based on query parameters
		// @Tags Twitter
		// @Accept  json
		// @Produce  json
		// @Param   query	string  true  "Search Query"
		// @Success 200 {array} Tweet "List of recent tweets"
		// @Failure 400 {object} ErrorResponse "Invalid query or error fetching tweets"
		// @Router /data/twitter/tweets/recent [post]
		v1.POST("/data/twitter/tweets/recent", API.SearchTweetsRecent())

		// @Summary Twitter Trends
		// @Description Retrieves the latest Twitter trending topics
		// @Tags Twitter
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} Trend "List of trending topics"
		// @Failure 400 {object} ErrorResponse "Error fetching Twitter trends"
		// @Router /data/twitter/tweets/trends [get]
		v1.GET("/data/twitter/tweets/trends", API.SearchTweetsTrends())

		// @Summary Search Discord Profile
		// @Description Retrieves a Discord user profile by user ID.
		// @Tags Discord
		// @Accept  json
		// @Produce  json
		// @Param   userID   path    string  true  "Discord User ID"
		// @Success 200 {object} UserProfile "Successfully retrieved Discord user profile"
		// @Failure 400 {object} ErrorResponse "Invalid user ID or error fetching profile"
		// @Router /discord/profile/{userID} [get]
		v1.GET("/data/discord/profile/:userID", API.SearchDiscordProfile())

		// @Summary Start Telegram Authentication
		// @Description Initiates the authentication process with Telegram by sending a code to the provided phone number.
		// @Tags Authentication
		// @Accept  json
		// @Produce  json
		// @Param   phone_number   body    string  true  "Phone Number"
		// @Success 200 {object} map[string]interface{} "Successfully sent authentication code"
		// @Failure 400 {object} ErrorResponse "Invalid request body"
		// @Failure 500 {object} ErrorResponse "Failed to initialize Telegram client or to start authentication"
		// @Router /auth/telegram/start [post]
		v1.POST("/auth/telegram/start", API.StartAuth())

		// @Summary Complete Telegram Authentication
		// @Description Completes the authentication process with Telegram using the code sent to the phone number.
		// @Tags Authentication
		// @Accept  json
		// @Produce  json
		// @Param   phone_number   body    string  true  "Phone Number"
		// @Param   code           body    string  true  "Authentication Code"
		// @Param   phone_code_hash body   string  true  "Phone Code Hash"
		// @Success 200 {object} map[string]interface{} "Successfully authenticated"
		// @Failure 400 {object} ErrorResponse "Invalid request body"
		// @Failure 401 {object} ErrorResponse "Two-factor authentication is required"
		// @Failure 500 {object} ErrorResponse "Failed to initialize Telegram client or to complete authentication"
		// @Router /auth/telegram/complete [post]
		v1.POST("/auth/telegram/complete", API.CompleteAuth())

		// @Summary Get Telegram Channel Messages
		// @Description Retrieves messages from a specified Telegram channel.
		// @Tags Telegram
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} map[string][]Message "Successfully retrieved messages"
		// @Failure 400 {object} ErrorResponse "Username must be provided"
		// @Failure 500 {object} ErrorResponse "Failed to fetch channel messages"
		// @Router /telegram/channel/{username}/messages [get]
		v1.POST("/data/telegram/channel/messages", API.GetChannelMessagesHandler())

		// oauth tests
		// v1.GET("/data/discord/exchangetoken/:code", API.ExchangeDiscordTokenHandler())

		// @Summary Get messages from a Discord channel
		// @Description Retrieves messages from a specified Discord channel.
		// @Tags Discord
		// @Accept  json
		// @Produce  json
		// @Param   channelID   path    string  true  "Discord Channel ID"
		// @Success 200 {array} ChannelMessage "Successfully retrieved messages from the Discord channel"
		// @Failure 400 {object} ErrorResponse "Invalid channel ID or error fetching messages"
		// @Router /channels/{channelID}/messages [get]
		v1.GET("data/discord/channels/:channelID/messages", API.SearchChannelMessages())

		// @Summary Get channels from a Discord guild
		// @Description Retrieves channels from a specified Discord guild.
		// @Tags Discord
		// @Accept  json
		// @Produce  json
		// @Param   guildID   path    string  true  "Discord Guild ID"
		// @Success 200 {array} GuildChannel "Successfully retrieved channels from the Discord guild"
		// @Failure 400 {object} ErrorResponse "Invalid guild ID or error fetching channels"
		// @Router /guilds/{guildID}/channels [get]
		v1.GET("/data/discord/guilds/:guildID/channels", API.SearchGuildChannels())

		// @Summary Get all guilds
		// @Description Retrieves all guilds that the Discord Bots are apart of.
		// @Tags Discord
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} Guild "Successfully retrieved all guilds for the Discord user"
		// @Failure 400 {object} ErrorResponse "Error fetching guilds or invalid access token"
		// @Router /discord/guilds/all [get]
		v1.GET("/data/discord/guilds/all", API.SearchAllGuilds())

		// @Summary Get guilds for a Discord user
		// @Description Retrieves guilds that the authorized Discord user is part of.
		// @Tags Discord
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} Guild "Successfully retrieved guilds for the Discord user"
		// @Failure 400 {object} ErrorResponse "Error fetching guilds"
		// @Router /user/guilds [get]
		v1.GET("/data/discord/user/guilds", API.SearchUserGuilds())

		// @Summary Web Data
		// @Description Retrieves data from the web
		// @Tags Web
		// @Accept  json
		// @Produce  json
		// @Param   url   body    object  true  "Web Data Request"  example({"url": "https://hedgey.finance/"})
		// @Success 200 {object} WebDataResponse "Successfully retrieved web data"
		// @Failure 400 {object} ErrorResponse "Invalid URL or error fetching web data"
		// @Router /data/web [post]
		v1.POST("/data/web", API.WebData())

		// @Summary Get DHT Data
		// @Description Retrieves data from the DHT (Distributed Hash Table)
		// @Tags DHT
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} DHTResponse "Successfully retrieved data from DHT"
		// @Failure 400 {object} ErrorResponse "Error retrieving data from DHT"
		// @Router /dht [get]
		v1.GET("/dht", API.GetFromDHT())

		// @Summary Post to DHT
		// @Description Adds data to the DHT (Distributed Hash Table)
		// @Tags DHT
		// @Accept  json
		// @Produce  json
		// @Param   data   body    string  true  "Data to store in DHT"
		// @Success 200 {object} SuccessResponse "Successfully added data to DHT"
		// @Failure 400 {object} ErrorResponse "Error adding data to DHT"
		// @Router /dht [post]
		v1.POST("/dht", API.PostToDHT())

		// @Summary Get LLM Models
		// @Description Retrieves the available LLM models
		// @Tags LLM
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} LLMModelsResponse "Successfully retrieved LLM models"
		// @Failure 400 {object} ErrorResponse "Error retrieving LLM models"
		// @Router /llm/models [get]
		v1.GET("/llm/models", API.GetLLMModelsHandler())

		// @Summary Node Data
		// @Description Retrieves data from the node
		// @Tags Node
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} NodeDataResponse "Successfully retrieved node data"
		// @Failure 400 {object} ErrorResponse "Error retrieving node data"
		// @Router /node/data [get]
		v1.GET("/node/data", API.GetNodeDataHandler())

		// @Summary Get Node Data by Peer ID
		// @Description Retrieves data for a specific node identified by peer ID
		// @Tags Node
		// @Accept  json
		// @Produce  json
		// @Param   peerid   path    string  true  "Peer ID"
		// @Success 200 {object} NodeDataResponse "Successfully retrieved node data by peer ID"
		// @Failure 400 {object} ErrorResponse "Error retrieving node data by peer ID"
		// @Router /node/data/{peerid} [get]
		v1.GET("/node/data/:peerid", API.GetNodeHandler())

		// @Summary Update Node Status
		// @Description Updates the status of the node
		// @Tags Node
		// @Accept  json
		// @Produce  json
		// @Param   status   body    string  true  "Status to update"
		// @Success 200 {object} SuccessResponse "Successfully updated node status"
		// @Failure 400 {object} ErrorResponse "Error updating node status"
		// @Router /node/status [post]
		v1.POST("/node/status", API.PostNodeStatusHandler())

		// @Summary Get Public Keys
		// @Description Retrieves a list of public keys from the node
		// @Tags PublicKeys
		// @Accept  json
		// @Produce  json
		// @Success 200 {array} string "Successfully retrieved public keys"
		// @Failure 400 {object} ErrorResponse "Error retrieving public keys"
		// @Router /publickeys [get]
		v1.GET("/publickeys", API.GetPublicKeysHandler())

		// @Summary Publish Public Key
		// @Description Publishes a new public key to the node
		// @Tags PublicKeys
		// @Accept  json
		// @Produce  json
		// @Param   publickey   body    string  true  "Public Key to publish"
		// @Success 200 {object} SuccessResponse "Successfully published public key"
		// @Failure 400 {object} ErrorResponse "Error publishing public key"
		// @Router /publickey/publish [post]
		v1.POST("/publickey/publish", API.PublishPublicKeyHandler())

		// @Summary Analyze Sentiment of Tweets
		// @Description Searches for tweets and analyzes their sentiment
		// @Tags Sentiment
		// @Accept  json
		// @Produce  json
		// @Param   query   body    string  true  "Search Query"
		// @Success 200 {object} SentimentAnalysisResponse "Successfully analyzed sentiment of tweets"
		// @Failure 400 {object} ErrorResponse "Error analyzing sentiment of tweets"
		// @Router /sentiment/tweets [post]
		v1.POST("/sentiment/tweets", API.SearchTweetsAndAnalyzeSentiment())

		// @Summary Analyze Sentiment of Tweets
		// @Description Searches for tweets and analyzes their sentiment
		// @Tags Sentiment
		// @Accept  json
		// @Produce  json
		// @Param   query   body    string  true  "Prompt"
		// @Success 200 {object} SentimentAnalysisResponse "Successfully analyzed sentiment of discord"
		// @Failure 400 {object} ErrorResponse "Error analyzing sentiment of discord"
		// @Router /sentiment/tweets [post]
		v1.POST("/sentiment/discord", API.SearchDiscordMessagesAndAnalyzeSentiment())

		// @Summary Analyze Sentiment of Telegram Messages
		// @Description Searches for Telegram messages and analyzes their sentiment
		// @Tags Sentiment
		// @Accept  json
		// @Produce  json
		// @Param   query   body    string  true  "Search Query"
		// @Success 200 {object} SentimentAnalysisResponse "Successfully analyzed sentiment of Telegram messages"
		// @Failure 400 {object} ErrorResponse "Error analyzing sentiment of Telegram messages"
		// @Router /sentiment/telegram [post]
		v1.POST("/sentiment/telegram", API.SearchTelegramMessagesAndAnalyzeSentiment())

		// @Summary Analyze Sentiment of Web Content
		// @Description Searches for web content and analyzes its sentiment
		// @Tags Sentiment
		// @Accept  json
		// @Produce  json
		// @Param   query   body    string  true  "Search Query"
		// @Success 200 {object} SentimentAnalysisResponse "Successfully analyzed sentiment of web content"
		// @Failure 400 {object} ErrorResponse "Error analyzing sentiment of web content"
		// @Router /sentiment/web [post]
		v1.POST("/sentiment/web", API.SearchWebAndAnalyzeSentiment())

		// @Summary Create New Topic
		// @Description Creates a new discussion topic
		// @Tags Topics
		// @Accept  json
		// @Produce  json
		// @Param   topic   body    Topic  true  "Topic to create"
		// @Success 201 {object} TopicResponse "Successfully created new topic"
		// @Failure 400 {object} ErrorResponse "Error creating new topic"
		// @Router /topic/create [post]
		v1.POST("/topic/create", API.CreateNewTopicHandler())

		// @Summary Post to a Topic
		// @Description Adds a post to an existing discussion topic
		// @Tags Topics
		// @Accept  json
		// @Produce  json
		// @Param   post   body    Post  true  "Post content"
		// @Success 200 {object} PostResponse "Successfully added post to topic"
		// @Failure 400 {object} ErrorResponse "Error adding post to topic"
		// @Router /topic/post [post]
		v1.POST("/topic/post", API.PostToTopicHandler())

		// @Summary Chat with Local Ollama AI
		// @Description Initiates a chat session with an AI model that accepts common ollama formatted requests
		// @Tags Chat
		// @Accept  json
		// @Produce  json
		// @Param   reqBody  body      LLMChat  true  "Chat Request"
		// @Success 200 {object} ChatResponse "Successfully received response from AI"
		// @Failure 400 {object} ErrorResponse "Error communicating with AI"
		// @Router /chat [post]
		v1.POST("/chat", API.LocalLlmChat())

		// @Summary Chat using Cloudflare AI Workers
		// @Description Initiates a chat session with a Cloudflare AI model
		// @Tags Chat
		// @Accept  json
		// @Produce  json
		// @Param   message   body    string  true  "Message to send to Cloudflare AI"
		// @Success 200 {object} ChatResponse "Successfully received response from Cloudflare AI"
		// @Failure 400 {object} ErrorResponse "Error communicating with Cloudflare AI"
		// @Router /chat/cf [post]
		v1.POST("/chat/cf", API.CfLlmChat())

		// @Summary Get Blocks
		// @Description Retrieves the list of blocks from the blockchain
		// @Tags Blocks
		// @Accept  json
		// @Produce  json
		// @Success 200 {object} Blocks "Successfully retrieved blocks"
		// @Failure 400 {object} ErrorResponse "Error retrieving blocks"
		// @Router /blocks [get]
		v1.GET("/blocks", API.GetBlocks())

		// @Summary Get Block by Hash
		// @Description Retrieves a specific block from the blockchain using its hash
		// @Tags Blocks
		// @Accept  json
		// @Produce  json
		// @Param   blockHash   path    string  true  "Hash of the block to retrieve"
		// @Success 200 {object} BlockResponse "Successfully retrieved block"
		// @Failure 400 {object} ErrorResponse "Invalid block hash"
		// @Failure 500 {object} ErrorResponse "Error retrieving block"
		// @Router /blocks/{blockHash} [get]
		v1.GET("/blocks/:blockHash", API.GetBlockByHash())

		// @note a test route
		v1.POST("/test", API.Test())

	}

	// @Summary Node Status Page
	// @Description Retrieves the status page of the node
	// @Tags Status
	// @Accept  html
	// @Produce  html
	// @Success 200 {object} string "Successfully retrieved node status page"
	// @Failure 400 {object} ErrorResponse "Error retrieving node status page"
	// @Router /status [get]
	router.GET("/status", API.NodeStatusPageHandler())

	// @Summary Chat Page
	// @Description Renders the chat page for user interaction with the AI
	// @Tags Chat
	// @Accept  html
	// @Produce  html
	// @Param
	// @Success 200 {object} string "Successfully rendered chat page"
	// @Failure 500 {object} ErrorResponse "Error rendering chat page"
	// @Router /chat [get]
	router.GET("/chat", API.ChatPageHandler())

	// @Summary Get Node API Key
	// @Description Retrieves the API key for the node
	// @Tags Authentication
	// @Accept  json
	// @Produce  json
	// @Success 200 {object} map[string]interface{} "Successfully retrieved API key"
	// @Failure 500 {object} ErrorResponse "Error generating API key"
	// @Router /auth [get]
	router.GET("/auth", API.GetNodeApiKey())

	// @Summary Health Check
	// @Description Checks the health status of the API
	// @Tags Health
	// @Accept  json
	// @Produce  json
	// @Success 200 {object} map[string]bool "API health status"
	// @Router /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})

	return router
}

func setupSwaggerHandler(router *gin.Engine) {
	router.GET("/swagger", func(c *gin.Context) {
		c.Request.URL.Path = "/swagger/index.html"
		router.HandleContext(c)
	})
	router.GET("/swagger/", func(c *gin.Context) {
		c.Request.URL.Path = "/swagger/index.html"
		router.HandleContext(c)
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/swagger/doc.json", func(c *gin.Context) {
		doc, err := swag.ReadDoc()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read Swagger doc"})
			return
		}

		var swaggerSpec map[string]interface{}
		if err := json.Unmarshal([]byte(doc), &swaggerSpec); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse Swagger doc"})
			return
		}

		// Determine the scheme
		scheme := "http"
		if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}

		// Update the schemes in the Swagger spec
		swaggerSpec["schemes"] = []string{scheme}

		c.JSON(http.StatusOK, swaggerSpec)
	})
}
