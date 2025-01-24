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
	return tweet, loginEvent, nil
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
		tweets = append(tweets, &TweetResult{Tweet: &tweet.Tweet})
	}
	return tweets, loginEvent, nil
}
