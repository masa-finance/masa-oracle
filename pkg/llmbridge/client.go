package llmbridge

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

type Response struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	Role         string            `json:"role"`
	Content      []ResponseContent `json:"content"`
	Model        string            `json:"model"`
	StopReason   string            `json:"stop_reason"`
	StopSequence *string           `json:"stop_sequence"` // Use *string for nullable fields
	Usage        Usage             `json:"usage"`
}

type ResponseContent struct {
	Type  string         `json:"type"`
	Text  string         `json:"text,omitempty"`
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// sanitizeResponse removes non-ASCII characters and unnecessary whitespace from a string.
// It also strips away double quotes for cleaner presentation.
// Parameters:
// - str: The input string to be sanitized.
// Returns: A sanitized string with only ASCII characters, reduced whitespace, and no double quotes.
func sanitizeResponse(str string) string {
	var result []rune
	for _, r := range str {
		if r >= 0 && r <= 127 {
			result = append(result, r)
		}
	}
	sanitizedString := string(result)
	sanitizedString = strings.ReplaceAll(sanitizedString, "\n\n", " ")
	sanitizedString = strings.ReplaceAll(sanitizedString, "\n", "")
	sanitizedString = strings.ReplaceAll(sanitizedString, "\"", "")
	return sanitizedString
}

// ParseResponse takes an http.Response, reads its body, and attempts to unmarshal it into a Response struct.
// It then sanitizes the text content of each ResponseContent within the Response and returns a summary string.
// Parameters:
// - resp: A pointer to an http.Response object that contains the server's response to an HTTP request.
// Returns:
// - A string that represents a sanitized summary of the response content.
// - An error if reading the response body or unmarshalling fails.
func ParseResponse(resp *http.Response) (string, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response Response
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", err
	}
	var summary = ""
	for _, t := range response.Content {
		summary = sanitizeResponse(t.Text)
	}
	return summary, nil
}
