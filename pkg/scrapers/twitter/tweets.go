package twitter

import (
	"context"
	"fmt"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

// auth initializes and returns a new Twitter scraper instance. It attempts to load cookies from a file to reuse an existing session.
// If no valid session is found, it performs a login with credentials specified in the application's configuration.
// On successful login, it saves the session cookies for future use. If the login fails, it returns nil.
func auth() *twitterscraper.Scraper {
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

	var err error
	if twoFACode != "" {
		err = Login(scraper, username, password, twoFACode)
	} else {
		err = Login(scraper, username, password)
	}

	if err != nil {
		logrus.WithError(err).Warning("Login failed")
		return nil
	}

	if err = SaveCookies(scraper, cookieFilePath); err != nil {
		logrus.WithError(err).Error("Failed to save cookies")
	}

	logrus.WithFields(logrus.Fields{
		"auth":     true,
		"username": username,
	}).Debug("Login successful")

	return scraper
}

// ScrapeTweetsForSentiment is a function that scrapes tweets based on a given query, analyzes their sentiment using a specified model, and returns the sentiment analysis results.
// Parameters:
//   - query: The search query string to find matching tweets.
//   - count: The maximum number of tweets to retrieve and analyze.
//   - model: The model to use for sentiment analysis.
//
// Returns:
//   - A string representing the sentiment analysis prompt.
//   - A string representing the sentiment analysis result.
//   - An error if the scraping or sentiment analysis process encounters any issues.
func ScrapeTweetsForSentiment(query string, count int, model string) (string, string, error) {
	scraper := auth()
	var tweets []*twitterscraper.Tweet

	if scraper == nil {
		return "", "", fmt.Errorf("there was an error authenticating with your Twitter credentials")
	}

	// Set search mode
	scraper.SetSearchMode(twitterscraper.SearchLatest)

	// Perform the search with the specified query and count
	for tweetResult := range scraper.SearchTweets(context.Background(), query, count) {
		if tweetResult.Error != nil {
			logrus.Printf("Error fetching tweet: %v", tweetResult.Error)
			continue
		}
		tweets = append(tweets, &tweetResult.Tweet)
	}
	sentimentPrompt := "Please perform a sentiment analysis on the following tweets, using an unbiased approach. Sentiment analysis involves identifying and categorizing opinions expressed in text, particularly to determine whether the writer's attitude towards a particular topic, product, etc., is positive, negative, or neutral. After analyzing, please provide a summary of the overall sentiment expressed in these tweets, including the proportion of positive, negative, and neutral sentiments if applicable."
	prompt, sentiment, err := llmbridge.AnalyzeSentimentTweets(tweets, model, sentimentPrompt)
	if err != nil {
		return "", "", err
	}
	return prompt, sentiment, nil
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
func ScrapeTweetsByQuery(query string, count int) ([]*twitterscraper.Tweet, error) {
	scraper := auth()
	var tweets []*twitterscraper.Tweet

	if scraper == nil {
		return nil, fmt.Errorf("there was an error authenticating with your Twitter credentials")
	}

	// Set search mode
	scraper.SetSearchMode(twitterscraper.SearchLatest)

	// Perform the search with the specified query and count
	for tweetResult := range scraper.SearchTweets(context.Background(), query, count) {
		if tweetResult.Error != nil {
			logrus.Printf("Error fetching tweet: %v", tweetResult.Error)
			continue
		}
		tweets = append(tweets, &tweetResult.Tweet)
	}
	return tweets, nil
}

// ScrapeTweetsByTrends scrapes the current trending topics on Twitter.
// It returns a slice of strings representing the trending topics.
// If an error occurs during the scraping process, it returns an error.
func ScrapeTweetsByTrends() ([]string, error) {
	scraper := auth()
	var tweets []string

	if scraper == nil {
		return nil, fmt.Errorf("there was an error authenticating with your Twitter credentials")
	}

	// Set search mode
	scraper.SetSearchMode(twitterscraper.SearchLatest)

	trends, err := scraper.GetTrends()
	if err != nil {
		logrus.Printf("Error fetching tweet: %v", err)
		return nil, err
	}

	tweets = append(tweets, trends...)

	return tweets, nil
}

// ScrapeTweetsProfile scrapes the profile and tweets of a specific Twitter user.
// It takes the username as a parameter and returns the scraped profile information and an error if any.
func ScrapeTweetsProfile(username string) (twitterscraper.Profile, error) {
	scraper := auth()

	if scraper == nil {
		return twitterscraper.Profile{}, fmt.Errorf("there was an error authenticating with your Twitter credentials")
	}

	// Set search mode
	scraper.SetSearchMode(twitterscraper.SearchLatest)

	profile, err := scraper.GetProfile(username)
	if err != nil {
		logrus.Printf("Error fetching profile: %v", err)
		return twitterscraper.Profile{}, err
	}

	return profile, nil
}
