package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/tg"

	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
)

// Fetch messages from a group
func FetchChannelMessages(username string) ([]*tg.Message, error) {
	// Initialize the Telegram client (if not already initialized)
	// client = GetClient()

	if client == nil {
		var err error
		client, err = InitializeClient()
		if err != nil {
			log.Printf("Failed to initialize Telegram client: %v", err)
			return nil, err
		}
	}

	var messagesSlice []*tg.Message // Define a slice to hold the messages
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := client.Run(ctx, func(ctx context.Context) error {
		resolved, err := client.API().ContactsResolveUsername(ctx, username)
		if err != nil {
			return err
		}

		channel := &tg.InputChannel{
			ChannelID:  resolved.Chats[0].GetID(),
			AccessHash: resolved.Chats[0].(*tg.Channel).AccessHash,
		}

		fmt.Printf("Channel ID: %d, Access Hash: %d\n", channel.ChannelID, channel.AccessHash)

		inputPeer := &tg.InputPeerChannel{ // Use InputPeerChannel instead of InputChannel
			ChannelID:  channel.ChannelID,
			AccessHash: channel.AccessHash,
		}
		result, err := client.API().MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
			Peer:  inputPeer, // Pass inputPeer here
			Limit: 100,       // Adjust the number of messages to fetch
		})
		if err != nil {
			return err
		}

		// Type assert the result to *tg.MessagesChannelMessages to access Messages field
		messages, ok := result.(*tg.MessagesChannelMessages)
		if !ok {
			return fmt.Errorf("unexpected type %T", result)
		}

		// Process the messages
		for _, m := range messages.Messages {
			message, ok := m.(*tg.Message) // Type assert to *tg.Message
			if !ok {
				// Handle the case where the message is not a regular message (e.g., service message)
				continue
			}
			messagesSlice = append(messagesSlice, message) // Append the message to the slice
		}
		return nil
	})

	return messagesSlice, err // Return the slice of messages and any error
}

// ScrapeTelegramMessagesForSentiment scrapes messages from a Telegram channel and analyzes their sentiment.
func ScrapeTelegramMessagesForSentiment(username string, model string, prompt string) (string, string, error) {
	// Fetch messages from the Telegram channel
	messages, err := FetchChannelMessages(username)
	if err != nil {
		return "", "", fmt.Errorf("error fetching messages from Telegram channel: %v", err)
	}

	// Analyze the sentiment of the fetched messages
	// Note: Ensure that llmbridge.AnalyzeSentimentTelegram is implemented and can handle the analysis
	analysisPrompt, sentiment, err := llmbridge.AnalyzeSentimentTelegram(messages, model, prompt)
	if err != nil {
		return "", "", fmt.Errorf("error analyzing sentiment of Telegram messages: %v", err)
	}
	return analysisPrompt, sentiment, nil
}
