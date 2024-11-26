package twitter

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	accountManager *TwitterAccountManager
	once           sync.Once
)

func initializeAccountManager() {
	accounts := loadAccountsFromConfig()
	accountManager = NewTwitterAccountManager(accounts)
}

func loadAccountsFromConfig() []*TwitterAccount {
	err := godotenv.Load()
	if err != nil {
		logrus.Errorf("error loading .env file: %v", err)
	}

	accountsEnv := os.Getenv("TWITTER_ACCOUNTS")
	if accountsEnv == "" {
		logrus.Fatal("TWITTER_ACCOUNTS not set in .env file")
	}

	return parseAccounts(strings.Split(accountsEnv, ","))
}

func parseAccounts(accountPairs []string) []*TwitterAccount {
	return filterMap(accountPairs, func(pair string) (*TwitterAccount, bool) {
		credentials := strings.Split(pair, ":")
		if len(credentials) != 2 {
			logrus.Warnf("invalid account credentials: %s", pair)
			return nil, false
		}
		return &TwitterAccount{
			Username: strings.TrimSpace(credentials[0]),
			Password: strings.TrimSpace(credentials[1]),
		}, true
	})
}

func getAuthenticatedScraper(baseDir string) (*Scraper, *TwitterAccount, error) {
	once.Do(initializeAccountManager)

	account := accountManager.GetNextAccount()
	if account == nil {
		return nil, nil, fmt.Errorf("all accounts are rate-limited")
	}
	scraper := NewScraper(account, baseDir)
	if scraper == nil {
		logrus.Errorf("Authentication failed for %s", account.Username)
		return nil, account, fmt.Errorf("Twitter authentication failed for %s", account.Username)
	}
	return scraper, account, nil
}

func handleRateLimit(err error, account *TwitterAccount) bool {
	if strings.Contains(err.Error(), "Rate limit exceeded") {
		accountManager.MarkAccountRateLimited(account)
		logrus.Warnf("rate limited: %s", account.Username)
		return true
	}
	return false
}

func filterMap[T any, R any](slice []T, f func(T) (R, bool)) []R {
	result := make([]R, 0, len(slice))
	for _, v := range slice {
		if r, ok := f(v); ok {
			result = append(result, r)
		}
	}
	return result
}
