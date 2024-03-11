package llmbridge

import "os"

type ClaudeAPIConfig struct {
	URL     string
	APIKey  string
	Version string
}

func NewClaudeAPIConfig() *ClaudeAPIConfig {
	return &ClaudeAPIConfig{
		URL:     "https://api.anthropic.com/v1/messages",
		APIKey:  os.Getenv("CLAUDE_API_KEY"),
		Version: "2023-06-01",
	}
}
