package twitter

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

func NewScraper(account *TwitterAccount, cookieDir string) (*Scraper, *data_types.LoginEvent) {
	scraper := &Scraper{Scraper: newTwitterScraper()}
	var loginEvent *data_types.LoginEvent
	if err := LoadCookies(scraper.Scraper, account, cookieDir); err == nil {
		logrus.Debugf("Cookies loaded for user %s.", account.Username)
		if scraper.IsLoggedIn() {
			logrus.Debugf("Already logged in as %s.", account.Username)
			// Log a successful login event
			loginEvent = data_types.NewLoginEvent("", account.Username, "Twitter", true, "")
			return scraper, loginEvent
		}
	}

	RandomSleep()

	if err := scraper.Login(account); err != nil {
		logrus.WithError(err).Warnf("Login failed for %s", account.Username)
		// Log a failed login event
		loginEvent = data_types.NewLoginEvent("", account.Username, "Twitter", false, err.Error())
		return nil, loginEvent
	}

	RandomSleep()

	if err := SaveCookies(scraper.Scraper, account, cookieDir); err != nil {
		logrus.WithError(err).Errorf("Failed to save cookies for %s", account.Username)
	}

	logrus.Debugf("Login successful for %s", account.Username)
	// Log a successful login event
	loginEvent = data_types.NewLoginEvent("", account.Username, "Twitter", true, "")

	return scraper, loginEvent
}

func (scraper *Scraper) Login(account *TwitterAccount) error {
	var err error
	if len(account.TwoFACode) > 0 {
		err = scraper.Scraper.Login(account.Username, account.Password, account.TwoFACode)
	} else {
		err = scraper.Scraper.Login(account.Username, account.Password)
	}
	if err != nil {
		account.LoginStatus = fmt.Sprintf("Failed - %v", err)
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
