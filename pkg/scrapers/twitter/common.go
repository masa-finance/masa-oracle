package twitter

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

var (
	accountManager *TwitterAccountManager
	once           sync.Once
)

func initializeAccountManager() {
	accounts := loadAccountsFromConfig()
	accountManager = NewTwitterAccountManager(accounts)
}

func GetAccountManager() *TwitterAccountManager {
	_, _, err := getAuthenticatedScraper()
	if err != nil {
		logrus.Errorf("error initializing account manager: %v", err)
	}
	return accountManager
}

func loadAccountsFromConfig() []*TwitterAccount {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("error loading .env file: %v", err)
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

func getAuthenticatedScraper() (*Scraper, *TwitterAccount, *data_types.LoginEvent, error) {
	once.Do(initializeAccountManager)
	baseDir := config.GetInstance().MasaDir

	account := accountManager.GetNextAccount()
	if account == nil {
		return nil, nil, nil, fmt.Errorf("all accounts are rate-limited")
	}
	scraper, loginEvent := NewScraper(account, baseDir)
	if scraper == nil {
		err := fmt.Errorf("twitter authentication failed for %s", account.Username)
		logrus.Error(err)
		return nil, account, err
	}
	account.LastScraped = time.Now()
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
