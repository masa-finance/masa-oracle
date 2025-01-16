package api

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/masa-finance/masa-oracle/docs"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers"

	"github.com/gin-contrib/cors"

	"path/filepath"
	"runtime"

	"github.com/masa-finance/masa-oracle/node"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // ginSwagger middleware
)

//go:embed templates/*.html
var htmlTemplates embed.FS

// SetupRoutes configures the router with all API routes.
// It takes an OracleNode instance and returns a configured gin.Engine.
// Routes are added for peers, ads, subscriptions, node data, public keys,
// topics, the DHT, node status, and serving HTML pages. Middleware is added
// for CORS and templates.
func SetupRoutes(node *node.OracleNode, workerManager *workers.WorkHandlerManager, pubkeySubscriptionHandler *pubsub.PublicKeySubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	API := NewAPI(node, workerManager, pubkeySubscriptionHandler)

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

		if API.Node.Options.IsStaked {
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
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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
		// @Description Retrieves recent tweets based on query parameters, supporting advanced search options
		// @Tags Twitter
		// @Accept json
		// @Produce json
		// @Param body body object true "Search Query"
		// @Success 200 {array} Tweet "List of recent tweets"
		// @Failure 400 {object} ErrorResponse "Invalid query or error fetching tweets"
		// @Router /data/twitter/tweets/recent [post]
		// @Param body body object true "Search Query" SchemaExample({"query": "#MasaNode", "count": 10})
		// @Example hashtag {"query": "#MasaNode", "count": 10}
		// @Example mention {"query": "@getmasafi", "count": 10}
		// @Example fromUser {"query": "from:getmasafi", "count": 10}
		// @Example toUser {"query": "to:getmasafi", "count": 10}
		// @Example language {"query": "Masa lang:en", "count": 10}
		// @Example dateRange {"query": "Masa since:2021-01-01 until:2021-12-31", "count": 10}
		// @Example excludeRetweets {"query": "Masa -filter:retweets", "count": 10}
		// @Example minLikes {"query": "Masa min_faves:100", "count": 10}
		// @Example minRetweets {"query": "Masa min_retweets:50", "count": 10}
		// @Example keywordExclusion {"query": "Masa -moon", "count": 10}
		// @Example orOperator {"query": "Masa OR Oracle", "count": 10}
		// @Example geoLocation {"query": "Masa geocode:37.781157,-122.398720,1mi", "count": 10}
		// @Example urlInclusion {"query": "url:\"http://example.com\"", "count": 10}
		// @Example questionFilter {"query": "Masa ?", "count": 10}
		// @Example safeSearch {"query": "Masa filter:safe", "count": 10}
		v1.POST("/data/twitter/tweets/recent", API.SearchTweetsRecent())

		// @Summary Get Tweet by ID
		// @Description Retrieves a specific tweet by its ID
		// @Tags Twitter
		// @Accept  json
		// @Produce  json
		// @Param   tweet_id   path    string  true  "Tweet ID"
		// @Success 200 {object} Tweet "Successfully retrieved tweet"
		// @Failure 400 {object} ErrorResponse "Invalid tweet ID or error fetching tweet"
		// @Router /data/twitter/tweet/{tweet_id} [get]
		v1.GET("/data/twitter/tweet/:tweet_id", API.GetTweetByID())

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
	// Get the current file's directory
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)

	// Construct the path to the swagger.html file
	swaggerTemplate := filepath.Join(currentDir, "templates", "swagger.html")

	// Create a custom handler that serves our HTML file
	customHandler := func(c *gin.Context) {
		if c.Request.URL.Path == "/swagger" || c.Request.URL.Path == "/swagger/" || c.Request.URL.Path == "/swagger/index.html" {
			c.File(swaggerTemplate)
			return
		}

		// For other swagger-related paths, use the default handler
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
			return
		}

		// If it's not a swagger path, pass it to the next handler
		c.Next()
	}

	// Use our custom handler for all /swagger paths
	router.GET("/swagger", customHandler)
	router.GET("/swagger/*any", customHandler)
}
