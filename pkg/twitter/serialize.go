package twitter

import (
	"encoding/json"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// SerializeTweets serializes a slice of tweets to JSON.
func SerializeTweets(tweets []*twitterscraper.Tweet) ([]byte, error) {
	tweetsData, err := json.Marshal(tweets)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize tweets data")
		return nil, err
	}
	return tweetsData, nil
}
