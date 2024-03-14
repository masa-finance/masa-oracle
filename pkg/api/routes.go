package api

import (
	"embed"
	"github.com/gin-gonic/gin"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"html/template"
)

// Before:
// //go:embed pkg/api/templates/*.html
// After, assuming the Go file is directly inside pkg/api:
//
//go:embed templates/*.html
var htmlTemplates embed.FS

func SetupRoutes(node *masa.OracleNode) *gin.Engine {
	router := gin.Default()
	// @TODO need to add a Authorization Bearer methodology for api security
	// add cors middleware

	API := NewAPI(node)

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

	// Serving node status html
	templ := template.Must(template.ParseFS(htmlTemplates, "templates/*.html"))
	router.SetHTMLTemplate(templ)

	router.GET("/status", API.NodeStatusPageHandler())

	return router
}
