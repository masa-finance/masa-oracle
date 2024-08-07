package workers

import (
	"time"
)

type WorkerConfig struct {
	WorkerTimeout         time.Duration
	WorkerResponseTimeout time.Duration
	MaxRetries            int
	MaxSpawnAttempts      int
	WorkerBufferSize      int
	MaxRemoteWorkers      int
}

var DefaultConfig = WorkerConfig{
	WorkerTimeout:         30 * time.Second,
	WorkerResponseTimeout: 25 * time.Second,
	MaxRetries:            1,
	MaxSpawnAttempts:      3,
	WorkerBufferSize:      100,
	MaxRemoteWorkers:      1,
}

func LoadConfig() (*WorkerConfig, error) {
	// For now, we'll just return the default config
	config := DefaultConfig
	return &config, nil
}
