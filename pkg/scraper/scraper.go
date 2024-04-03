package scraper

import (
	"github.com/gocolly/colly/v2"
)

type Section struct {
	Title      string
	Paragraphs []string
}

type CollectedData struct {
	Sections []Section
}

// Collect a basic test
func Collect(uri []string) (CollectedData, error) {
	var collectedData CollectedData

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(100),
	)

	c.OnHTML("h1, h2", func(e *colly.HTMLElement) {
		// Directly append a new Section to collectedData.Sections
		collectedData.Sections = append(collectedData.Sections, Section{Title: e.Text})
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		// Check if there are any sections to append paragraphs to
		if len(collectedData.Sections) > 0 {
			// Get a reference to the last section
			lastSection := &collectedData.Sections[len(collectedData.Sections)-1]
			// Append the paragraph to the last section
			// Check for duplicate paragraphs before appending
			isDuplicate := false
			for _, paragraph := range lastSection.Paragraphs {
				if paragraph == e.Text {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				lastSection.Paragraphs = append(lastSection.Paragraphs, e.Text)
			}
		}
	})

	// OnRequest and OnError handlers remain the same

	// Visit each URL
	for _, u := range uri {
		err := c.Visit(u)
		if err != nil {
			return CollectedData{}, err
		}
	}

	// Wait for all requests to finish
	c.Wait()

	return collectedData, nil

}
