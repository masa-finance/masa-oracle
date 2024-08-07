package workers

import (
	"time"
)

type WorkerConfig struct {
	ConnectTimeout        time.Duration
	WorkerResponseTimeout time.Duration
	MaxRetries            int
	MaxSpawnAttempts      int
	WorkerBufferSize      int
	MaxRemoteWorkers      int
}

var DefaultConfig = WorkerConfig{
	ConnectTimeout:        250 * time.Millisecond,
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
