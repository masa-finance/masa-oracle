package twitter

import (
	"fmt"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"

	"github.com/sirupsen/logrus"
)

func ScrapeFollowersForProfile(username string, count int) ([]*twitterscraper.Profile, error) {
	scraper, account, err := getAuthenticatedScraper()
	if err != nil {
		return nil, err
	}

	followingResponse, errString, _ := scraper.FetchFollowers(username, count, "")
	if errString != "" {
		if handleRateLimit(fmt.Errorf(errString), account) {
			return nil, fmt.Errorf("rate limited")
		}
		logrus.Errorf("Error fetching followers: %v", errString)
		return nil, fmt.Errorf("%v", errString)
	}
	account.LastScraped = time.Now()
	return followingResponse, nil
}
