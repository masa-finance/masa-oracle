package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
)

func SaveCookies(scraper *twitterscraper.Scraper, filePath string) error {
	cookies := scraper.GetCookies()
	js, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("error marshaling cookies: %v", err)
	}
	err = os.WriteFile(filePath, js, 0644)
	if err != nil {
		return fmt.Errorf("error saving cookies to file: %v", err)
	}

	// Load the saved cookies back into the scraper
	if err := LoadCookies(scraper, filePath); err != nil {
		return fmt.Errorf("error loading saved cookies: %v", err)
	}

	return nil
}

func LoadCookies(scraper *twitterscraper.Scraper, filePath string) error {
	js, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading cookies from file: %v", err)
	}
	var cookies []*http.Cookie
	err = json.Unmarshal(js, &cookies)
	if err != nil {
		return fmt.Errorf("error unmarshaling cookies: %v", err)
	}
	scraper.SetCookies(cookies)
	return nil
}
