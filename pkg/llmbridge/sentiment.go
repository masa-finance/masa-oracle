package llmbridge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gotd/td/tg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/ollama/ollama/api"
	"github.com/sirupsen/logrus"
)

// AnalyzeSentimentTweets analyzes the sentiment of the provided tweets by sending them to the Claude API.
// It concatenates the tweets, creates a payload, sends a request to Claude, parses the response,
// and returns the concatenated tweets content, a sentiment summary, and any error.
func AnalyzeSentimentTweets(tweets []*twitterscraper.TweetResult, model string, prompt string) (string, string, error) {
	// check if we are using claude or gpt, can add others easily
	if strings.Contains(model, "claude-") {
		client := NewClaudeClient() // Adjusted to call without arguments

		var validTweets []*twitterscraper.TweetResult
		for _, tweet := range tweets {
			if tweet.Error != nil {
				logrus.WithError(tweet.Error).Warn("Error in tweet")
				continue
			}
			validTweets = append(validTweets, tweet)
		}

		tweetsContent := ConcatenateTweets(validTweets)
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
func ConcatenateTweets(tweets []*twitterscraper.TweetResult) string {
	var tweetsTexts []string
	for _, t := range tweets {
		tweetsTexts = append(tweetsTexts, t.Tweet.Text)
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
	} else if strings.HasPrefix(model, "@") {
		genReq := api.ChatRequest{
			Model: model,
			Messages: []api.Message{
				{Role: "user", Content: data},
				{Role: "assistant", Content: prompt},
			},
		}

		requestJSON, err := json.Marshal(genReq)
		if err != nil {
			return "", "", err
		}
		cfUrl := config.GetInstance().LLMCfUrl
		if cfUrl == "" {
			return "", "", errors.New("cloudflare workers url not set")
		}
		uri := fmt.Sprintf("%s%s", cfUrl, model)
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
	} else if strings.HasPrefix(model, "ollama/") {
		stream := false

		genReq := api.ChatRequest{
			Model: strings.TrimPrefix(model, "ollama/"),
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
	} else {
		return "", "", errors.New("model not supported")
	}
}

// AnalyzeSentimentDiscord analyzes the sentiment of the provided Discord messages by sending them to the sentiment analysis API.
// It concatenates the messages, creates a payload, sends a request to the sentiment analysis service, parses the response,
// and returns the concatenated messages content, a sentiment summary, and any error.
func AnalyzeSentimentDiscord(messages []string, model string, prompt string) (string, string, error) {
	// Concatenate messages with a newline character
	messagesContent := strings.Join(messages, "\n")

	// The rest of the code follows the same pattern as AnalyzeSentimentTweets
	// Replace with the actual logic you have for sending requests to your sentiment analysis service
	// For example, if you're using the Claude API:
	if strings.Contains(model, "claude-") {
		client := NewClaudeClient() // Adjusted to call without arguments
		payloadBytes, err := CreatePayload(messagesContent, model, prompt)
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
		return messagesContent, sentimentSummary, nil

	} else {
		stream := false

		genReq := api.ChatRequest{
			Model: model,
			Messages: []api.Message{
				{Role: "user", Content: messagesContent},
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
			logrus.Errorf("Error marshaling request JSON: %v", err)
			return "", "", err
		}
		uri := config.GetInstance().LLMChatUrl
		if uri == "" {
			errMsg := "ollama api url not set"
			logrus.Errorf("%v", errMsg)
			return "", "", errors.New(errMsg)
		}
		resp, err := http.Post(uri, "application/json", bytes.NewReader(requestJSON))
		if err != nil {
			logrus.Errorf("Error sending request to API: %v", err)
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Error reading response body: %v", err)
			return "", "", err
		}

		var payload api.ChatResponse
		err = json.Unmarshal(body, &payload)
		if err != nil {
			logrus.Errorf("Error unmarshaling response JSON: %v", err)
			return "", "", err
		}

		sentimentSummary := payload.Message.Content
		return messagesContent, SanitizeResponse(sentimentSummary), nil
	}
}

// AnalyzeSentimentTelegram analyzes the sentiment of the provided Telegram messages by sending them to the sentiment analysis API.
func AnalyzeSentimentTelegram(messages []*tg.Message, model string, prompt string) (string, string, error) {
	// Concatenate messages with a newline character
	var messageTexts []string
	for _, msg := range messages {
		if msg != nil {
			messageTexts = append(messageTexts, msg.Message)
		}
	}
	messagesContent := strings.Join(messageTexts, "\n")

	// The rest of the code follows the same pattern as AnalyzeSentimentDiscord
	if strings.Contains(model, "claude-") {
		client := NewClaudeClient() // Adjusted to call without arguments
		payloadBytes, err := CreatePayload(messagesContent, model, prompt)
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
		return messagesContent, sentimentSummary, nil

	} else {
		stream := false

		genReq := api.ChatRequest{
			Model: model,
			Messages: []api.Message{
				{Role: "user", Content: messagesContent},
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
			logrus.Errorf("Error marshaling request JSON: %v", err)
			return "", "", err
		}
		uri := config.GetInstance().LLMChatUrl
		if uri == "" {
			errMsg := "ollama api url not set"
			logrus.Errorf(errMsg)
			return "", "", errors.New(errMsg)
		}
		resp, err := http.Post(uri, "application/json", bytes.NewReader(requestJSON))
		if err != nil {
			logrus.Errorf("Error sending request to API: %v", err)
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Error reading response body: %v", err)
			return "", "", err
		}

		var payload api.ChatResponse
		err = json.Unmarshal(body, &payload)
		if err != nil {
			logrus.Errorf("Error unmarshaling response JSON: %v", err)
			return "", "", err
		}

		sentimentSummary := payload.Message.Content
		return messagesContent, SanitizeResponse(sentimentSummary), nil
	}
}
