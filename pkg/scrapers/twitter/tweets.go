package twitter

import (
	"context"

	twitterscraper "github.com/imperatrona/twitter-scraper"

	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type TweetResult struct {
	Tweet *twitterscraper.Tweet
	Error error
}

// Function to remove cyclic references
func removeCyclicReferences(tweet *twitterscraper.Tweet) *twitterscraper.Tweet {
	if tweet == nil {
		return nil
	}

	// Create a copy of the tweet to avoid modifying the original
	cleanTweet := *tweet

	// Set cyclic reference fields to nil
	cleanTweet.InReplyToStatus = nil
	cleanTweet.QuotedStatus = nil
	cleanTweet.RetweetedStatus = nil
	cleanTweet.Thread = nil

	return &cleanTweet
}

func ScrapeTweetByID(id string) (*twitterscraper.Tweet, *data_types.LoginEvent, error) {
	scraper, account, loginEvent, err := getAuthenticatedScraper()
	if err != nil {
		return nil, loginEvent, err
	}

	tweet, err := scraper.GetTweet(id)
	if err != nil {
		if handleRateLimit(err, account) {
			return nil, loginEvent, err
		}
		return nil, loginEvent, err
	}

	// Remove cyclic references before conversion
	cleanTweet := removeCyclicReferences(tweet)
	return cleanTweet, loginEvent, nil
}

func ScrapeTweetsByQuery(query string, count int) ([]*TweetResult, *data_types.LoginEvent, error) {
	scraper, account, loginEvent, err := getAuthenticatedScraper()
	if err != nil {
		return nil, loginEvent, err
	}

	var tweets []*TweetResult
	ctx := context.Background()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	for tweet := range scraper.SearchTweets(ctx, query, count) {
		if tweet.Error != nil {
			if handleRateLimit(tweet.Error, account) {
				return nil, loginEvent, tweet.Error
			}
			return nil, loginEvent, tweet.Error
		}

		// Remove cyclic references before conversion
		cleanTweet := removeCyclicReferences(&tweet.Tweet)
		tweets = append(tweets, &TweetResult{Tweet: cleanTweet})
	}
	return tweets, loginEvent, nil
}
