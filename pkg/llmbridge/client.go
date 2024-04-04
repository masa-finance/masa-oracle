package llmbridge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type ClaudeClient struct {
	config *ClaudeAPIConfig
}

type GPTClient struct {
	config *GPTAPIConfig
}

// NewClaudeClient creates a new ClaudeClient instance with default configuration.
func NewClaudeClient() *ClaudeClient {
	cnf := NewClaudeAPIConfig()
	return &ClaudeClient{config: cnf}
}

// NewGPTClient creates a new GPTClient instance with default configuration.
func NewGPTClient() *GPTClient {
	cnf := NewGPTConfig()
	return &GPTClient{config: cnf}
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

func (c *GPTClient) SendRequest(tweetsContent string, model string, prompt string) (string, error) {
	var openAiModel string
	switch model {
	case "gpt-4":
		openAiModel = openai.GPT4
	case "gpt-4-turbo-preview":
		openAiModel = openai.GPT40613
	case "gpt-3.5-turbo":
		openAiModel = openai.GPT3Dot5Turbo
	default:
		break
	}

	cfg := config.GetInstance()
	key := cfg.GPTApiKey
	if key == "" {
		return "", errors.New("OPENAI_API_KEY is not set")
	}
	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openAiModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Please perform a sentiment analysis on the following tweets, using an unbiased approach. Sentiment analysis involves identifying and categorizing opinions expressed in text, particularly to determine whether the writer's attitude towards a particular topic, product, etc., is positive, negative, or neutral. After analyzing, please provide a summary of the overall sentiment expressed in these tweets, including the proportion of positive, negative, and neutral sentiments if applicable.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: tweetsContent,
				},
			},
		},
	)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
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

// SanitizeResponse removes non-ASCII characters and unnecessary whitespace from a string.
// It also strips away double quotes for cleaner presentation.
// Parameters:
// - str: The input string to be sanitized.
// Returns: A sanitized string with only ASCII characters, reduced whitespace, and no double quotes.
func SanitizeResponse(str string) string {
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
	if response.Content != nil {
		for _, t := range response.Content {
			summary = SanitizeResponse(t.Text)
		}
	} else {
		var responseError map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseError); err == nil {
			if errVal, ok := responseError["error"].(map[string]interface{}); ok {
				if message, ok := errVal["message"].(string); ok {
					summary = fmt.Sprintf("error from llm: Service %v", message)
				}
			}
		}
		return summary, nil
	}
	return summary, nil
}
