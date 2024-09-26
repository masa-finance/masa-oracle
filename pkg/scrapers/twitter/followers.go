package twitter

import (
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
	twitterscraper "github.com/masa-finance/twitter-scraper"
	"github.com/sirupsen/logrus"
)

// ScrapeFollowersForProfile scrapes the profile and tweets of a specific Twitter user.
// It takes the username as a parameter and returns the scraped profile information and an error if any.
func ScrapeFollowersForProfile(username string, count int) ([]*twitterscraper.Profile, error) {
	scraper := auth()

	if scraper == nil {
		return nil, fmt.Errorf("there was an error authenticating with your Twitter credentials")
	}

	followingResponse, errString, _ := scraper.FetchFollowers(username, count, "")
	if errString != "" {
		logrus.Printf("Error fetching profile: %v", errString)
		return nil, fmt.Errorf("%v", errString)
	}

	// Marshal the followingResponse into a JSON string for logging
	responseJSON, err := json.Marshal(followingResponse)
	if err != nil {
		// Log the error if the marshaling fails
		logrus.Errorf("[-] Error marshaling followingResponse: %v", err)
	} else {
		// Log the JSON string of followingResponse
		logrus.Debugf("Following response: %s", responseJSON)
	}

	return followingResponse, nil
}
