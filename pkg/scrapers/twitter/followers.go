package twitter

import (
	"fmt"

	twitterscraper "github.com/imperatrona/twitter-scraper"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

func ScrapeFollowersForProfile(username string, count int) ([]*twitterscraper.Profile, *data_types.LoginEvent, error) {
	scraper, account, loginEvent, err := getAuthenticatedScraper()
	if err != nil {
		return nil, loginEvent, err
	}

	followingResponse, errString, _ := scraper.FetchFollowers(username, count, "")
	if errString != "" {
		if handleRateLimit(fmt.Errorf(errString), account) {
			return nil, loginEvent, fmt.Errorf("rate limited")
		}
		logrus.Errorf("Error fetching followers: %v", errString)
		return nil, loginEvent, fmt.Errorf("%v", errString)
	}

	return followingResponse, loginEvent, nil
}
