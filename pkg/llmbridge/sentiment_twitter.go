package llmbridge

import (
	"strings"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

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

func ConcatenateTweets(tweets []*twitterscraper.Tweet) string {
	var tweetsTexts []string
	for _, tweet := range tweets {
		tweetsTexts = append(tweetsTexts, tweet.Text)
	}
	return strings.Join(tweetsTexts, "\n")
}
