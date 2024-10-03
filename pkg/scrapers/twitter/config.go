package twitter

import "time"

type TwitterConfig struct {
	SleepTime time.Duration
}

func NewTwitterConfig() *TwitterConfig {
	return &TwitterConfig{
		SleepTime: 100 * time.Millisecond,
	}
}
