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
func CreatePayload(tweetsContent string, model string, prompt string) ([]byte, error) {
	payload := Payload{
		Model:       model,
		MaxTokens:   4000,
		Temperature: 0,
		System:      prompt,
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
