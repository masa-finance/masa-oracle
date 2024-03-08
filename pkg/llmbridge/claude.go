package llmbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// AnalyzeSentiment sends tweets to Claude for sentiment analysis and returns the analysis result.
func AnalyzeSentiment(tweets []*twitterscraper.Tweet) (string, error) {
	// Concatenate the text of each tweet into a single string
	var tweetsTexts []string
	for _, tweet := range tweets {
		tweetsTexts = append(tweetsTexts, tweet.Text)
	}
	tweetsContent := strings.Join(tweetsTexts, "\n")

	// Construct the request payload with actual tweets text
	payload := map[string]interface{}{
		"model":       "claude-3-opus-20240229",
		"max_tokens":  4000,
		"temperature": 0,
		"system":      "Please analyze the sentiment of the following tweets without bias and summarize the overall sentiment:",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": tweetsContent,
					},
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Error marshaling payload: %v", err)
		return "", fmt.Errorf("error marshaling payload: %v", err)
	}

	logrus.Infof("Payload for Claude API: %s", string(payloadBytes))

	// Send the request to Claude API
	logrus.Info("Sending request to Claude API")
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(payloadBytes))
	if err != nil {
		logrus.Errorf("Error creating request to Claude API: %v", err)
		return "", fmt.Errorf("error creating request to Claude API: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "-")
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error sending request to Claude API: %v", err)
		return "", fmt.Errorf("error sending request to Claude API: %v", err)
	}
	defer resp.Body.Close()

	logrus.Info("Request to Claude API sent successfully")

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error reading response body: %v", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	logrus.Infof("Response from Claude API: %s", string(body))

	var response struct {
		SentimentSummary string `json:"sentiment_summary"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		logrus.Errorf("Error parsing response from Claude: %v", err)
		return "", fmt.Errorf("error parsing response from Claude: %v", err)
	}

	logrus.Infof("Sentiment Summary: %s", response.SentimentSummary)
	return response.SentimentSummary, nil
}
