package api

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"strings"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/docs"
	masa "github.com/masa-finance/masa-oracle/pkg"
)

// Before:
// //go:embed pkg/api/templates/*.html
// After, assuming the Go file is directly inside pkg/api:
//
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
		AllowAllOrigins: true,                                      // Allow requests from any origin
		AllowMethods:    []string{"GET", "POST", "PUT", "OPTIONS"}, // Specify allowed methods
		AllowHeaders:    []string{"Origin", "Authorization"},       // Specify allowed headers
	}))

	// Define a list of routes that should not require authentication.
	ignoredRoutes := []string{
		"/api/v1/status",
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
		if token != os.Getenv("API_KEY") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API token"})
			return
		}
		c.Next() // Proceed to the next middleware or handler since authentication is successful.
	})

	// Serving html
	templ := template.Must(template.ParseFS(htmlTemplates, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	//	@contact.name	Masa API Support
	//	@contact.url	https://api.masa.ai
	//	@contact.email	support@masa.ai

	//	@license.name	MIT
	//	@license.url	https://opensource.org/license/mit

	//	@BasePath		/api/v1

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "Masa Oracle"
	docs.SwaggerInfo.Host = "https://api.masa.ai"
	docs.SwaggerInfo.Description = "The World's Personal Data Network Masa Oracle Node API."
	docs.SwaggerInfo.Version = "0.0.10-alpha"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	v1 := router.Group("/api/v1")
	{
		//	@Summary		Get peers
		//	@Description	Retrieves a list of peers
		//	@Tags			peers
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{array}	api.Peer
		//	@Router			/peers [get]
		v1.GET("/peers", API.GetPeersHandler())

		//	@Summary		Get peer addresses
		//	@Description	Retrieves the addresses of connected peers
		//	@Tags			peers
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{array}	string
		//	@Router			/peer/addresses [get]
		v1.GET("/peer/addresses", API.GetPeerAddresses())

		//	@Summary		Post an ad
		//	@Description	Creates a new ad
		//	@Tags			ads
		//	@Accept			json
		//	@Produce		json
		//	@Param			ad	body		api.Ad	true	"Ad object"
		//	@Success		200	{object}	api.Ad
		//	@Router			/ads [post]
		v1.POST("/ads", API.PostAd())

		//	@Summary		Get ads
		//	@Description	Retrieves a list of ads
		//	@Tags			ads
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{array}	api.Ad
		//	@Router			/ads [get]
		v1.GET("/ads", API.GetAds())

		//	@Summary		Subscribe to ads
		//	@Description	Subscribes the node to receive ads
		//	@Tags			ads
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{object}	api.SubscriptionResponse
		//	@Router			/ads/subscribe [post]
		v1.POST("/ads/subscribe", API.SubscribeToAds())

		//	@Summary		Search tweets by profile
		//	@Description	Retrieves tweets from a specific user profile
		//	@Tags			twitter
		//	@Accept			json
		//	@Produce		json
		//	@Param			username	path		string	true	"Twitter username"
		//	@Success		200			{object}	gin.H
		//	@Failure		400			{object}	gin.H
		//	@Failure		500			{object}	gin.H
		//	@Router			/data/twitter/profile/{username} [get]
		v1.GET("/data/twitter/profile/:username", API.SearchTweetsProfile())

		//	@Summary		Search recent tweets
		//	@Description	Searches for recent tweets based on a query and count
		//	@Tags			twitter
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.SearchTweetsRecentRequest	true	"Search request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/data/twitter/tweets/recent [post]
		v1.POST("/data/twitter/tweets/recent", API.SearchTweetsRecent())

		//	@Summary		Search trending tweets
		//	@Description	Retrieves trending tweets from Twitter
		//	@Tags			twitter
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{array}		string
		//	@Failure		500	{object}	gin.H
		//	@Router			/data/twitter/tweets/trends [get]
		v1.GET("/data/twitter/tweets/trends", API.SearchTweetsTrends())

		//	@Summary		Scrape web data
		//	@Description	Scrapes web data from a given URL up to a specified depth
		//	@Tags			web
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.WebDataRequest	true	"Web data request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/data/web [post]
		v1.POST("/data/web", API.WebData())

		//	@Summary		Get data from DHT
		//	@Description	Retrieves data from the Distributed Hash Table (DHT)
		//	@Tags			DHT
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{object}	gin.H
		//	@Failure		400	{object}	gin.H
		//	@Failure		500	{object}	gin.H
		//	@Router			/dht [get]
		v1.GET("/dht", API.GetFromDHT())

		//	@Summary		Post data to DHT
		//	@Description	Posts data to the Distributed Hash Table (DHT)
		//	@Tags			DHT
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.PostToDHTRequest	true	"Post request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/dht [post]
		v1.POST("/dht", API.PostToDHT())

		//	@Summary		Get node data
		//	@Description	Retrieves data from the node
		//	@Tags			node
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{object}	gin.H
		//	@Failure		400	{object}	gin.H
		//	@Failure		500	{object}	gin.H
		//	@Router			/node/data [get]
		v1.GET("/node/data", API.GetNodeDataHandler())

		//	@Summary		Get node data by peer ID
		//	@Description	Retrieves data from the node by peer ID
		//	@Tags			node
		//	@Accept			json
		//	@Produce		json
		//	@Param			peerid	path		string	true	"Peer ID"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/node/data/{peerid} [get]
		v1.GET("/node/data/:peerid", API.GetNodeHandler())

		//	@Summary		Update node status
		//	@Description	Updates the status of the node
		//	@Tags			node
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.PostNodeStatusRequest	true	"Post request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/node/status [post]
		v1.POST("/node/status", API.PostNodeStatusHandler())

		//	@Summary		Get public keys
		//	@Description	Retrieves the public keys
		//	@Tags			keys
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{object}	gin.H
		//	@Failure		400	{object}	gin.H
		//	@Failure		500	{object}	gin.H
		//	@Router			/publickeys [get]
		v1.GET("/publickeys", API.GetPublicKeysHandler())

		//	@Summary		Publish public key
		//	@Description	Publishes the public key to the network
		//	@Tags			keys
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.PublishPublicKeyRequest	true	"Publish request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/publickey/publish [post]
		v1.POST("/publickey/publish", API.PublishPublicKeyHandler())

		//	@Summary		Analyze sentiment of tweets
		//	@Description	Searches for tweets and analyzes their sentiment
		//	@Tags			sentiment
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.SearchTweetsAndAnalyzeSentimentRequest	true	"Search and Analyze request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/sentiment/tweets [post]
		v1.POST("/sentiment/tweets", API.SearchTweetsAndAnalyzeSentiment())

		//	@Summary		Analyze sentiment of web data
		//	@Description	Searches for web data and analyzes their sentiment
		//	@Tags			sentiment
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.SearchWebAndAnalyzeSentimentRequest	true	"Search and Analyze request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/sentiment/web [post]
		v1.POST("/sentiment/web", API.SearchWebAndAnalyzeSentiment())

		//	@Summary		Get node status
		//	@Description	Retrieves the status of the node
		//	@Tags			status
		//	@Accept			json
		//	@Produce		json
		//	@Success		200	{object}	api.NodeStatus
		//	@Router			/status [get]
		v1.GET("/status", API.NodeStatusPageHandler())

		//	@Summary		Create a new topic
		//	@Description	Creates a new topic in the network
		//	@Tags			topics
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.CreateNewTopicRequest	true	"Create new topic request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/topic/create [post]
		v1.POST("/topic/create", API.CreateNewTopicHandler())

		//	@Summary		Post to a topic
		//	@Description	Posts a message to a specific topic
		//	@Tags			topics
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		api.PostToTopicRequest	true	"Post to topic request"
		//	@Success		200		{object}	gin.H
		//	@Failure		400		{object}	gin.H
		//	@Failure		500		{object}	gin.H
		//	@Router			/topic/post [post]
		v1.POST("/topic/post", API.PostToTopicHandler())
	}

	//	@Summary		Swagger UI
	//	@Description	Serves the Swagger UI for API documentation
	//	@Tags			swagger
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{object}	swagger.WrapHandler
	//	@Router			/swagger/{any} [get]
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
