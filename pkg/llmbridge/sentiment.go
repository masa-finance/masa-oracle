package llmbridge

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

// AnalyzeSentimentTweets analyzes the sentiment of the provided tweets by sending them to the Claude API.
// It concatenates the tweets, creates a payload, sends a request to Claude, parses the response,
// and returns the concatenated tweets content, a sentiment summary, and any error.
func AnalyzeSentimentTweets(tweets []*twitterscraper.Tweet, model string, prompt string) (string, string, error) {
	// check if we are using claude or gpt, can add others easily
	if strings.Contains(model, "claude-") {
		client := NewClaudeClient() // Adjusted to call without arguments
		tweetsContent := ConcatenateTweets(tweets)
		payloadBytes, err := CreatePayload(tweetsContent, model, prompt)
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

	} else if strings.Contains(model, "gpt-") {
		client := NewGPTClient()
		tweetsContent := ConcatenateTweets(tweets)
		sentimentSummary, err := client.SendRequest(tweetsContent, model, prompt)
		if err != nil {
			logrus.Errorf("Error sending request to GPT: %v", err)
			return "", "", err
		}
		return tweetsContent, sentimentSummary, nil
	} else {
		stream := false
		tweetsContent := ConcatenateTweets(tweets)

		genReq := api.ChatRequest{
			Model: model,
			Messages: []api.Message{
				{Role: "user", Content: tweetsContent},
				{Role: "assistant", Content: prompt},
			},
			Stream: &stream,
			Options: map[string]interface{}{
				"temperature": 0.0,
				"seed":        42,
				"num_ctx":     4096,
			},
		}

		requestJSON, err := json.Marshal(genReq)
		if err != nil {
			return "", "", err
		}
		uri := config.GetInstance().LLMChatUrl
		if uri == "" {
			return "", "", errors.New("ollama api url not set")
		}
		resp, err := http.Post(uri, "application/json", bytes.NewReader(requestJSON))
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}

		var payload api.ChatResponse
		err = json.Unmarshal(body, &payload)
		if err != nil {
			return "", "", err
		}

		sentimentSummary := payload.Message.Content
		return tweetsContent, SanitizeResponse(sentimentSummary), nil
	}

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

// AnalyzeSentimentWeb analyzes the sentiment of the provided web page text data by sending them to the Claude API.
// It concatenates the text, creates a payload, sends a request to Claude, parses the response,
// and returns the concatenated content, a sentiment summary, and any error.
func AnalyzeSentimentWeb(data string, model string, prompt string) (string, string, error) {
	// check if we are using claude or gpt, can add others easily
	if strings.Contains(model, "claude-") {
		client := NewClaudeClient() // Adjusted to call without arguments
		payloadBytes, err := CreatePayload(data, model, prompt)
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
		return data, sentimentSummary, nil

	} else if strings.Contains(model, "gpt-") {
		client := NewGPTClient()
		sentimentSummary, err := client.SendRequest(data, model, prompt)
		if err != nil {
			logrus.Errorf("Error sending request to GPT: %v", err)
			return "", "", err
		}
		return data, sentimentSummary, nil
	} else {
		stream := false

		genReq := api.ChatRequest{
			Model: model,
			Messages: []api.Message{
				{Role: "user", Content: data},
				{Role: "assistant", Content: prompt},
			},
			Stream: &stream,
			Options: map[string]interface{}{
				"temperature": 0.0,
				"seed":        42,
				"num_ctx":     4096,
			},
		}

		requestJSON, err := json.Marshal(genReq)
		if err != nil {
			return "", "", err
		}
		uri := config.GetInstance().LLMChatUrl
		if uri == "" {
			return "", "", errors.New("ollama api url not set")
		}
		resp, err := http.Post(uri, "application/json", bytes.NewReader(requestJSON))
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}

		var payload api.ChatResponse
		err = json.Unmarshal(body, &payload)
		if err != nil {
			return "", "", err
		}

		sentimentSummary := payload.Message.Content
		return data, SanitizeResponse(sentimentSummary), nil
	}
}
