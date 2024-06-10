// discordguilds.go
package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Guild represents a Discord guild (server) structure
type Guild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int64  `json:"permissions"`
}

// GetUserGuilds fetches the guilds (servers) that the current user is part of
func GetUserGuilds() ([]Guild, error) {
	botToken := os.Getenv("DISCORD_BOT_TOKEN") // Replace with your actual environment variable name
	if botToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN environment variable not set")
	}

	url := "https://discord.com/api/users/@me/guilds"
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
		return nil, fmt.Errorf("error fetching guilds, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var guilds []Guild
	if err := json.Unmarshal(body, &guilds); err != nil {
		return nil, err
	}

	return guilds, nil
}
