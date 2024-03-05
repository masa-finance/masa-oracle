package twitter

import (
	"context"
	"fmt"
	"log"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

// ScrapeTweetsByQuery searches tweets based on a query, with options for filtering and search mode.
// This function assumes that the caller has already logged in and will manage logout separately.
func ScrapeTweetsByQuery(scraper *twitterscraper.Scraper, query string, count int, searchMode twitterscraper.SearchMode) {
	if scraper == nil {
		log.Fatal("Scraper instance is nil. Please initialize and log in before calling ScrapeTweetsByQuery.")
		return
	}

	// Set search mode
	scraper.SetSearchMode(searchMode)

	// Perform the search with the specified query and count
	for tweet := range scraper.SearchTweets(context.Background(), query, count) {
		if tweet.Error != nil {
			log.Printf("Error fetching tweet: %v", tweet.Error)
			continue
		}
		fmt.Println(tweet.Tweet.Text)
	}
}
