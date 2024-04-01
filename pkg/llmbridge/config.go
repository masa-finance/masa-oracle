package llmbridge

import "github.com/masa-finance/masa-oracle/pkg/config"

type ClaudeAPIConfig struct {
	URL     string
	APIKey  string
	Version string
}

type GPTAPIConfig struct {
	APIKey string
}

// NewClaudeAPIConfig creates a new ClaudeAPIConfig instance with values loaded
// from the application config.
func NewClaudeAPIConfig() *ClaudeAPIConfig {
	appConfig := config.GetInstance()

	// need to add these to the config package
	return &ClaudeAPIConfig{
		URL:     appConfig.ClaudeApiURL,
		APIKey:  appConfig.ClaudeApiKey,
		Version: appConfig.ClaudeApiVersion,
	}
}

// NewGPTConfig creates a new GPTConfig instance with values loaded
// from the application config.
func NewGPTConfig() *GPTAPIConfig {
	appConfig := config.GetInstance()

	// need to add these to the config package
	return &GPTAPIConfig{
		APIKey: appConfig.GPTApiKey,
	}
}
