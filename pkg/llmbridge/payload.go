package llmbridge

import "encoding/json"

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

func CreatePayload(tweetsContent string) ([]byte, error) {
	payload := Payload{
		Model:       "claude-3-opus-20240229",
		MaxTokens:   4000,
		Temperature: 0,
		System:      "Please analyze the sentiment of the following tweets without bias and summarize the overall sentiment:",
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
