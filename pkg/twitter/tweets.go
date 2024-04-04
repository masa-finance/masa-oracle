package twitter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anthdm/hollywood/actor"
	_ "github.com/lib/pq"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

var sentimentCh = make(chan string)

// TweetRequest represents the parameters for a request to fetch and analyze tweets.
// Count specifies the number of tweets to fetch.
// Query is a list of keywords or phrases to search for within tweets.
// Model indicates the language model to use for analyzing the sentiment of the tweets.
type TweetRequest struct {
	Count int      // Number of tweets to fetch
	Query []string // Keywords or phrases to search for
	Model string   // Language model for sentiment analysis
}

// Manager is a struct that maintains a map of worker actors.
// The map keys are pointers to actor.PID (Process Identifiers) and the values are booleans indicating the worker's availability.
type Manager struct {
	workers map[*actor.PID]bool // Map of worker actors to their availability status
}

// Worker represents a worker entity capable of processing TweetRequests.
// It embeds TweetRequest to inherit its fields and adds a pid field to hold the worker's process identifier.
type Worker struct {
	TweetRequest            // Embedding TweetRequest to inherit its fields
	pid          *actor.PID // pid is the process identifier for the worker
}

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
		"username": os.Getenv("TWITTER_USER"),
	}).Debug("Login successful")

	return scraper
}

// NewManager creates and returns a new actor.Producer function.
// This function, when invoked, will produce an instance of a Manager actor receiver.
// The Manager is responsible for managing a pool of worker actors that process tweet requests.
func NewManager() actor.Producer {
	return func() actor.Receiver {
		return &Manager{
			workers: make(map[*actor.PID]bool),
		}
	}
}

// Receive processes messages sent to the Manager actor. It handles different types of messages
// such as TweetRequest, actor.Started, and actor.Stopped by switching over the type of the message.
// For TweetRequest messages, it delegates the handling to the handleTweetRequest method.
// For actor.Started and actor.Stopped messages, it logs the state of the actor engine.
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

// handleTweetRequest processes a TweetRequest message by iterating over each tweet query within the message.
// For each tweet query, if the worker actor for the current PID does not exist, it spawns a new child actor
// as a TweetWorker to handle the tweet scraping and analysis. It then marks the worker as active in the manager's worker map.
// Parameters:
// - c: The actor context, providing access to actor system functionalities.
// - msg: The TweetRequest containing the details of the tweet queries to process.
// Returns:
// - An error if any issues occur during the handling of the tweet request. Currently, it always returns nil.
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

// NewTweetWorker creates and returns a new actor.Producer function specific for tweet workers.
// This function, when invoked, will produce an instance of a Worker actor receiver, initialized with the provided TweetRequest and PID.
// Parameters:
// - t: A pointer to a TweetRequest struct containing the details of the tweet queries to process.
// - pid: A pointer to an actor.PID representing the process identifier for the actor system.
// Returns:
// - An actor.Producer function that produces Worker actor receivers.
func NewTweetWorker(t *TweetRequest, pid *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Worker{TweetRequest{
			Count: t.Count,
			Query: t.Query,
			Model: t.Model,
		}, pid}
	}
}

// Receive processes messages sent to the Worker actor. It handles different types of messages
// such as actor.Started and actor.Stopped by switching over the type of the message.
// For actor.Started messages, it initiates the process of scraping tweets based on the query and count specified in the Worker,
// analyzes the sentiment of the scraped tweets, and then stops the worker actor.
// For actor.Stopped messages, it simply logs that the worker has stopped.
func (w *Worker) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Started:

		logrus.Infof("Worker started with pid %+v", c.PID())
		tweets, err := ScrapeTweetsByQuery(w.Query[0], w.Count)
		if err != nil {
			logrus.Errorf("ScrapeTweetsByQuery worker error %v", err)
		}
		_, sentimentSummary, err := llmbridge.AnalyzeSentiment(tweets, w.Model)
		if err != nil {
			sentimentCh <- err.Error()
		}
		sentimentCh <- sentimentSummary
		c.Engine().Poison(c.PID()).Wait() // stop this worker by pid when job is complete
	case actor.Stopped:
		logrus.Info("Worker stopped")
	}
}

// ScrapeTweetsUsingActors initiates the process of scraping tweets based on a given query, count, and model.
// It leverages actor-based concurrency to manage the scraping and analysis tasks.
// The function spawns a new actor engine and sends a TweetRequest message to the Manager actor.
// It then waits for a sentiment analysis result to be sent back through a channel.
// Parameters:
//   - query: The search query for fetching tweets.
//   - count: The number of tweets to fetch and analyze.
//   - model: The sentiment analysis model to use.
//     -- claude-3-opus-20240229
//     -- claude-3-sonnet-20240229
//     -- claude-3-haiku-20240307
//     -- gpt-4
//     -- gpt-4-turbo-preview
//     -- gpt-3.5-turbo
//
// Returns:
// - A string containing the sentiment analysis summary.
// - An error if the process fails at any point.
func ScrapeTweetsUsingActors(query string, count int, model string) (string, error) {
	done := make(chan bool)
	var err error
	var engine *actor.Engine

	go func() {
		engine, err = actor.NewEngine(actor.NewEngineConfig())
		if err != nil {
			logrus.Errorf("new actor engine error %v", err)
		} else {
			pid := engine.Spawn(NewManager(), "Manager")
			time.Sleep(time.Millisecond * 200)
			logrus.Infof("Started new actor engine with pid %v \n", pid)
			engine.Send(pid, TweetRequest{Count: count, Query: []string{query}, Model: model})
		}
	}()

	for {
		select {
		case sentiment := <-sentimentCh:
			{
				return sentiment, nil
			}
		case <-done:
			engine = nil
		}
	}
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
