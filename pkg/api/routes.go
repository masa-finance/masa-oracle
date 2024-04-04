package api

import (
	"embed"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Initialize CORS middleware with a configuration that allows all origins and specifies
	// the HTTP methods and headers that can be used in requests.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,                                      // Allow requests from any origin
		AllowMethods:    []string{"GET", "POST", "PUT", "OPTIONS"}, // Specify allowed methods
		AllowHeaders:    []string{"Origin", "Authorization"},       // Specify allowed headers
	}))

	// Define a list of routes that should not require authentication.
	ignoredRoutes := []string{
		"/status",
	}

	// Middleware to enforce API token authentication, excluding ignored routes.
	router.Use(func(c *gin.Context) {
		// Iterate over the ignored routes to determine if the current request should bypass authentication.
		for _, route := range ignoredRoutes {
			if c.Request.URL.Path == route {
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

	API := NewAPI(node)

	// Serving html
	templ := template.Must(template.ParseFS(htmlTemplates, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	router.GET("/peers", API.GetPeersHandler())

	router.GET("/peerAddresses", API.GetPeerAddresses())

	router.POST("/ads", API.PostAd())

	router.GET("/ads", API.GetAds())

	router.POST("/subscribeToAds", API.SubscribeToAds())

	router.GET("/nodeData", API.GetNodeDataHandler())

	router.GET("/nodeData/:peerID", API.GetNodeHandler())

	router.GET("/publicKeys", API.GetPublicKeysHandler())

	router.POST("/publishPublicKey", API.PublishPublicKeyHandler())

	router.POST("/createTopic", API.CreateNewTopicHandler())

	router.POST("/postToTopic", API.PostToTopicHandler())

	router.GET("/dht", API.GetFromDHT())

	router.POST("/dht", API.PostToDHT())

	router.POST("/nodestatus", API.PostNodeStatusHandler())

	// @todo
	// search/tweets - X
	// search/tweets/recent
	// search/tweets/popular
	// search/tweets/profile/:username

	router.POST("/search/tweets", API.SearchTweets())

	router.POST("/search/tweets/recent", nil)

	router.POST("/search/tweets/popular", nil)

	router.GET("/search/tweets/profile/{username}", nil)

	router.POST("/sentiment/tweets", API.SearchTweetsAndAnalyzeSentiment())

	router.POST("/sentiment/web", API.SearchWebAndAnalyzeSentiment())

	router.GET("/status", API.NodeStatusPageHandler())

	return router
}
