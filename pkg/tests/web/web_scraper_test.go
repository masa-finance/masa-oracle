// Package scrapers_test contains tests for web scraping functionality.
//
// Dev Notes:
// - This test suite uses Ginkgo and Gomega for BDD-style testing.
// - It tests the ScrapeWebData function from the web package.
// - The test currently scrapes data from a single URL (CoinMarketCap) with a depth of 1.
// - TODO: Consider adding more diverse test cases with multiple URLs and depths.
// - TODO: Add more specific assertions on the scraped data structure and content.
package scrapers_test

import (
	"encoding/json"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Web Scraper", func() {
	Context("ScrapeWebData", func() {
		It("scrapes data from URLs", func() {
			urls := []string{"https://coinmarketcap.com/currencies/masa-network/"}
			depth := 1

			data, err := web.ScrapeWebData(urls, depth)
			Expect(err).ToNot(HaveOccurred())

			var result web.CollectedData
			Expect(json.Unmarshal(data, &result)).To(Succeed())
			Expect(result.Sections).ToNot(BeEmpty())

			logrus.Infof("Scraped Data: %+v", result)
		})
	})
})
