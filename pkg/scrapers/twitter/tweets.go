package twitter

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

var (
	accountManager *TwitterAccountManager
	once           sync.Once
	maxRetries     = 3
)

type TweetResult struct {
	Tweet *twitterscraper.Tweet
	Error error
}

func initializeAccountManager() {
	accounts := loadAccountsFromConfig()
	accountManager = NewTwitterAccountManager(accounts)
}

// loadAccountsFromConfig reads Twitter accounts from the .env file
func loadAccountsFromConfig() []*TwitterAccount {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("error loading .env file: %v", err)
	}

	accountsEnv := os.Getenv("TWITTER_ACCOUNTS")
	if accountsEnv == "" {
		logrus.Fatal("TWITTER_ACCOUNTS not set in .env file")
	}

	accountPairs := strings.Split(accountsEnv, ",")
	var accounts []*TwitterAccount

	for _, pair := range accountPairs {
		credentials := strings.Split(pair, ":")
		if len(credentials) != 2 {
			logrus.Warnf("invalid account credentials: %s", pair)
			continue
		}
		account := &TwitterAccount{
			Username: strings.TrimSpace(credentials[0]),
			Password: strings.TrimSpace(credentials[1]),
		}
		accounts = append(accounts, account)
	}

	return accounts
}

func ScrapeTweetsByQuery(query string, count int) ([]*TweetResult, error) {
	once.Do(initializeAccountManager)

	getAuthenticatedScraper := func() (*twitterscraper.Scraper, *TwitterAccount, error) {
		account := accountManager.GetNextAccount()
		if account == nil {
			return nil, nil, fmt.Errorf("all accounts are rate-limited")
		}
		scraper := NewScraper(account)
		if scraper == nil {
			logrus.Errorf("authentication failed for %s", account.Username)
			return nil, account, fmt.Errorf("authentication failed for %s", account.Username)
		}
		return scraper, account, nil
	}

	scrapeTweets := func(scraper *twitterscraper.Scraper) ([]*TweetResult, error) {
		var tweets []*TweetResult
		ctx := context.Background()
		scraper.SetSearchMode(twitterscraper.SearchLatest)
		for tweet := range scraper.SearchTweets(ctx, query, count) {
			if tweet.Error != nil {
				return nil, tweet.Error
			}
			tweets = append(tweets, &TweetResult{Tweet: &tweet.Tweet})
		}
		return tweets, nil
	}

	handleRateLimit := func(err error, account *TwitterAccount) bool {
		if strings.Contains(err.Error(), "Rate limit exceeded") {
			accountManager.MarkAccountRateLimited(account)
			logrus.Warnf("rate limited: %s", account.Username)
			return true
		}
		return false
	}

	return retryTweets(func() ([]*TweetResult, error) {
		scraper, account, err := getAuthenticatedScraper()
		if err != nil {
			return nil, err
		}

		tweets, err := scrapeTweets(scraper)
		if err != nil {
			if handleRateLimit(err, account) {
				return nil, err
			}
			return nil, err
		}
		return tweets, nil
	}, maxRetries)
}

func ScrapeTweetsProfile(username string) (twitterscraper.Profile, error) {
	once.Do(initializeAccountManager)

	getAuthenticatedScraper := func() (*twitterscraper.Scraper, *TwitterAccount, error) {
		account := accountManager.GetNextAccount()
		if account == nil {
			return nil, nil, fmt.Errorf("all accounts are rate-limited")
		}
		scraper := NewScraper(account)
		if scraper == nil {
			logrus.Errorf("authentication failed for %s", account.Username)
			return nil, account, fmt.Errorf("authentication failed for %s", account.Username)
		}
		return scraper, account, nil
	}

	getProfile := func(scraper *twitterscraper.Scraper) (twitterscraper.Profile, error) {
		return scraper.GetProfile(username)
	}

	handleRateLimit := func(err error, account *TwitterAccount) bool {
		if strings.Contains(err.Error(), "Rate limit exceeded") {
			accountManager.MarkAccountRateLimited(account)
			logrus.Warnf("rate limited: %s", account.Username)
			return true
		}
		return false
	}

	return retryProfile(func() (twitterscraper.Profile, error) {
		scraper, account, err := getAuthenticatedScraper()
		if err != nil {
			return twitterscraper.Profile{}, err
		}

		profile, err := getProfile(scraper)
		if err != nil {
			if handleRateLimit(err, account) {
				return twitterscraper.Profile{}, err
			}
			return twitterscraper.Profile{}, err
		}
		return profile, nil
	}, maxRetries)
}

func retryTweets(operation func() ([]*TweetResult, error), maxAttempts int) ([]*TweetResult, error) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}
		logrus.Errorf("retry attempt %d failed: %v", attempt, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return nil, fmt.Errorf("operation failed after %d attempts", maxAttempts)
}

func retryProfile(operation func() (twitterscraper.Profile, error), maxAttempts int) (twitterscraper.Profile, error) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}
		logrus.Errorf("retry attempt %d failed: %v", attempt, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return twitterscraper.Profile{}, fmt.Errorf("operation failed after %d attempts", maxAttempts)
}
