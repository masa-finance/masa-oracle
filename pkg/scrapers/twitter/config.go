package twitter

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	minSleepDuration  = 500 * time.Millisecond
	maxSleepDuration  = 2 * time.Second
	RateLimitDuration = 15 * time.Minute
)

var (
	rng *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomSleep() {
	duration := minSleepDuration + time.Duration(rng.Int63n(int64(maxSleepDuration-minSleepDuration)))
	logrus.Debugf("Sleeping for %v", duration)
	time.Sleep(duration)
}

func GetRateLimitDuration() time.Duration {
	return RateLimitDuration
}
