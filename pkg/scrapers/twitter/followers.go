package twitter

import (
	"fmt"
	"time"

	twitterscraper "github.com/imperatrona/twitter-scraper"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"

	"github.com/sirupsen/logrus"
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
	account.LastScraped = time.Now()
	return followingResponse, nil
}
