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

// NewClaudeClient creates a new ClaudeClient instance with default configuration.
func NewClaudeClient() *ClaudeClient {
	config := NewClaudeAPIConfig()
	return &ClaudeClient{config: config}
}

// SendRequest sends an HTTP request to the Claude API with the given payload.
// It sets the required headers like Content-Type, x-api-key etc.
// Returns the HTTP response and any error.
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

// ParseResponse unmarshals the JSON response body from the Claude API
// into a struct to extract the sentiment summary string.
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
