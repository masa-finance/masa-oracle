package twitter

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/masa-finance/masa-oracle/pkg/config"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

var (
	accountManager *TwitterAccountManager
	once           sync.Once
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
	baseDir := config.GetInstance().MasaDir

	getAuthenticatedScraper := func() (*Scraper, *TwitterAccount, error) {
		account := accountManager.GetNextAccount()
		if account == nil {
			return nil, nil, fmt.Errorf("all accounts are rate-limited")
		}
		scraper := NewScraper(account, baseDir)
		if scraper == nil {
			logrus.Errorf("Authentication failed for %s", account.Username)
			return nil, account, fmt.Errorf("Twitter authentication failed for %s", account.Username)
		}
		return scraper, account, nil
	}

	scrapeTweets := func(scraper *Scraper) ([]*TweetResult, error) {
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

	return Retry(func() ([]*TweetResult, error) {
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
	}, MaxRetries)
}

func ScrapeTweetsProfile(username string) (twitterscraper.Profile, error) {
	once.Do(initializeAccountManager)
	baseDir := config.GetInstance().MasaDir

	getAuthenticatedScraper := func() (*Scraper, *TwitterAccount, error) {
		account := accountManager.GetNextAccount()
		if account == nil {
			return nil, nil, fmt.Errorf("all accounts are rate-limited")
		}
		scraper := NewScraper(account, baseDir)
		if scraper == nil {
			logrus.Errorf("Authentication failed for %s", account.Username)
			return nil, account, fmt.Errorf("Twitter authentication failed for %s", account.Username)
		}
		return scraper, account, nil
	}

	getProfile := func(scraper *Scraper) (twitterscraper.Profile, error) {
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

	return Retry(func() (twitterscraper.Profile, error) {
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
	}, MaxRetries)
}
