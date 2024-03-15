package twitter

import (
	"context"
	"fmt"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// ScrapeTweetsByQuery searches tweets based on a query, with options for filtering and search mode.
// This function assumes that the caller has already logged in and will manage logout separately.
// It now returns a slice of Tweet pointers and an error.
func ScrapeTweetsByQuery(scraper *twitterscraper.Scraper, query string, count int, searchMode twitterscraper.SearchMode) ([]*twitterscraper.Tweet, error) {
	var tweets []*twitterscraper.Tweet

	if scraper == nil {
		logrus.Debug("Scraper instance is nil. Please initialize and log in before calling ScrapeTweetsByQuery.")
		return nil, fmt.Errorf("scraper instance is nil")
	}

	// Set search mode
	scraper.SetSearchMode(searchMode)

	// Perform the search with the specified query and count
	for tweetResult := range scraper.SearchTweets(context.Background(), query, count) {
		if tweetResult.Error != nil {
			logrus.Printf("Error fetching tweet: %v", tweetResult.Error)
			continue
		}
		tweets = append(tweets, &tweetResult.Tweet)
	}

	return tweets, nil
}
