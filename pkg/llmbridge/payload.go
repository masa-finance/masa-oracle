package llmbridge

import (
	"encoding/json"
)

type Payload struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	System      string    `json:"system"`
	Messages    []Message `json:"messages"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// CreatePayload generates a JSON payload for the OpenAI API from the given
// tweetsContent string. This payload configures the model, max tokens,
// temperature and prompt to analyze the sentiment of the tweets without
// bias and summarize the overall sentiment.
func CreatePayload(tweetsContent string, model string) ([]byte, error) {
	payload := Payload{
		Model:       model,
		MaxTokens:   4000,
		Temperature: 0,
		System:      "Please perform a sentiment analysis on the following tweets, using an unbiased approach. Sentiment analysis involves identifying and categorizing opinions expressed in text, particularly to determine whether the writer's attitude towards a particular topic, product, etc., is positive, negative, or neutral. After analyzing, please provide a summary of the overall sentiment expressed in these tweets, including the proportion of positive, negative, and neutral sentiments if applicable.",
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: tweetsContent,
					},
				},
			},
		},
	}
	return json.Marshal(payload)
}
