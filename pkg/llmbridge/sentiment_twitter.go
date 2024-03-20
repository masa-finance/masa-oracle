package llmbridge

import (
	"strings"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// AnalyzeSentiment analyzes the sentiment of the provided tweets by sending them to the Claude API.
// It concatenates the tweets, creates a payload, sends a request to Claude, parses the response,
// and returns the concatenated tweets content, a sentiment summary, and any error.
func AnalyzeSentiment(tweets []*twitterscraper.Tweet) (string, string, error) {
	client := NewClaudeClient() // Adjusted to call without arguments

	tweetsContent := ConcatenateTweets(tweets)
	payloadBytes, err := CreatePayload(tweetsContent)
	if err != nil {
		logrus.Errorf("Error creating payload: %v", err)
		return "", "", err
	}

	resp, err := client.SendRequest(payloadBytes)
	if err != nil {
		logrus.Errorf("Error sending request to Claude API: %v", err)
		return "", "", err
	}
	defer resp.Body.Close()

	sentimentSummary, err := ParseResponse(resp)
	if err != nil {
		logrus.Errorf("Error parsing response from Claude: %v", err)
		return "", "", err
	}

	return tweetsContent, sentimentSummary, nil
}

// ConcatenateTweets concatenates the text of the provided tweets into a single string,
// with each tweet separated by a newline character.
func ConcatenateTweets(tweets []*twitterscraper.Tweet) string {
	var tweetsTexts []string
	for _, tweet := range tweets {
		tweetsTexts = append(tweetsTexts, tweet.Text)
	}
	return strings.Join(tweetsTexts, "\n")
}
