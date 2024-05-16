package discord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GuildMember represents a member within a guild, including their roles and statuses
type GuildMember struct {
	User         UserProfile `json:"user"` // Reuse UserProfile from discordprofile.go
	Nick         string      `json:"nick,omitempty"`
	Roles        []string    `json:"roles"`
	JoinedAt     string      `json:"joined_at"`
	PremiumSince string      `json:"premium_since,omitempty"`
	Deaf         bool        `json:"deaf"`
	Mute         bool        `json:"mute"`
}

// ListGuildMemberships lists the guild memberships for a given user ID
func ListGuildMemberships(userID, botToken string) ([]GuildMember, error) {
	// This endpoint and method are hypothetical and do not directly correspond to Discord's API.
	// You would need to iterate through guilds the bot is part of and check memberships individually.
	url := fmt.Sprintf("https://discord.com/api/users/%s/guilds", userID)
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
		return nil, fmt.Errorf("error fetching guild memberships, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var guildMemberships []GuildMember
	if err := json.Unmarshal(body, &guildMemberships); err != nil {
		return nil, err
	}

	return guildMemberships, nil
}
