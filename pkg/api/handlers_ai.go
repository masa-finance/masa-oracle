package api

import (
	"github.com/gin-gonic/gin"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/twitter"
	"net/http"
	"strconv"
)

func (api *API) PostSentiment() gin.HandlerFunc {
	return func(c *gin.Context) {

		sharedData := db.SharedData{}
		if err := c.BindJSON(&sharedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "invalid request",
			})
			return
		}

		query := sharedData["query"].(string)
		count, _ := strconv.Atoi(sharedData["count"].(string))

		// count, _ := strconv.ParseInt(sharedData["count"], 10, 64)

		//jsonData, err := json.Marshal(sharedData["value"])
		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"success": false,
		//		"message": "invalid json",
		//	})
		//	return
		//}
		// WIP testing Twitter scrape to sentiment
		// created a tweet reader with worker threads
		// simplified the flow using channels
		// store AI request and sentiment to datastore
		// cnt, err := strconv.Atoi(count)
		//if err !=nil{
		//
		//}
		twitter.Scrape(query, count)

		//if err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{
		//		"success": success,
		//		"message": keyStr,
		//	})
		//	return
		//}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": sharedData,
		})
	}
}
