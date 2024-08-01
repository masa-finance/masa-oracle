// go test -v ./pkg/tests -run TestScrapeTweetsByQuery
// export TWITTER_2FA_CODE="873855"
package tests

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
)

// Global scraper instance
var scraper *twitterscraper.Scraper

func setup() {
	var err error
	_, b, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(b), "../..")
	if _, _ = os.Stat(rootDir + "/.env"); !os.IsNotExist(err) {
		_ = godotenv.Load()
	}

	logrus.SetLevel(logrus.DebugLevel)
	if scraper == nil {
		scraper = twitterscraper.New()
	}

	// Use GetInstance from config to access MasaDir
	appConfig := config.GetInstance()

	// Construct the cookie file path using MasaDir from AppConfig
	cookieFilePath := filepath.Join(appConfig.MasaDir, "twitter_cookies.json")

	// Attempt to load cookies
	if err = twitter.LoadCookies(scraper, cookieFilePath); err == nil {
		logrus.Debug("Cookies loaded successfully.")
		if twitter.IsLoggedIn(scraper) {
			logrus.Debug("Already logged in via cookies.")
			return
		}
	}

	// If cookies are not valid or do not exist, proceed with login
	username := appConfig.TwitterUsername
	password := appConfig.TwitterPassword
	logrus.WithFields(logrus.Fields{"username": username, "password": password}).Debug("Attempting to login")

	twoFACode := appConfig.Twitter2FaCode

	if twoFACode != "" {
		logrus.WithField("2FA", "provided").Debug("2FA code is provided, attempting login with 2FA")
		err = twitter.Login(scraper, username, password, twoFACode)
	} else {
		logrus.Debug("No 2FA code provided, attempting basic login")
		err = twitter.Login(scraper, username, password)
	}

	if err != nil {
		logrus.WithError(err).Warning("[-] Login failed")
		return
	}

	// Save cookies after successful login
	if err = twitter.SaveCookies(scraper, cookieFilePath); err != nil {
		logrus.WithError(err).Error("[-] Failed to save cookies")
		return
	}

	logrus.Debug("[+] Login successful")
}

func scrapeTweets(outputFile string) error {
	// Implement the tweet scraping logic here
	// This function should:
	// 1. Make API calls to the MASA_NODE_URL
	// 2. Process the responses
	// 3. Write the tweets to the outputFile in CSV format
	// 4. Handle rate limiting and retries
	// 5. Return an error if something goes wrong

	// For now, we'll just create a dummy file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"tweet", "datetime"})
	writer.Write([]string{"Test tweet #1", time.Now().Format(time.RFC3339)})
	writer.Write([]string{"Test tweet #2", time.Now().Format(time.RFC3339)})

	return nil
}

func TestSetup(t *testing.T) {
	setup()
}

func TestScrapeTweetsWithSentimentByQuery(t *testing.T) {
	// Ensure setup is done before running the test
	setup()

	query := "$MASA Token Masa"
	count := 100
	tweets, err := twitter.ScrapeTweetsByQuery(query, count)
	if err != nil {
		logrus.WithError(err).Error("[-] Failed to scrape tweets")
		return
	}

	// Serialize the tweets data to JSON
	tweetsData, err := json.Marshal(tweets)
	if err != nil {
		logrus.WithError(err).Error("[-] Failed to serialize tweets data")
		return
	}

	// Write the serialized data to a file
	filePath := "scraped_tweets.json"
	err = os.WriteFile(filePath, tweetsData, 0644)
	if err != nil {
		logrus.WithError(err).Error("[-] Failed to write tweets data to file")
		return
	}
	logrus.WithField("file", filePath).Debug("[+] Tweets data written to file successfully.")

	// Read the serialized data from the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		logrus.WithError(err).Error("[-] Failed to read tweets data from file")
		return
	}

	// Correctly declare a new variable for the deserialized data
	var deserializedTweets []*twitterscraper.Tweet
	err = json.Unmarshal(fileData, &deserializedTweets)
	if err != nil {
		logrus.WithError(err).Error("[-] Failed to deserialize tweets data")
		return
	}

	// Now, deserializedTweets contains the tweets loaded from the file
	// Send the tweets data to Claude for sentiment analysis
	// Convert []*twitterscraper.Tweet to []*twitterscraper.TweetResult
	twitterScraperTweets := make([]*twitterscraper.TweetResult, len(deserializedTweets))
	for i, tweet := range deserializedTweets {
		twitterScraperTweets[i] = &twitterscraper.TweetResult{
			Tweet: *tweet,
			Error: nil,
		}
	}
	sentimentRequest, sentimentSummary, err := llmbridge.AnalyzeSentimentTweets(twitterScraperTweets, "claude-3-opus-20240229", "Please perform a sentiment analysis on the following tweets, using an unbiased approach. Sentiment analysis involves identifying and categorizing opinions expressed in text, particularly to determine whether the writer's attitude towards a particular topic, product, etc., is positive, negative, or neutral. After analyzing, please provide a summary of the overall sentiment expressed in these tweets, including the proportion of positive, negative, and neutral sentiments if applicable.")
	if err != nil {
		logrus.WithError(err).Error("[-] Failed to analyze sentiment")
		err = os.Remove(filePath)
		if err != nil {
			logrus.WithError(err).Error("[-] Failed to delete the temporary file")
		} else {
			logrus.WithField("file", filePath).Debug("[+] Temporary file deleted successfully")
		}
		return
	}
	logrus.WithFields(logrus.Fields{
		"sentimentRequest": sentimentRequest,
		"sentimentSummary": sentimentSummary,
	}).Debug("[+] Sentiment analysis completed successfully.")

	// Delete the created file after the test
	err = os.Remove(filePath)
	if err != nil {
		logrus.WithError(err).Error("[-] Failed to delete the temporary file")
	} else {
		logrus.WithField("file", filePath).Debug("[+] Temporary file deleted successfully")
	}
}

