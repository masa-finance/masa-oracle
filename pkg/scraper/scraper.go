package scraper

import (
	"github.com/gocolly/colly/v2"
)

// Section represents a distinct part of a scraped webpage, typically defined by a heading.
// It contains a Title, representing the heading of the section, and Paragraphs, a slice of strings
// containing the text content found within that section.
type Section struct {
	Title      string   // Title is the heading text of the section.
	Paragraphs []string // Paragraphs contains all the text content of the section.
}

// CollectedData represents the aggregated result of the scraping process.
// It contains a slice of Section structs, each representing a distinct part of a scraped webpage.
type CollectedData struct {
	Sections []Section // Sections is a collection of webpage sections that have been scraped.
}

// Collect initiates the scraping process for the given list of URIs.
// It returns a CollectedData struct containing the scraped sections from each URI,
// and an error if any occurred during the scraping process.
func Collect(uri []string, depth int) (CollectedData, error) {
	var collectedData CollectedData

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.MaxDepth(depth),
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
