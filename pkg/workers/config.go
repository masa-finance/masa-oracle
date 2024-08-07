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
	WorkerTimeout:         500 * time.Millisecond, // 500 milliseconds
	WorkerResponseTimeout: 250 * time.Millisecond, // 250 milliseconds
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
