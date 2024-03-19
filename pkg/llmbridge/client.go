package llmbridge

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ClaudeClient struct {
	config *ClaudeAPIConfig
}

func NewClaudeClient() *ClaudeClient {
	config := NewClaudeAPIConfig()
	return &ClaudeClient{config: config}
}

func (c *ClaudeClient) SendRequest(payloadBytes []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.config.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", c.config.Version)

	client := &http.Client{}
	return client.Do(req)
}

func ParseResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		SentimentSummary string `json:"sentiment_summary"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.SentimentSummary, nil
}
