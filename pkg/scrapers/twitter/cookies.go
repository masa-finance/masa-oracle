package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
)

func SaveCookies(scraper *twitterscraper.Scraper, account *TwitterAccount, baseDir string) error {
	cookieFile := filepath.Join(baseDir, fmt.Sprintf("%s_twitter_cookies.json", account.Username))
	cookies := scraper.GetCookies()
	data, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("error marshaling cookies: %v", err)
	}
	if err = os.WriteFile(cookieFile, data, 0644); err != nil {
		return fmt.Errorf("error saving cookies: %v", err)
	}
	return nil
}

func LoadCookies(scraper *twitterscraper.Scraper, account *TwitterAccount, baseDir string) error {
	cookieFile := filepath.Join(baseDir, fmt.Sprintf("%s_twitter_cookies.json", account.Username))
	data, err := os.ReadFile(cookieFile)
	if err != nil {
		return fmt.Errorf("error reading cookies: %v", err)
	}
	var cookies []*http.Cookie
	if err = json.Unmarshal(data, &cookies); err != nil {
		return fmt.Errorf("error unmarshaling cookies: %v", err)
	}
	scraper.SetCookies(cookies)
	return nil
}
