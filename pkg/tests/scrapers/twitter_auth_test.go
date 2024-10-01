package scrapers_test

import (
	"os"
	"path/filepath"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Twitter Auth Function", func() {
	var (
		twitterUsername string
		twitterPassword string
		twoFACode       string
	)

	BeforeEach(func() {
		tempDir := GinkgoT().TempDir()
		config.GetInstance().MasaDir = tempDir

		twitterUsername = os.Getenv("TWITTER_USERNAME")
		twitterPassword = os.Getenv("TWITTER_PASSWORD")
		twoFACode = os.Getenv("TWITTER_2FA_CODE")

		Expect(twitterUsername).NotTo(BeEmpty(), "TWITTER_USERNAME environment variable is not set")
		Expect(twitterPassword).NotTo(BeEmpty(), "TWITTER_PASSWORD environment variable is not set")

		config.GetInstance().TwitterUsername = twitterUsername
		config.GetInstance().TwitterPassword = twitterPassword
		config.GetInstance().Twitter2FaCode = twoFACode
	})

	authenticate := func() *twitterscraper.Scraper {
		return twitter.Auth()
	}

	checkLoggedIn := func(scraper *twitterscraper.Scraper) bool {
		return twitter.IsLoggedIn(scraper)
	}

	It("authenticates and logs in successfully", func() {
		scraper := authenticate()
		Expect(scraper).NotTo(BeNil())

		cookieFile := filepath.Join(config.GetInstance().MasaDir, "twitter_cookies.json")
		Expect(cookieFile).To(BeAnExistingFile())

		Expect(checkLoggedIn(scraper)).To(BeTrue())
		logrus.Info("Authenticated and logged in to Twitter")
	})

	It("reuses session from cookies", func() {
		firstScraper := authenticate()
		Expect(firstScraper).NotTo(BeNil())

		secondScraper := authenticate()
		Expect(secondScraper).NotTo(BeNil())

		Expect(checkLoggedIn(secondScraper)).To(BeTrue())
		logrus.Info("Reused session from cookies")
	})

	AfterEach(func() {
		os.RemoveAll(config.GetInstance().MasaDir)
	})
})
