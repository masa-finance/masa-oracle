// discordchannels.go
package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GuildChannel represents a Discord guild channel structure
type GuildChannel struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
}

// GetGuildChannels fetches the channels for a specific guild from the Discord API
func GetGuildChannels(guildID string) ([]GuildChannel, error) {
	botToken := os.Getenv("DISCORD_BOT_TOKEN") // Replace with your actual environment variable name
	if botToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", guildID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

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
