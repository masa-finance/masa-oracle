package workers

import (
	"time"

	"github.com/sirupsen/logrus"
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
	WorkerTimeout:         45 * time.Second,
	WorkerResponseTimeout: 35 * time.Second,
	ConnectionTimeout:     500 * time.Millisecond,
	MaxRetries:            1,
	MaxSpawnAttempts:      1,
	WorkerBufferSize:      100,
	MaxRemoteWorkers:      25,
}

var workerConfig *WorkerConfig

func init() {
	var err error
	workerConfig, err = LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load worker config: %v", err)
	}
}

func LoadConfig() (*WorkerConfig, error) {
	// For now, we'll just return the default config
	config := DefaultConfig
	return &config, nil
}
