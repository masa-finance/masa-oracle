package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gocolly/colly/v2"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	"github.com/sirupsen/logrus"
)

// Section represents a distinct part of a scraped webpage, typically defined by a heading.
// It contains a Title, representing the heading of the section, and Paragraphs, a slice of strings
// containing the text content found within that section.
type Section struct {
	Title      string   `json:"title"`      // Title is the heading text of the section.
	Paragraphs []string `json:"paragraphs"` // Paragraphs contains all the text content of the section.
	Images     []string `json:"images"`     // Images storing base64 - maybe!!?
}

// CollectedData represents the aggregated result of the scraping process.
// It contains a slice of Section structs, each representing a distinct part of a scraped webpage.
type CollectedData struct {
	Sections []Section `json:"sections"` // Sections is a collection of webpage sections that have been scraped.
	Pages    []string  `json:"pages"`
}

// ScrapeWebDataForSentiment initiates the scraping process for the given list of URIs.
// It returns a CollectedData struct containing the scraped sections from each URI,
// and an error if any occurred during the scraping process.
// Usage:
// @param uri: string - url to scrape
// @param depth: int - depth of how many subpages to scrape
// @param model: string - model to use for sentiment analysis
// Example:
//
//	go func() {
//		res, err := scraper.ScrapeWebDataForSentiment([]string{"https://en.wikipedia.org/wiki/Maize"}, 5)
//		if err != nil {
//			logrus.Errorf("[-] Error collecting data: %s", err.Error())
//		return
//	  }
//	logrus.Infof("%+v", res)
//	}()
func ScrapeWebDataForSentiment(uri []string, depth int, model string) (string, string, error) {
	var collectedData CollectedData

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.IgnoreRobotsTxt(),
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
			// Handle dupes
			if !isDuplicate {
				lastSection.Paragraphs = append(lastSection.Paragraphs, e.Text)
			}
		}
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		imageURL := e.Request.AbsoluteURL(e.Attr("src"))
		if len(collectedData.Sections) > 0 {
			lastSection := &collectedData.Sections[len(collectedData.Sections)-1]
			lastSection.Images = append(lastSection.Images, imageURL)
		}

		// @todo implement get image search for text and retrieve and add to data struct
		// // Fetch the image content
		// resp, err := http.Get(imageURL)
		// if err != nil {
		// 	// Handle error
		// 	return
		// }
		// defer resp.Body.Close()
		// imgBytes, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	// Handle error
		// 	return
		// }
		// // Convert image content to base64
		// base64Image := base64.StdEncoding.EncodeToString(imgBytes)
		// // Append this base64Image to your structured data
		// if len(collectedData.Sections) > 0 {
		// 	lastSection := &collectedData.Sections[len(collectedData.Sections)-1]
		// 	// Assuming Section struct has an Images field which is a slice of strings
		// 	lastSection.Images = append(lastSection.Images, base64Image)
		// }
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		pageURL := e.Request.AbsoluteURL(e.Attr("href"))
		collectedData.Pages = append(collectedData.Pages, pageURL)
		_ = e.Request.Visit(pageURL)
	})

	// Visit each URL
	for _, u := range uri {
		err := c.Visit(u)
		if err != nil {
			return "", "", err
		}
	}

	// Wait for all requests to finish
	c.Wait()

	j, _ := json.Marshal(collectedData)
	sentimentPrompt := "Please perform a sentiment analysis on the following text, using an unbiased approach. Sentiment analysis involves identifying and categorizing opinions expressed in text, particularly to determine whether the writer's attitude towards a particular topic, product, etc., is positive, negative, or neutral. After analyzing, please provide a summary of the overall sentiment expressed in this text, including the proportion of positive, negative, and neutral sentiments if applicable."
	sentimentRequest, sentimentSummary, err := llmbridge.AnalyzeSentimentWeb(string(j), model, sentimentPrompt)

	if err != nil {
		return "", "", err
	}
	return sentimentRequest, sentimentSummary, nil
}

// ScrapeWebData initiates the scraping process for the given list of URIs.
// It returns a CollectedData struct containing the scraped sections from each URI,
// and an error if any occurred during the scraping process.
// Usage:
// @param uri: string - url to scrape
// @param depth: int - depth of how many subpages to scrape
// Example:
//
//	go func() {
//		res, err := scraper.ScrapeWebDataForSentiment([]string{"https://en.wikipedia.org/wiki/Maize"}, 5)
//		if err != nil {
//			logrus.Errorf("[-] Error collecting data: %s", err.Error())
//		return
//	  }
//	logrus.Infof("%+v", res)
//	}()
func ScrapeWebData(uri []string, depth int) ([]byte, error) {
	// Set default depth to 1 if 0 is provided
	if depth <= 0 {
		depth = 1
	}

	var collectedData CollectedData

	c := colly.NewCollector(
		colly.Async(true), // Enable asynchronous requests
		colly.AllowURLRevisit(),
		colly.IgnoreRobotsTxt(),
		colly.MaxDepth(depth),
	)

	// Adjust the parallelism and delay based on your needs and server capacity
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,                      // Increased parallelism
		Delay:       500 * time.Millisecond, // Reduced delay
	})

	// Increase the timeout slightly if necessary
	c.SetRequestTimeout(240 * time.Second) // Increased to 4 minutes

	// Initialize a backoff strategy
	backoffStrategy := backoff.NewExponentialBackOff()

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == http.StatusTooManyRequests {
			// Parse the Retry-After header (in seconds)
			retryAfter, convErr := strconv.Atoi(r.Headers.Get("Retry-After"))
			if convErr != nil {
				// If not in seconds, it might be a date. Handle accordingly.
				logrus.Debugf("[-] Retry-After: %s", r.Headers.Get("Retry-After"))
			}
			// Calculate the next delay
			nextDelay := backoffStrategy.NextBackOff()
			if retryAfter > 0 {
				nextDelay = time.Duration(retryAfter) * time.Second
			}
			logrus.Warnf("[-] Rate limited. Retrying after %v", nextDelay)
			time.Sleep(nextDelay)
			// Retry the request
			_ = r.Request.Retry()
		} else {
			logrus.Errorf("[-] Request URL: %s failed with error: %v", r.Request.URL, err)
		}
	})
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
			// Handle dupes
			if !isDuplicate {
				lastSection.Paragraphs = append(lastSection.Paragraphs, e.Text)
			}
		}
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
		imageURL := e.Request.AbsoluteURL(e.Attr("src"))
		if len(collectedData.Sections) > 0 {
			lastSection := &collectedData.Sections[len(collectedData.Sections)-1]
			lastSection.Images = append(lastSection.Images, imageURL)
		}
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		pageURL := e.Request.AbsoluteURL(e.Attr("href"))
		// Check if the URL protocol is supported (http or https)
		if strings.HasPrefix(pageURL, "http://") || strings.HasPrefix(pageURL, "https://") {
			collectedData.Pages = append(collectedData.Pages, pageURL)
			_ = e.Request.Visit(pageURL)
		} else {
			logrus.Warnf("Unsupported URL protocol, skipping: %s", pageURL)
		}
	})

	for _, u := range uri {
		err := c.Visit(u)
		if err != nil {
			return nil, err
		}
	}

	// Wait for all requests to finish
	c.Wait()

	j, _ := json.Marshal(collectedData)
	return j, nil
}
