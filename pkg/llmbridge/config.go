package llmbridge

import "github.com/masa-finance/masa-oracle/pkg/config"

type ClaudeAPIConfig struct {
	URL     string
	APIKey  string
	Version string
}

func NewClaudeAPIConfig() *ClaudeAPIConfig {
	appConfig := config.GetInstance()

	// need to add these to the config package
	return &ClaudeAPIConfig{
		URL:     appConfig.ClaudeApiURL,
		APIKey:  appConfig.ClaudeApiKey,
		Version: appConfig.ClaudeApiVersion,
	}
}
