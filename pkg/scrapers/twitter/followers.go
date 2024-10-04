package twitter

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

// ScrapeFollowersForProfile scrapes the followers of a specific Twitter user.
// It takes the username and count as parameters and returns the scraped followers information and an error if any.
func ScrapeFollowersForProfile(username string, count int) ([]twitterscraper.Legacy, error) {
	once.Do(initializeAccountManager)

	for {
		account := accountManager.GetNextAccount()
		if account == nil {
			return nil, fmt.Errorf("all accounts are rate-limited")
		}

		scraper := NewScraper(account)
		if scraper == nil {
			logrus.Errorf("Authentication failed for %s", account.Username)
			continue
		}

		followingResponse, errString, _ := scraper.FetchFollowers(username, count, "")
		if errString != "" {
			if strings.Contains(errString, "Rate limit exceeded") {
				accountManager.MarkAccountRateLimited(account)
				logrus.Warnf("Rate limited: %s", account.Username)
				continue
			}
			logrus.Errorf("Error fetching followers: %v", errString)
			return nil, fmt.Errorf("%v", errString)
		}

		// Marshal the followingResponse into a JSON string for logging
		responseJSON, err := json.Marshal(followingResponse)
		if err != nil {
			logrus.Errorf("Error marshaling followingResponse: %v", err)
		} else {
			logrus.Debugf("Following response: %s", responseJSON)
		}

		return followingResponse, nil
	}
}
