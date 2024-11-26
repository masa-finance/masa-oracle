package twitter

import (
	"fmt"

	twitterscraper "github.com/imperatrona/twitter-scraper"

	"github.com/sirupsen/logrus"
)

func ScrapeFollowersForProfile(baseDir string, username string, count int) ([]*twitterscraper.Profile, error) {
	scraper, account, err := getAuthenticatedScraper(baseDir)
	if err != nil {
		return nil, err
	}

	followingResponse, errString, _ := scraper.FetchFollowers(username, count, "")
	if errString != "" {
		err := fmt.Errorf("rate limited: %s", errString)
		if handleRateLimit(err, account) {
			return nil, err
		}

		logrus.Errorf("[-] Error fetching followers: %s", errString)
		return nil, fmt.Errorf("error fetching followers: %s", errString)
	}

	return followingResponse, nil
}
