package twitter

import (
	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type Scraper struct {
	*twitterscraper.Scraper
}

func newTwitterScraper() *twitterscraper.Scraper {
	return twitterscraper.New()
}

func (scraper *Scraper) IsLoggedIn() bool {
	return scraper.Scraper.IsLoggedIn()
}
