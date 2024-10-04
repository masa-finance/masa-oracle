package twitter

import "time"

const (
	ShortSleepDuration = 100 * time.Millisecond
	RateLimitDuration  = time.Hour
)

func ShortSleep() {
	time.Sleep(ShortSleepDuration)
}

func GetRateLimitDuration() time.Duration {
	return RateLimitDuration
}
