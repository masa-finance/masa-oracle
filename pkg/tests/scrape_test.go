// go test -v ./pkg/tests -run TestScrapeTweetsByQuery
// export TWITTER_2FA_CODE="873855"
package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/twitter"
)

// Global scraper instance
var scraper *twitterscraper.Scraper

func setup() {
	logrus.SetLevel(logrus.DebugLevel)
	if scraper == nil {
		scraper = twitterscraper.New()
	}

	// Use GetInstance from config to access MasaDir
	appConfig := config.GetInstance()

	// Construct the cookie file path using MasaDir from AppConfig
	cookieFilePath := filepath.Join(appConfig.MasaDir, "twitter_cookies.json")

	// Attempt to load cookies
	if err := twitter.LoadCookies(scraper, cookieFilePath); err == nil {
		logrus.Debug("Cookies loaded successfully.")
		if twitter.IsLoggedIn(scraper) {
			logrus.Debug("Already logged in via cookies.")
			return
		}
	}

	// If cookies are not valid or do not exist, proceed with login
	username := appConfig.TwitterUsername
	password := appConfig.TwitterPassword
	logrus.WithFields(logrus.Fields{"username": username}).Debug("Attempting to login")

	twoFACode := appConfig.Twitter2FaCode
	var err error
	if twoFACode != "" {
		logrus.WithField("2FA", "provided").Debug("2FA code is provided, attempting login with 2FA")
		err = twitter.Login(scraper, username, password, twoFACode)
	} else {
		logrus.Debug("No 2FA code provided, attempting basic login")
		err = twitter.Login(scraper, username, password)
	}

	if err != nil {
		logrus.WithError(err).Fatal("Login failed")
		return
	}

	// Save cookies after successful login
	if err := twitter.SaveCookies(scraper, cookieFilePath); err != nil {
		logrus.WithError(err).Error("Failed to save cookies")
		return
	}

	logrus.Debug("Login successful")
}

func TestScrapeTweetsByQuery(t *testing.T) {
	// Ensure setup is done before running the test
	setup()

	query := "TRUMP"
	count := 200
	searchMode := twitterscraper.SearchLatest
	tweets, err := twitter.ScrapeTweetsByQuery(scraper, query, count, searchMode)
	if err != nil {
		logrus.WithError(err).Error("Failed to scrape tweets")
		return
	}

	// Serialize the tweets data to JSON
	tweetsData, err := json.Marshal(tweets)
	if err != nil {
		logrus.WithError(err).Error("Failed to serialize tweets data")
		return
	}

	// Write the serialized data to a file
	filePath := "scraped_tweets.json"
	err = os.WriteFile(filePath, tweetsData, 0644)
	if err != nil {
		logrus.WithError(err).Error("Failed to write tweets data to file")
		return
	}
	logrus.WithField("file", filePath).Debug("Tweets data written to file successfully.")
}
