package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/tg"
)

// Fetch messages from a group
func fetchChannelMessages(ctx context.Context, username string) ([]*tg.Message, error) {
	client, err := InitializeClient()
	if err != nil {
		log.Printf("Failed to initialize Telegram client: %v", err)
		return nil, err
	}

	channel, err := resolveChannelUsername(ctx, username) // Edit: Assign the second value to err
	if err != nil {
		return nil, err // Handle the error if resolveChannelUsername fails
	}
	var messagesSlice []*tg.Message // Define a slice to hold the messages

	err = client.Run(ctx, func(ctx context.Context) error {
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

func resolveChannelUsername(ctx context.Context, username string) (*tg.InputChannel, error) {
	client, err := InitializeClient()
	if err != nil {
		log.Printf("Failed to initialize Telegram client: %v", err)
		return nil, err
	}

	var channel *tg.InputChannel
	err = client.Run(ctx, func(ctx context.Context) error {
		resolved, err := client.API().ContactsResolveUsername(ctx, username)
		if err != nil {
			return err
		}

		channel = &tg.InputChannel{
			ChannelID:  resolved.Chats[0].GetID(),
			AccessHash: resolved.Chats[0].(*tg.Channel).AccessHash,
		}

		fmt.Printf("Channel ID: %d, Access Hash: %d\n", channel.ChannelID, channel.AccessHash)
		return nil
	})
	return channel, err
}
