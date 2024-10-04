package twitter

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/masa-finance/masa-oracle/pkg/config"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

// ScrapeFollowersForProfile scrapes the followers of a specific Twitter user.
// It takes the username and count as parameters and returns the scraped followers information and an error if any.
func ScrapeFollowersForProfile(username string, count int) ([]twitterscraper.Legacy, error) {
	once.Do(initializeAccountManager)
	baseDir := config.GetInstance().MasaDir

	for {
		account := accountManager.GetNextAccount()
		if account == nil {
			return nil, fmt.Errorf("all accounts are rate-limited")
		}

		scraper := NewScraper(account, baseDir)
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

		return followingResponse, nil
	}
}
