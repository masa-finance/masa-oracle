package api

import "time"

// APIConfig contains configuration settings for the API
type APIConfig struct {
	WorkerResponseTimeout time.Duration
	// Add other API-specific configuration fields here
}

var DefaultConfig = APIConfig{
	WorkerResponseTimeout: 120 * time.Second,
	// Set default values for other fields here
}

// LoadConfig loads the API configuration
// This can be expanded later to load from environment variables or a file
func LoadConfig() (*APIConfig, error) {
	// For now, we'll just return the default config
	config := DefaultConfig
	return &config, nil
}
