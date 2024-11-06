// TODO: This is a WIP
// Package scrapers_test contains integration tests for the Twitter scraper functionality.
//
// Dev Notes:
// - These tests require valid Twitter credentials set in environment variables.
// - The tests use a temporary directory for storing cookies and other data.
// - Make sure to run these tests in a controlled environment to avoid rate limiting.
// - The tests cover authentication, session reuse, and basic scraping operations.
package scrapers_test

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
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

	loadEnv := func() {
		_, filename, _, _ := runtime.Caller(0)
		projectRoot := filepath.Join(filepath.Dir(filename), "..", "..", "..")
		envPath := filepath.Join(projectRoot, ".env")

		err := godotenv.Load(envPath)
		if err != nil {
			logrus.Warnf("Error loading .env file from %s: %v", envPath, err)
		} else {
			logrus.Infof("Loaded .env from %s", envPath)
		}
	}

	BeforeEach(func() {
		loadEnv()

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
		return nil
		//return twitter.Auth()
	}

	PIt("authenticates and logs in successfully", func() {
		// Ensure cookie file doesn't exist before authentication
		cookieFile := filepath.Join(config.GetInstance().MasaDir, "twitter_cookies.json")
		Expect(cookieFile).NotTo(BeAnExistingFile())

		// Authenticate
		scraper := authenticate()
		Expect(scraper).NotTo(BeNil())

		// Check if cookie file was created
		Expect(cookieFile).To(BeAnExistingFile())

		// Verify logged in state
		Expect(scraper.IsLoggedIn()).To(BeTrue())

		// Attempt a simple operation to verify the session is valid
		profile, err := twitter.ScrapeTweetsProfile("twitter")
		Expect(err).To(BeNil())
		Expect(profile.Username).To(Equal("twitter"))

		logrus.Info("Authenticated and logged in to Twitter successfully")
	})

	PIt("reuses session from cookies", func() {
		// First authentication
		firstScraper := authenticate()
		Expect(firstScraper).NotTo(BeNil())

		// Verify cookie file is created
		cookieFile := filepath.Join(config.GetInstance().MasaDir, "twitter_cookies.json")
		Expect(cookieFile).To(BeAnExistingFile())

		// Clear the scraper to force cookie reuse
		firstScraper = nil // nolint: ineffassign

		// Second authentication (should use cookies)
		secondScraper := authenticate()
		Expect(secondScraper).NotTo(BeNil())

		// Verify logged in state
		Expect(secondScraper.IsLoggedIn()).To(BeTrue())

		// Attempt a simple operation to verify the session is valid
		profile, err := twitter.ScrapeTweetsProfile("twitter")
		Expect(err).To(BeNil())
		Expect(profile.Username).To(Equal("twitter"))

		logrus.Info("Reused session from cookies successfully")
	})

	PIt("scrapes the profile of 'god' and recent #Bitcoin tweets using saved cookies", func() {
		// First authentication
		firstScraper := authenticate()
		Expect(firstScraper).NotTo(BeNil())

		// Verify cookie file is created
		cookieFile := filepath.Join(config.GetInstance().MasaDir, "twitter_cookies.json")
		Expect(cookieFile).To(BeAnExistingFile())

		// Clear the scraper to force cookie reuse
		firstScraper = nil // nolint: ineffassign

		// Second authentication (should use cookies)
		secondScraper := authenticate()
		Expect(secondScraper).NotTo(BeNil())

		// Verify logged in state
		Expect(secondScraper.IsLoggedIn()).To(BeTrue())

		// Attempt to scrape profile
		profile, err := twitter.ScrapeTweetsProfile("god")
		Expect(err).To(BeNil())
		logrus.Infof("Profile of 'god': %+v", profile)

		// Scrape recent #Bitcoin tweets
		tweets, err := twitter.ScrapeTweetsByQuery("#Bitcoin", 3)
		Expect(err).To(BeNil())
		Expect(tweets).To(HaveLen(3))

		logrus.Info("Recent #Bitcoin tweets:")
		for i, tweet := range tweets {
			logrus.Infof("Tweet %d: %s", i+1, tweet.Tweet.Text)
		}
	})

	AfterEach(func() {
		os.RemoveAll(config.GetInstance().MasaDir)
	})
})
