// discordprofile.go
package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UserProfile holds the structure for a Discord user profile response
type UserProfile struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

// GetUserProfile fetches a user's profile from Discord API
func GetUserProfile(userID, botToken string) (*UserProfile, error) {
	url := fmt.Sprintf("https://discord.com/api/users/%s", userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the authorization header to your bot token
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching user profile, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var profile UserProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
