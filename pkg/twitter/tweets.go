package twitter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anthdm/hollywood/actor"
	_ "github.com/lib/pq"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

var sentimentCh = make(chan string)

type TweetRequest struct {
	Count int
	Query []string
}

type Manager struct {
	workers map[*actor.PID]bool
}

type Worker struct {
	TweetRequest
	pid *actor.PID
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

func NewManager() actor.Producer {
	return func() actor.Receiver {
		return &Manager{
			workers: make(map[*actor.PID]bool),
		}
	}
}

func (m *Manager) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case TweetRequest:
		_ = m.handleTweetRequest(c, msg)
	case actor.Started:
		logrus.Info("Actor engine initialized")
	case actor.Stopped:
		logrus.Info("Actor engine stopped")
	}
}

func (m *Manager) handleTweetRequest(c *actor.Context, msg TweetRequest) error {
	for i, tweet := range msg.Query {
		if _, ok := m.workers[c.PID()]; !ok {
			logrus.Debugf("Tweet %+v with %v", tweet, c.PID())
			c.SpawnChild(NewTweetWorker(&msg, c.PID()), fmt.Sprintf("worker/%d", i))
			m.workers[c.PID()] = true
		}
	}
	return nil
}

func NewTweetWorker(t *TweetRequest, pid *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Worker{TweetRequest{
			Count: t.Count,
			Query: t.Query,
		}, pid}
	}
}

func (w *Worker) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Started:

		logrus.Infof("Worker started with pid %+v", c.PID())
		tweets, err := ScrapeTweetsByQuery(w.Query[0], w.Count)
		if err != nil {
			logrus.Errorf("ScrapeTweetsByQuery worker error %v", err)
		}
		_, sentimentSummary, err := llmbridge.AnalyzeSentiment(tweets)
		if err != nil {
			sentimentCh <- err.Error()
		}
		sentimentCh <- sentimentSummary
		c.Engine().Poison(c.PID()) // stop this worker by pid
	case actor.Stopped:
		logrus.Info("Worker stopped")
	}
}

func ScrapeTweetsUsingActors(query string, count int) (string, error) {

	done := make(chan bool)

	go func() {
		// Concurrent Actor Framework ref: https://en.wikipedia.org/wiki/Actor_model
		e, err := actor.NewEngine(actor.NewEngineConfig())
		if err != nil {
			logrus.Errorf("%v", err)
		} else {
			pid := e.Spawn(NewManager(), "Manager")
			time.Sleep(time.Millisecond * 200)
			logrus.Printf("Started new actor engine with pid %v \n", pid)
			e.Send(pid, TweetRequest{Count: count, Query: []string{query}})
		}

	}()

	for {
		select {
		case sentiment := <-sentimentCh:
			{
				return sentiment, nil
			}
		case <-done:
			break
		}
	}

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

func ScrapeTweetsByQueryWaitGroups(query string, count int) ([]*twitterscraper.Tweet, error) {
	rowChan := make(chan []*twitterscraper.Tweet)
	scraper := auth()
	go scrapeTweetsToChannel(scraper, query, count, rowChan)

	var wg sync.WaitGroup
	var numWorkers = 5

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processTweets(&wg, rowChan)
	}
	wg.Wait()
	// @todo pass results to topic
	return nil, nil
}

func processTweets(wg *sync.WaitGroup, rowChan chan []*twitterscraper.Tweet) {
	defer wg.Done()

	for rows := range rowChan {
		if rows != nil {
			logrus.Println("rows", rows)
			for _, row := range rows {
				// @todo do we want to process each tweet or all tweets
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
