package twitter

import (
	"context"
	"fmt"
	"log"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

// ScrapeTweetsByQuery searches tweets based on a query, with options for filtering and search mode.
// It now accepts variadic `credentials` to support different login methods.
func ScrapeTweetsByQuery(query string, count int, searchMode twitterscraper.SearchMode, credentials ...string) {
	scraper := twitterscraper.New()

	// Login using the modular Login function from auth.go
	// Pass the variadic `credentials` directly to the Login function.
	err := Login(scraper, credentials...)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
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

	// Optionally, log out after the operation is complete
	err = Logout(scraper) // Use the modular Logout function from auth.go
	if err != nil {
		log.Printf("Error logging out: %v", err)
	}
}
