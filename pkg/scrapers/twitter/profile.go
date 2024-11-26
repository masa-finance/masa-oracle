package twitter

import (
	twitterscraper "github.com/imperatrona/twitter-scraper"
)

func ScrapeTweetsProfile(username string) (twitterscraper.Profile, error) {
	scraper, account, err := getAuthenticatedScraper()
	if err != nil {
		return twitterscraper.Profile{}, err
	}

	profile, err := scraper.GetProfile(username)
	if err != nil {
		if handleRateLimit(err, account) {
			return twitterscraper.Profile{}, err
		}
		return twitterscraper.Profile{}, err
	}
	return profile, nil
}
