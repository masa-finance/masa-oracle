// discordmessages.go
package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
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
func GetChannelMessages(channelID string, limit int, before string) ([]ChannelMessage, error) {
	// Print the parameters for debugging
	fmt.Printf("Parameters - channelID: %s, limit: %d, before: %s\n", channelID, limit, before)

	botToken := os.Getenv("DISCORD_BOT_TOKEN") // Replace with your actual environment variable name
	if botToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN environment variable not set")
	}

	url := fmt.Sprintf("https://discord.com/api/channels/%s/messages", channelID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters if they are provided
	q := req.URL.Query()
	if limit > 0 && limit <= 100 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	if before != "" {
		q.Add("before", before)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	// Print the full request URL for debugging
	fmt.Printf("Request URL: %s\n", req.URL.String())

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

// ScrapeDiscordMessagesForSentiment scrapes messages from a Discord channel and analyzes their sentiment.
func ScrapeDiscordMessagesForSentiment(channelID string, model string, prompt string) (string, string, error) {
	// Fetch messages from the Discord channel
	messages, err := GetChannelMessages(channelID, 100, "")
	if err != nil {
		return "", "", fmt.Errorf("error fetching messages from Discord channel: %v", err)
	}

	// Extract the content of the messages
	var messageContents []string
	for _, message := range messages {
		messageContents = append(messageContents, message.Content)
	}

	// Analyze the sentiment of the fetched messages
	// Note: Ensure that llmbridge.AnalyzeSentimentDiscord is implemented and can handle the analysis
	analysisPrompt, sentiment, err := llmbridge.AnalyzeSentimentDiscord(messageContents, model, prompt)
	if err != nil {
		return "", "", fmt.Errorf("error analyzing sentiment of Discord messages: %v", err)
	}
	return analysisPrompt, sentiment, nil

}
