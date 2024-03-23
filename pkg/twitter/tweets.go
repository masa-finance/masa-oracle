package twitter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	_ "github.com/lib/pq"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

type Sentiment struct {
	ID             int64           `json:"id"`
	ConversationId int64           `json:"conversation_id"`
	Tweet          json.RawMessage `json:"tweet"`
	PromptId       int64           `json:"prompt_id"`
}

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
		"username": os.Getenv("TWITTER_USER"),
	}).Debug("Login successful")

	return scraper
}

func ScrapeTweetsByQuery(query string, count int) ([]*twitterscraper.Tweet, error) {
	scraper := auth()
	var tweets []*twitterscraper.Tweet

	if scraper == nil {
		return nil, fmt.Errorf("scraper instance is nil")
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

func Scrape(query string, count int) ([]*twitterscraper.Tweet, error) {
	rowChan := make(chan []*twitterscraper.Tweet)
	scraper := auth()
	go scrapeTweetsToChannel(scraper, query, count, rowChan)

	var wg sync.WaitGroup
	var size int64
	var numWorkers = 5

	//conn, err := connectToClickHouse()
	//if err != nil {
	//	logrus.Errorf("clickhouse connect err %s", err)
	//}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processTweets(&wg, rowChan, &size)
	}
	wg.Wait()
	logrus.Println("size", size)

	return nil, nil
}

func processTweets(wg *sync.WaitGroup, rowChan chan []*twitterscraper.Tweet, size *int64) {
	defer wg.Done()

	for rows := range rowChan {
		if rows != nil {
			logrus.Println("rows", rows)
			for _, row := range rows {
				atomic.AddInt64(size, 1) // testing counts
				logrus.Printf("row ===> %v\n", row)
			}
			deserializedTweets, err := serializeTweets(rows)
			if err != nil {
				logrus.WithError(err).Error("Failed to deserialize tweets data")
				return
			}
			// also getting sentiment request to save to datastore
			sentimentRequest, sentimentSummary, e := llmbridge.AnalyzeSentiment(deserializedTweets)
			if e != nil {
				logrus.WithError(e).Error("Failed to analyze sentiment")
				return
			}

			logrus.Println("sentimentRequest", sentimentRequest)
			logrus.Println("sentimentSummary", sentimentSummary)
		}
	}
}

// connectToClickHouse tests
func connectToClickHouse() (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{os.Getenv("CH_URL")},
			Auth: clickhouse.Auth{
				Database: "masa",
				Username: "default",
				// Password: "<password>",
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if pe := conn.Ping(ctx); pe != nil {
		var ex *clickhouse.Exception
		if errors.As(pe, &ex) {
			logrus.Errorf("Exception [%d %s \n%s\n", ex.Code, ex.Message, ex.StackTrace)
		}
		return nil, pe
	}

	err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS sentiment (id UInt64, conversation_id UInt32, tweets String, prompt_id Uint32, PRIMARY KEY (id)) ENGINE = MergeTree() ORDER BY id`)
	if err != nil {
		return nil, err
	}

	batch, e := conn.PrepareBatch(ctx, `INSERT INTO sentiment`)
	if e != nil {
		return nil, e
	}

	if bae := batch.Append(uint64(1), uint32(1), "hello", uint32(1)); bae != nil {
		return nil, bae
	}

	if be := batch.Send(); be != nil {
		return nil, be
	}

	rows, qe := conn.Query(ctx, `SELECT conversation_id, tweet FROM sentiment`)
	if qe != nil {
		return nil, qe
	}
	for rows.Next() {
		var (
			conversation_id uint32
			tweet           string
		)
		if se := rows.Scan(&conversation_id, &tweet); err != nil {
			return nil, se
		}
		logrus.Printf("row: convesation_id=%d, tweet=%s\n", conversation_id, tweet)
	}
	_ = rows.Close()

	return conn, rows.Err()

	// return conn, nil
}

// ingestTweets tests
func ingestTweets(wg *sync.WaitGroup, rowChan <-chan []*twitterscraper.Tweet, conn driver.Conn, batchSize int) {
	defer wg.Done()

	newBatch := func() driver.Batch {
		ctx := context.Background()
		batch, err := conn.PrepareBatch(ctx, `INSERT INTO sentiment (conversation_id, tweet)`)
		if err != nil {
			logrus.Errorf("err %v", err)
		}
		return batch
	}
	batch := newBatch()
	tweetsProcessed := 0
	for row := range rowChan {
		conversationID := row[0]
		body := row[1]

		err := batch.Append(conversationID, body)
		if err != nil {
			logrus.Errorf("err %v", err)
		}
		tweetsProcessed++
		if tweetsProcessed%tweetsProcessed == 0 {
			if err := batch.Send(); err != nil {
				logrus.Errorf("%v", err)
			}
			batch = newBatch()
		}
	}
	if err := batch.Send(); err != nil {
		logrus.Errorf("%v", err)
	}
}
