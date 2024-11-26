package twitter

import (
	"context"

	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type TweetResult struct {
	Tweet *twitterscraper.Tweet
	Error error
}

func ScrapeTweetsByQuery(baseDir string, query string, count int) ([]*TweetResult, error) {
	scraper, account, err := getAuthenticatedScraper(baseDir)
	if err != nil {
		return nil, err
	}

	var tweets []*TweetResult
	ctx := context.Background()
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	for tweet := range scraper.SearchTweets(ctx, query, count) {
		if tweet.Error != nil {
			if handleRateLimit(tweet.Error, account) {
				return nil, tweet.Error
			}
			return nil, tweet.Error
		}
		tweets = append(tweets, &TweetResult{Tweet: &tweet.Tweet})
	}
	return tweets, nil
}
