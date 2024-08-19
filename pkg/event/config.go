package event

import "time"

const (
	// APIVersion is the version of the analytics API
	APIVersion = "v1"

	// DefaultBaseURL is the default URL for the external API
	DefaultBaseURL = "https://api.example.com/analytics"

	// DefaultHTTPTimeout is the default timeout for HTTP requests
	DefaultHTTPTimeout = 10 * time.Second

	// MaxEventsInMemory is the maximum number of events to keep in memory
	MaxEventsInMemory = 1000
)

// Config holds the configuration for the analytics package
type Config struct {
	BaseURL     string
	HTTPTimeout time.Duration
	LogLevel    string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL:     DefaultBaseURL,
		HTTPTimeout: DefaultHTTPTimeout,
		LogLevel:    "info",
	}
}
