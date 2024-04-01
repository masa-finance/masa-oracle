package twitter

import (
	"encoding/json"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// SerializeTweets serializes a slice of tweets to JSON.
func serializeTweets(tweets []*twitterscraper.Tweet) ([]*twitterscraper.Tweet, error) {
	// Serialize the tweets data to JSON
	tweetsData, err := json.Marshal(tweets)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize tweets data")
		return nil, err
	}
	// Correctly declare a new variable for the deserialized data
	var deserializedTweets []*twitterscraper.Tweet
	err = json.Unmarshal(tweetsData, &deserializedTweets)
	if err != nil {
		logrus.WithError(err).Error("Failed to deserialize tweets data")
		return nil, err
	}
	return deserializedTweets, nil
}
