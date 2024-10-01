package scrapers_test

import (
	"os"
	"path/filepath"

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
		err := godotenv.Load()
		if err != nil {
			logrus.Warn("Error loading .env file")
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
		return twitter.Auth()
	}

	checkLoggedIn := func(scraper *twitterscraper.Scraper) bool {
		return twitter.IsLoggedIn(scraper)
	}

	It("authenticates and logs in successfully", func() {
		// Ensure cookie file doesn't exist before authentication
		cookieFile := filepath.Join(config.GetInstance().MasaDir, "twitter_cookies.json")
		Expect(cookieFile).NotTo(BeAnExistingFile())

		// Authenticate
		scraper := authenticate()
		Expect(scraper).NotTo(BeNil())

		// Check if cookie file was created
		Expect(cookieFile).To(BeAnExistingFile())

		// Verify logged in state
		Expect(checkLoggedIn(scraper)).To(BeTrue())

		// Attempt a simple operation to verify the session is valid
		profile, err := twitter.ScrapeTweetsProfile("twitter")
		Expect(err).To(BeNil())
		Expect(profile.Username).To(Equal("twitter"))

		logrus.Info("Authenticated and logged in to Twitter successfully")
	})

	It("reuses session from cookies", func() {
		// First authentication
		firstScraper := authenticate()
		Expect(firstScraper).NotTo(BeNil())

		// Verify cookie file is created
		cookieFile := filepath.Join(config.GetInstance().MasaDir, "twitter_cookies.json")
		Expect(cookieFile).To(BeAnExistingFile())

		// Clear the scraper to force cookie reuse
		firstScraper = nil

		// Second authentication (should use cookies)
		secondScraper := authenticate()
		Expect(secondScraper).NotTo(BeNil())

		// Verify logged in state
		Expect(checkLoggedIn(secondScraper)).To(BeTrue())

		// Attempt a simple operation to verify the session is valid
		profile, err := twitter.ScrapeTweetsProfile("twitter")
		Expect(err).To(BeNil())
		Expect(profile.Username).To(Equal("twitter"))

		logrus.Info("Reused session from cookies successfully")
	})

	It("scrapes the profile of 'god' and recent #Bitcoin tweets using saved cookies", func() {
		// First authentication
		firstScraper := authenticate()
		Expect(firstScraper).NotTo(BeNil())

		// Verify cookie file is created
		cookieFile := filepath.Join(config.GetInstance().MasaDir, "twitter_cookies.json")
		Expect(cookieFile).To(BeAnExistingFile())

		// Clear the scraper to force cookie reuse
		firstScraper = nil

		// Second authentication (should use cookies)
		secondScraper := authenticate()
		Expect(secondScraper).NotTo(BeNil())

		// Verify logged in state
		Expect(twitter.IsLoggedIn(secondScraper)).To(BeTrue())

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
