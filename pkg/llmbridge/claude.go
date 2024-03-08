package llmbridge

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/masa-finance/masa-oracle/pkg/config" // Assuming this is where your configuration is managed
	twitterscraper "github.com/n0madic/twitter-scraper"
)

// AnalyzeSentiment takes a slice of tweets, sends them to Claude for sentiment analysis, and returns a summary.
func AnalyzeSentiment(tweets []*twitterscraper.Tweet) (string, error) {
	// Retrieve the configuration instance
	appConfig := config.GetInstance()

	// Retrieve the Claude API key from the configuration
	ClaudeAPIKey := appConfig.ClaudeApiKey

	// Prepare the tweets for analysis
	var tweetsText []string
	for _, tweet := range tweets {
		tweetsText = append(tweetsText, tweet.Text)
	}
	tweetsForAnalysis := strings.Join(tweetsText, "\n")

	// Craft a natural language prompt for Claude
	prompt := fmt.Sprintf("Please analyze the sentiment of the following tweets without bias and summarize the overall sentiment:\n\n%s", tweetsForAnalysis)

	// Prepare the request to Claude API
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":    "claude-3-opus-20240229",
		"messages": []map[string]interface{}{{"role": "user", "content": prompt}},
		"api_key":  ClaudeAPIKey,
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	// Send the request
	resp, err := http.Post("https://api.anthropic.com/claude/analyze", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return "", fmt.Errorf("error sending request to Claude API: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the JSON response from Claude
	// Assuming Claude returns a structured response with a field for the sentiment summary
	var response struct {
		SentimentSummary string `json:"sentiment_summary"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parsing response from Claude: %v", err)
	}

	return response.SentimentSummary, nil
}
