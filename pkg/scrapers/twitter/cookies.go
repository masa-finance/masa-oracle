package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	twitterscraper "github.com/imperatrona/twitter-scraper"

	"github.com/sirupsen/logrus"
)

func SaveCookies(scraper *twitterscraper.Scraper, account *TwitterAccount, baseDir string) error {
	logrus.Debugf("Saving cookies for user %s", account.Username)
	cookieFile := filepath.Join(baseDir, fmt.Sprintf("%s_twitter_cookies.json", account.Username))
	cookies := scraper.GetCookies()
	logrus.Debugf("Got %d cookies to save", len(cookies))

	data, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("error marshaling cookies: %v", err)
	}

	logrus.Debugf("Writing cookies to file: %s", cookieFile)
	if err = os.WriteFile(cookieFile, data, 0644); err != nil {
		return fmt.Errorf("error saving cookies: %v", err)
	}
	logrus.Debug("Successfully saved cookies")
	return nil
}

func LoadCookies(scraper *twitterscraper.Scraper, account *TwitterAccount, baseDir string) error {
	logrus.Debugf("Loading cookies for user %s", account.Username)
	cookieFile := filepath.Join(baseDir, fmt.Sprintf("%s_twitter_cookies.json", account.Username))

	logrus.Debugf("Reading cookie file: %s", cookieFile)
	data, err := os.ReadFile(cookieFile)
	if err != nil {
		return fmt.Errorf("error reading cookies: %v", err)
	}

	var cookies []*http.Cookie
	if err = json.Unmarshal(data, &cookies); err != nil {
		return fmt.Errorf("error unmarshaling cookies: %v", err)
	}
	logrus.Debugf("Loaded %d cookies from file", len(cookies))

	// Verify critical cookies are present
	var hasAuthToken, hasCSRFToken bool
	for _, cookie := range cookies {
		if cookie.Name == "auth_token" {
			hasAuthToken = true
			logrus.Debug("Found auth_token cookie")
		}
		if cookie.Name == "ct0" {
			hasCSRFToken = true
			logrus.Debug("Found CSRF token cookie")
		}
	}

	if !hasAuthToken || !hasCSRFToken {
		logrus.Debug("Missing critical authentication cookies")
		return fmt.Errorf("missing critical authentication cookies")
	}

	logrus.Debug("Setting cookies in scraper")
	scraper.SetCookies(cookies)
	logrus.Debug("Successfully loaded and set cookies")
	return nil
}
