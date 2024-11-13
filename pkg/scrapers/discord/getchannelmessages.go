// discordmessages.go
package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// ChannelMessage represents a Discord channel message structure
type ChannelMessage struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Author    struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
	} `json:"author"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// GetChannelMessages fetches messages for a specific channel from the Discord API
func GetChannelMessages(channelID string, limit string, before string) ([]ChannelMessage, error) {
	botToken := os.Getenv("DISCORD_BOT_TOKEN") // Replace with your actual environment variable name
	if botToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("https://discord.com/api/channels/%s/messages", channelID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	limitCheck, _ := strconv.Atoi(limit)

	// Add query parameters if they are provided
	q := req.URL.Query()
	if limitCheck > 0 && limitCheck <= 100 {
		q.Add("limit", limit)
	}
	if before != "" {
		q.Add("before", before)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching channel messages, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var messages []ChannelMessage
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
