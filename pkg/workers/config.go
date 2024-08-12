package workers

import (
	"time"
)

type WorkerConfig struct {
	WorkerTimeout         time.Duration
	WorkerResponseTimeout time.Duration
	ConnectionTimeout     time.Duration
	MaxRetries            int
	MaxSpawnAttempts      int
	WorkerBufferSize      int
	MaxRemoteWorkers      int
}

var DefaultConfig = WorkerConfig{
	WorkerTimeout:         55 * time.Second,
	WorkerResponseTimeout: 30 * time.Second,
	ConnectionTimeout:     1 * time.Second,
	MaxRetries:            1,
	MaxSpawnAttempts:      1,
	WorkerBufferSize:      100,
	MaxRemoteWorkers:      1,
}

func LoadConfig() (*WorkerConfig, error) {
	// For now, we'll just return the default config
	config := DefaultConfig
	return &config, nil
}
