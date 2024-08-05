package workers

import (
	"time"
)

type WorkerConfig struct {
	WorkerTimeout    time.Duration
	MaxRetries       int
	MaxSpawnAttempts int
	WorkerBufferSize int
}

var DefaultConfig = WorkerConfig{
	WorkerTimeout:    30 * time.Second,
	MaxRetries:       3,
	MaxSpawnAttempts: 3,
	WorkerBufferSize: 100,
}

func LoadConfig() (*WorkerConfig, error) {
	// For now, we'll just return the default config
	config := DefaultConfig
	return &config, nil
}
