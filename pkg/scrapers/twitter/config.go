package twitter

import (
	"time"
)

const (
	ShortSleepDuration = 20 * time.Millisecond
	RateLimitDuration  = 15 * time.Minute
)

func ShortSleep() {
	time.Sleep(ShortSleepDuration)
}

func GetRateLimitDuration() time.Duration {
	return RateLimitDuration
}
