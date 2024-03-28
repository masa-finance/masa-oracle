package twitter

import (
	"context"
	_ "encoding/json"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// scrapeTweetsToChannel searches tweets based on a query, with options for filtering and search mode.
// This function assumes that the caller has already logged in and will manage logout separately.
// It now returns the tweets to a channel to be processed.
func scrapeTweetsToChannel(scraper *twitterscraper.Scraper, query string, count int, rowChan chan<- []*twitterscraper.Tweet) {
	var tweets []*twitterscraper.Tweet

	if scraper == nil {
		logrus.Debug("Scraper instance is nil. Please initialize and log in before calling scrapeTweetsToChannel.")
		return
	}

	// Set search mode
	scraper.SetSearchMode(twitterscraper.SearchLatest)

	// Perform the search with the specified query and count
	for tweetResult := range scraper.SearchTweets(context.Background(), query, count) {
		if tweetResult.Error != nil {
			logrus.Errorf("Error fetching tweet: %v", tweetResult.Error)
			continue
		}
		tweets = append(tweets, &tweetResult.Tweet)
	}

	defer close(rowChan)
	rowChan <- tweets

}
