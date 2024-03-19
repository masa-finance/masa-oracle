package masa

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewBridge() error {

	//logrus.Info("starting server")
	//router := gin.Default()
	////router.Use(cors.Default())
	//// router.SetTrustedProxies([]string{add values here})
	//
	//// Use the auth middleware for the /webhook route
	//router.POST("/webhook", authMiddleware(), webhookHandler)
	//
	//// Paths to the certificate and key files
	//certFile := os.Getenv(Cert)
	//keyFile := os.Getenv(CertPem)
	//
	//if err := router.RunTLS(":8080", certFile, keyFile); err != nil {
	//	return err
	//}
	return nil
}

// authMiddleware returns a middleware function that checks for a valid
// authorization token in the request header. If the token is not valid,
// it aborts the request with a 401 Unauthorized status code. Otherwise
// it calls the next handler in the chain. This can be used to protect
// routes that require authentication.
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		token := c.GetHeader("Authorization")

		// Check the token
		if token != "your_expected_token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

// webhookHandler handles incoming webhook requests.
// It simply returns a 200 OK response with a message.
func webhookHandler(c *gin.Context) {
	// Handle the webhook request here

	c.JSON(http.StatusOK, gin.H{"message": "Webhook called"})
}
