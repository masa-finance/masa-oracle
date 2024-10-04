package twitter

import (
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
)

type Scraper struct {
	*twitterscraper.Scraper
}

func newTwitterScraper() *twitterscraper.Scraper {
	return twitterscraper.New()
}

func (s *Scraper) IsLoggedIn() bool {
	return s.Scraper.IsLoggedIn()
}
