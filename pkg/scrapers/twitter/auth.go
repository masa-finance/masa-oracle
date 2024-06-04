package twitter

import (
	"fmt"

	twitterscraper "github.com/masa-finance/masa-twitter-scraper"
)

// Login attempts to log in to the Twitter scraper service.
// It supports three modes of operation:
// 1. Basic login using just a username and password.
// 2. Login requiring an email confirmation, using a username, password, and email address.
// 3. Login with two-factor authentication, using a username, password, and 2FA code.
// Parameters:
//   - scraper: A pointer to an instance of the twitterscraper.Scraper.
//   - credentials: A variadic list of strings representing login credentials.
//     The function expects either two strings (username, password) for basic login,
//     or three strings (username, password, email/2FA code) for email confirmation or 2FA.
//
// Returns an error if login fails or if an invalid number of credentials is provided.
func Login(scraper *twitterscraper.Scraper, credentials ...string) error {
	var err error
	switch len(credentials) {
	case 2:
		// Basic login with username and password.
		err = scraper.Login(credentials[0], credentials[1])
	case 3:
		// The third parameter is used for either email confirmation or a 2FA code.
		// This design assumes the Twitter scraper's Login method can contextually handle both cases.
		err = scraper.Login(credentials[0], credentials[1], credentials[2])
	default:
		// Return an error if the number of provided credentials is neither 2 nor 3.
		return fmt.Errorf("invalid number of login credentials provided")
	}
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func IsLoggedIn(scraper *twitterscraper.Scraper) bool {
	return scraper.IsLoggedIn()
}

func Logout(scraper *twitterscraper.Scraper) error {
	err := scraper.Logout()
	if err != nil {
		return fmt.Errorf("logout failed: %v", err)
	}
	return nil
}
