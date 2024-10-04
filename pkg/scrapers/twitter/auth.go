package twitter

import (
	"fmt"
	"sync"
	"time"

	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

type Scraper struct {
	*twitterscraper.Scraper
}

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
	account.RateLimitedUntil = time.Now().Add(GetRateLimitDuration())
}

func NewScraper(account *TwitterAccount, cookieDir string) *Scraper {
	ts := twitterscraper.New()
	scraper := &Scraper{Scraper: ts}

	if err := LoadCookies(scraper.Scraper, account, cookieDir); err == nil {
		logrus.Debugf("Cookies loaded for user %s.", account.Username)
		if scraper.IsLoggedIn() {
			logrus.Debugf("Already logged in as %s.", account.Username)
			return scraper
		}
	}

	ShortSleep()

	if err := scraper.Login(account.Username, account.Password, account.TwoFACode); err != nil {
		logrus.WithError(err).Warnf("Login failed for %s", account.Username)
		return nil
	}

	ShortSleep()

	if err := SaveCookies(scraper.Scraper, account, cookieDir); err != nil {
		logrus.WithError(err).Errorf("Failed to save cookies for %s", account.Username)
	}

	logrus.Debugf("Login successful for %s", account.Username)
	return scraper
}

func (scraper *Scraper) Login(username, password string, twoFACode ...string) error {
	var err error
	if len(twoFACode) > 0 {
		err = scraper.Scraper.Login(username, password, twoFACode[0])
	} else {
		err = scraper.Scraper.Login(username, password)
	}
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}
	return nil
}

func (scraper *Scraper) Logout() error {
	if err := scraper.Scraper.Logout(); err != nil {
		return fmt.Errorf("logout failed: %v", err)
	}
	return nil
}

func (scraper *Scraper) IsLoggedIn() bool {
	return scraper.Scraper.IsLoggedIn()
}
