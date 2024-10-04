package twitter

import (
	"fmt"

	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

func ScrapeFollowersForProfile(username string, count int) ([]twitterscraper.Legacy, error) {
	return Retry(func() ([]twitterscraper.Legacy, error) {
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

		return followingResponse, nil
	}, MaxRetries)
}
