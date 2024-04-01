package scraper

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
)

// Collect a basic test
func Collect(uri string) (string, error) {

	fmt.Println(uri)

	c := colly.NewCollector(
		colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		err := c.Visit(e.Request.AbsoluteURL(link))
		if err != nil {
			logrus.Errorf("%v", err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	err := c.Visit("https://hackerspaces.org/")
	if err != nil {
		logrus.Errorf("%v", err)
	}
	return "", nil

}
