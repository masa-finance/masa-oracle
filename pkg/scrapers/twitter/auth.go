package twitter

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func NewScraper(account *TwitterAccount, cookieDir string) *Scraper {
	scraper := &Scraper{Scraper: newTwitterScraper()}

	if err := LoadCookies(scraper.Scraper, account, cookieDir); err == nil {
		logrus.Debugf("Cookies loaded for user %s.", account.Username)
		if scraper.IsLoggedIn() {
			logrus.Debugf("Already logged in as %s.", account.Username)
			return scraper
		}
	}

	RandomSleep()

	if err := scraper.Login(account.Username, account.Password, account.TwoFACode); err != nil {
		logrus.WithError(err).Warnf("Login failed for %s", account.Username)
		return nil
	}

	RandomSleep()

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
