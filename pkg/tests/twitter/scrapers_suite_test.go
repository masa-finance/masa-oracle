package scrapers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestScrapers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scrapers Suite")
}

func TestMain(m *testing.M) {
	// Override os.Args to prevent flag parsing errors
	os.Args = []string{os.Args[0]}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("Failed to get current directory: %v", err)
	}

	// Define the project root path
	projectRoot := filepath.Join(cwd, "..", "..", "..")
	envPath := filepath.Join(projectRoot, ".env")

	// Load the .env file
	err = godotenv.Load(envPath)
	if err != nil {
		logrus.Warnf("Error loading .env file from %s: %v", envPath, err)
	} else {
		logrus.Info("Loaded .env from project root")
	}

	// Verify that the required environment variables are set
	requiredEnvVars := []string{"USER_AGENTS", "TWITTER_USERNAME", "TWITTER_PASSWORD"}
	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			logrus.Warnf("%s environment variable is not set", envVar)
		} else {
			logrus.Debugf("%s: %s", envVar, value)
		}
	}

	// Run the tests
	exitCode := m.Run()
	os.Exit(exitCode)
}
