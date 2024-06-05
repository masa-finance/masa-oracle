// discordchannels.go
package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GuildChannel represents a Discord guild channel structure
type GuildChannel struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
}

// GetGuildChannels fetches the channels for a specific guild from the Discord API
func GetGuildChannels(guildID, accessToken string) ([]GuildChannel, error) {
	url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", guildID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the authorization header to your OAuth2 access token
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the response body
		errorBody, err := io.ReadAll(resp.Body)
		if err != nil {
			// If we can't read the body, return the status code only
			return nil, fmt.Errorf("error fetching guild channels, status code: %d, error reading response body: %v", resp.StatusCode, err)
		}
		// Return the status code and the response body as a string
		return nil, fmt.Errorf("error fetching guild channels, status code: %d, response: %s", resp.StatusCode, string(errorBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var channels []GuildChannel
	if err := json.Unmarshal(body, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}
