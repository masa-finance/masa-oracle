package twitter

import (
	"fmt"
	"sync"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/config"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

type TwitterAccount struct {
	Username         string
	Password         string
	TwoFACode        string
	RateLimitedUntil time.Time
}

type TwitterAccountManager struct {
	accounts []*TwitterAccount
	index    int
	mutex    sync.Mutex
}

func NewTwitterAccountManager(accounts []*TwitterAccount) *TwitterAccountManager {
	return &TwitterAccountManager{
		accounts: accounts,
		index:    0,
	}
}

func (manager *TwitterAccountManager) GetNextAccount() *TwitterAccount {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	for i := 0; i < len(manager.accounts); i++ {
		account := manager.accounts[manager.index]
		manager.index = (manager.index + 1) % len(manager.accounts)
		if time.Now().After(account.RateLimitedUntil) {
			return account
		}
	}
	return nil
}

func (manager *TwitterAccountManager) MarkAccountRateLimited(account *TwitterAccount) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	account.RateLimitedUntil = time.Now().Add(time.Hour)
}

func Auth(account *TwitterAccount) *twitterscraper.Scraper {
	scraper := twitterscraper.New()
	baseDir := config.GetInstance().MasaDir

	if err := LoadCookies(scraper, account, baseDir); err == nil {
		logrus.Debugf("Cookies loaded for user %s.", account.Username)
		if IsLoggedIn(scraper) {
			logrus.Debugf("Already logged in as %s.", account.Username)
			return scraper
		}
	}

	time.Sleep(100 * time.Millisecond)

	var err error
	if account.TwoFACode != "" {
		err = Login(scraper, account.Username, account.Password, account.TwoFACode)
	} else {
		err = Login(scraper, account.Username, account.Password)
	}

	if err != nil {
		logrus.WithError(err).Warnf("Login failed for %s", account.Username)
		return nil
	}

	time.Sleep(100 * time.Millisecond)

	if err = SaveCookies(scraper, account, baseDir); err != nil {
		logrus.WithError(err).Errorf("Failed to save cookies for %s", account.Username)
	}

	logrus.Debugf("Login successful for %s", account.Username)
	return scraper
}

func Login(scraper *twitterscraper.Scraper, credentials ...string) error {
	var err error
	switch len(credentials) {
	case 2:
		err = scraper.Login(credentials[0], credentials[1])
	case 3:
		err = scraper.Login(credentials[0], credentials[1], credentials[2])
	default:
		return fmt.Errorf("invalid number of credentials")
	}
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func IsLoggedIn(scraper *twitterscraper.Scraper) bool {
	return scraper.IsLoggedIn()
}

func Logout(scraper *twitterscraper.Scraper) error {
	if err := scraper.Logout(); err != nil {
		return fmt.Errorf("logout failed: %v", err)
	}
	return nil
}
