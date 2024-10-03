package twitter

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"

	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

type TweetResult struct {
	Tweet *twitterscraper.Tweet
	Error error
}

// auth initializes and returns a new Twitter scraper instance. It attempts to load cookies from a file to reuse an existing session.
// If no valid session is found, it performs a login with credentials specified in the application's configuration.
// On successful login, it saves the session cookies for future use. If the login fails, it returns nil.
func Auth() *twitterscraper.Scraper {
	scraper := twitterscraper.New()
	appConfig := config.GetInstance()
	cookieFilePath := filepath.Join(appConfig.MasaDir, "twitter_cookies.json")

	if err := LoadCookies(scraper, cookieFilePath); err == nil {
		logrus.Debug("Cookies loaded successfully.")
		if IsLoggedIn(scraper) {
			logrus.Debug("Already logged in via cookies.")
			return scraper
		}
	}

	username := appConfig.TwitterUsername
	password := appConfig.TwitterPassword
	twoFACode := appConfig.Twitter2FaCode

	time.Sleep(100 * time.Millisecond)

	var err error
	if twoFACode != "" {
		err = Login(scraper, username, password, twoFACode)
	} else {
		err = Login(scraper, username, password)
	}

	if err != nil {
		logrus.WithError(err).Warning("[-] Login failed")
		return nil
	}

	time.Sleep(100 * time.Millisecond)

	if err = SaveCookies(scraper, cookieFilePath); err != nil {
		logrus.WithError(err).Error("[-] Failed to save cookies")
	}

	logrus.WithFields(logrus.Fields{
		"auth":     true,
		"username": username,
	}).Debug("Login successful")

	return scraper
}

// ScrapeTweetsByQuery performs a search on Twitter for tweets matching the specified query.
// It fetches up to the specified count of tweets and returns a slice of Tweet pointers.
// Parameters:
//   - query: The search query string to find matching tweets.
//   - count: The maximum number of tweets to retrieve.
//
// Returns:
//   - A slice of pointers to twitterscraper.Tweet objects that match the search query.
//   - An error if the scraping process encounters any issues.
func ScrapeTweetsByQuery(query string, count int) ([]*TweetResult, error) {
	scraper := Auth()
	var tweets []*TweetResult
	var lastError error

	if scraper == nil {
		return nil, fmt.Errorf("there was an error authenticating with your Twitter credentials")
	}

	// Set search mode
	scraper.SetSearchMode(twitterscraper.SearchLatest)

	// Perform the search with the specified query and count
	for tweetResult := range scraper.SearchTweets(context.Background(), query, count) {
		if tweetResult.Error != nil {
			lastError = tweetResult.Error
			logrus.Warnf("[+] Error encountered while scraping tweet: %v", tweetResult.Error)
			if strings.Contains(tweetResult.Error.Error(), "Rate limit exceeded") {
				return nil, fmt.Errorf("Twitter API rate limit exceeded (429 error)")
			}
			continue
		}
		tweets = append(tweets, &TweetResult{Tweet: &tweetResult.Tweet, Error: nil})
	}

	if len(tweets) == 0 && lastError != nil {
		return nil, lastError
	}

	return tweets, nil
}

// ScrapeTweetsProfile scrapes the profile and tweets of a specific Twitter user.
// It takes the username as a parameter and returns the scraped profile information and an error if any.
func ScrapeTweetsProfile(username string) (twitterscraper.Profile, error) {
	scraper := Auth()

	if scraper == nil {
		return twitterscraper.Profile{}, fmt.Errorf("there was an error authenticating with your Twitter credentials")
	}

	// Set search mode
	scraper.SetSearchMode(twitterscraper.SearchLatest)

	profile, err := scraper.GetProfile(username)
	if err != nil {
		return twitterscraper.Profile{}, err
	}

	return profile, nil
}
