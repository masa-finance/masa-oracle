package twitterscraper_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
	"github.com/sirupsen/logrus"
)

var (
	username     string
	password     string
	email        string
	skipAuthTest bool
	testScraper  *twitterscraper.Scraper
	logger       *logrus.Logger
)

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Warnf("Error loading .env file: %v", err)
	}

	username = os.Getenv("TWITTER_USERNAME")
	logger.Infof("Loaded TWITTER_USERNAME: %s", username)

	password = os.Getenv("TWITTER_PASSWORD")
	logger.Infof("Loaded TWITTER_PASSWORD: %s", maskPassword(password))

	email = os.Getenv("TWITTER_EMAIL")
	logger.Infof("Loaded TWITTER_EMAIL: %s", email)

	skipAuthTest = os.Getenv("SKIP_AUTH_TEST") != ""
	logger.Infof("Loaded SKIP_AUTH_TEST: %v", skipAuthTest)

	testScraper = twitterscraper.New()

	if username != "" && password != "" && !skipAuthTest {
		err := testScraper.Login(username, password, email)
		if err != nil {
			logger.Warnf("Login() error = %v", err)
		} else {
			logger.Info("Successfully logged in with test scraper")
		}
	}
}

func TestAuth(t *testing.T) {
	if skipAuthTest {
		t.Skip("Skipping test due to environment variable")
	}
	scraper := twitterscraper.New()
	if err := scraper.Login(username, password, email); err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if !scraper.IsLoggedIn() {
		t.Fatalf("Expected IsLoggedIn() = true")
	}
	cookies := scraper.GetCookies()
	scraper2 := twitterscraper.New()
	scraper2.SetCookies(cookies)
	if !scraper2.IsLoggedIn() {
		t.Error("Expected restored IsLoggedIn() = true")
	}
	if err := scraper.Logout(); err != nil {
		t.Errorf("Logout() error = %v", err)
	}
	if scraper.IsLoggedIn() {
		t.Error("Expected IsLoggedIn() = false")
	}
}

func TestLoginOpenAccount(t *testing.T) {
	scraper := twitterscraper.New()
	if err := scraper.LoginOpenAccount(); err != nil {
		t.Fatalf("LoginOpenAccount() error = %v", err)
	}
}

func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
