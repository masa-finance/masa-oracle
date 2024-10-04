package twitter

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	ShortSleepDuration = 20 * time.Millisecond
	RateLimitDuration  = time.Hour
	MaxRetries         = 3
)

func ShortSleep() {
	time.Sleep(ShortSleepDuration)
}

func GetRateLimitDuration() time.Duration {
	return RateLimitDuration
}

func Retry[T any](operation func() (T, error), maxAttempts int) (T, error) {
	var zero T
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}
		logrus.Errorf("retry attempt %d failed: %v", attempt, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return zero, fmt.Errorf("operation failed after %d attempts", maxAttempts)
}
