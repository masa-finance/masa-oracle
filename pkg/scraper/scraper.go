package scraper

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
)

// Collect a basic test
func Collect(uri []string) (string, error) {
	var collectedTexts []string

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(100),
	)

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		text := e.Text
		collectedTexts = append(collectedTexts, text)
	})

	c.OnHTML("h2", func(e *colly.HTMLElement) {
		text := e.Text
		collectedTexts = append(collectedTexts, text)
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		text := e.Text
		collectedTexts = append(collectedTexts, text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		logrus.Errorf("OnError: %v", err)
	})

	// Visit each URL in the slice
	for _, u := range uri {
		err := c.Visit(u)
		if err != nil {
			return "", err
		}
	}

	// Combine all collected texts into a single string
	result := strings.Join(collectedTexts, "\n")

	return result, nil

}