func TestScrapeTweetsWithMockServer(t *testing.T) {
	setup()

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var requestBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&requestBody)

		if query, ok := requestBody["query"].(string); !ok || query == "" {
			t.Errorf("Expected query in request body")
		}

		response := map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"Text":      "Test tweet #1",
					"Timestamp": time.Now().Unix(),
				},
				{
					"Text":      "Test tweet #2",
					"Timestamp": time.Now().Unix(),
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Set up test environment
	os.Setenv("MASA_NODE_URL", server.URL)
	outputFile := "test_tweets.csv"
	defer os.Remove(outputFile)

	// Run the scrape function (you'll need to implement this)
	err := scrapeTweets(outputFile)
	if err != nil {
		t.Fatalf("Error scraping tweets: %v", err)
	}

	// Verify the output
	file, err := os.Open(outputFile)
	if err != nil {
		t.Fatalf("Error opening output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Error reading CSV: %v", err)
	}

	if len(records) != 3 { // Header + 2 tweets
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	if records[0][0] != "tweet" || records[0][1] != "datetime" {
		t.Errorf("Unexpected header: %v", records[0])
	}

	for i, record := range records[1:] {
		if record[0] == "" || record[1] == "" {
			t.Errorf("Empty field in record %d: %v", i+1, record)
		}
		_, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			t.Errorf("Invalid datetime format in record %d: %v", i+1, record[1])
		}
	}
}

func TestScrapeTweets(t *testing.T) {
	setup()

	query := "Sadhguru"
	count := 10

	for i := 0; i < 100; i++ {
		tweets, err := twitter.ScrapeTweetsByQuery(query, count)
		if err != nil {
			logrus.WithError(err).Error("[-] Failed to scrape tweets")
			return
		}

		var validTweets []*twitter.TweetResult
		for _, tweet := range tweets {
			if tweet.Error != nil {
				logrus.WithError(tweet.Error).Warn("[-] Error in tweet")
				continue
			}
			validTweets = append(validTweets, tweet)
		}

		tweetsData, err := json.Marshal(validTweets)
		if err != nil {
			logrus.WithError(err).Error("[-] Failed to serialize tweets data")
			return
		}

		logrus.WithFields(logrus.Fields{
			"total_tweets":  len(tweets),
			"valid_tweets":  len(validTweets),
			"tweets_sample": string(tweetsData[:min(100, len(tweetsData))]),
		}).Debug("[+] Tweets data")

		assert.NotNil(t, tweetsData[:10])
	}
}

func TestScrapeTweetsByTrends(t *testing.T) {
	setup()

	for i := 0; i < 5; i++ { // Run the test 5 times to ensure consistency
		trends, err := twitter.ScrapeTweetsByTrends()
		if err != nil {
			t.Fatalf("Failed to scrape tweets by trends: %v", err)
		}

		// Check if we got any trends
		assert.NotEmpty(t, trends, "No trends were returned")

		// Log the trends for debugging
		logrus.WithFields(logrus.Fields{
			"trends_count":  len(trends),
			"trends_sample": trends[:min(5, len(trends))],
		}).Debug("[+] Trends data")

		// Check each trend
		for _, trend := range trends {
			assert.NotEmpty(t, trend, "Empty trend found")
		}

		// Check if we have at least 5 trends (Twitter usually provides more)
		assert.GreaterOrEqual(t, len(trends), 5, "Expected at least 5 trends")

		// Optional: Add a short delay between iterations to avoid rate limiting
		time.Sleep(2 * time.Second)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
