package twitter

import (
	"context"
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

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
		logrus.WithError(err).Fatal("Login failed")
		return nil
	}

	if err = SaveCookies(scraper, cookieFilePath); err != nil {
		logrus.WithError(err).Error("Failed to save cookies")
	}

	logrus.WithFields(logrus.Fields{
		"auth":     true,
		"username": os.Getenv("TWITTER_USER"),
	}).Debug("Login successful")

	return scraper
}

func Scrape(query string, count int) {

	rowChan := make(chan []*twitterscraper.Tweet)

	scraper := auth()
	go scrapeTweetsToChannel(scraper, query, count, rowChan)

	var wg sync.WaitGroup
	var size int64
	var numWorkers = 5

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processTweets(&wg, rowChan, &size)
	}
	wg.Wait()
	logrus.Println(size)
}

func processTweets(wg *sync.WaitGroup, rowChan chan []*twitterscraper.Tweet, size *int64) {
	defer wg.Done()
	for row := range rowChan {
		if row != nil {
			atomic.AddInt64(size, 1)
			logrus.Println("row", row)

			deserializedTweets, err := serializeTweets(row)
			if err != nil {
				logrus.WithError(err).Error("Failed to deserialize tweets data")
				return
			}

			sentimentRequest, sentimentSummary, e := llmbridge.AnalyzeSentiment(deserializedTweets)
			if e != nil {
				logrus.WithError(e).Error("Failed to analyze sentiment")
				return
			}

			logrus.Println("sentimentRequest", sentimentRequest)
			logrus.Println("sentimentSummary", sentimentSummary)

			//var response struct {
			//	SentimentSummary string `json:"sentiment_summary"`
			//}
			//if err := json.Unmarshal(body, &response); err != nil {
			//	logrus.Errorf("Error parsing response from Claude: %v", err)
			//	return "", fmt.Errorf("error parsing response from Claude: %v", err)
			//}
			//
			//logrus.Infof("Sentiment Summary: %s", response.SentimentSummary)
			//return response.SentimentSummary, nil

		}
	}
}

func storeTweets() error {
	ctx := context.Background()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{os.Getenv("CH_API_URL")},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			// Password: "<password>",
		},
		TLS: &tls.Config{},
	})

	if e := conn.Ping(ctx); err != nil {
		if ex, ok := e.(*clickhouse.Exception); ok {
			logrus.Printf("Exception [%d %s \n%s\n", ex.Code, ex.Message, ex.StackTrace)
		}
		return e
	}
	logrus.Printf("%v", conn)

	return nil
}
