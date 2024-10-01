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
			depth := 2

			data, err := web.ScrapeWebData(urls, depth)
			Expect(err).ToNot(HaveOccurred())

			var result web.CollectedData
			Expect(json.Unmarshal(data, &result)).To(Succeed())
			Expect(result.Sections).ToNot(BeEmpty())

			logrus.Infof("Scraped Data: %+v", result)
		})
	})
})
