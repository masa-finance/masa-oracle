package twitter

import (
	twitterscraper "github.com/imperatrona/twitter-scraper"

	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

func ScrapeTweetsProfile(username string) (twitterscraper.Profile, *data_types.LoginEvent, error) {
	scraper, account, loginEvent, err := getAuthenticatedScraper()
	if err != nil {
		return twitterscraper.Profile{}, loginEvent, err
	}

	profile, err := scraper.GetProfile(username)
	if err != nil {
		if handleRateLimit(err, account) {
			return twitterscraper.Profile{}, loginEvent, err
		}
		return twitterscraper.Profile{}, loginEvent, err
	}
	return profile, loginEvent, nil
}
