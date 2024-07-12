package telegram

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// Fetch messages from a group
func fetchChannelMessages(ctx context.Context, client *telegram.Client, channelID int64, accessHash int64) error {
	inputPeer := &tg.InputPeerChannel{ // Use InputPeerChannel instead of InputChannel
		ChannelID:  channelID,
		AccessHash: accessHash,
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
		fmt.Printf("Message ID: %d, Content: %s\n", message.ID, message.Message)
	}

	return nil
}

func resolveChannelUsername(ctx context.Context, client *telegram.Client, username string) (*tg.InputChannel, error) {
	resolved, err := client.API().ContactsResolveUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	channel := &tg.InputChannel{
		ChannelID:  resolved.Chats[0].GetID(),
		AccessHash: resolved.Chats[0].(*tg.Channel).AccessHash,
	}

	fmt.Printf("Channel ID: %d, Access Hash: %d\n", channel.ChannelID, channel.AccessHash)
	return channel, nil
}

func main() {
	// Define your Telegram app credentials
	appID := 28423325                             // Your actual app ID
	appHash := "c60c0a268973ea3f7d52e16e4ab2a0d3" // Your actual app hash

	// Define the path to the directory where session data will be stored
	sessionDir := filepath.Join(os.Getenv("HOME"), ".telegram-sessions")

	// Create the session directory if it doesn't already exist
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create session directory: %v\n", err)
		os.Exit(1)
	}

	// Create a session storage
	storage := &session.FileStorage{
		Path: filepath.Join(sessionDir, "session.json"),
	}

	// Create a random seed for the client
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate random seed: %v\n", err)
		os.Exit(1)
	}

	client := telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: storage,
	})

	// Continue with the rest of your application logic
	ctx := context.Background()

	// Check if the session file already exists and has content
	sessionExists := false
	if info, err := os.Stat(filepath.Join(sessionDir, "session.json")); err == nil && info.Size() > 0 {
		sessionExists = true
	}

	// Define the phone number and password (if 2FA is enabled)
	phone := "+13053398321"
	password := "" // Leave empty if 2FA is not enabled

	// Define the code prompt function
	codePrompt := func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
		fmt.Print("Enter code: ")
		code, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(code), nil
	}

	// Set up and perform the authentication flow only if no valid session exists
	if !sessionExists {
		authFlow := auth.NewFlow(
			auth.Constant(phone, password, auth.CodeAuthenticatorFunc(codePrompt)),
			auth.SendCodeOptions{},
		)
		if err := client.Run(ctx, func(ctx context.Context) error {
			return authFlow.Run(ctx, client.Auth())
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Authentication failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		// If a session exists, simply start the client
		if err := client.Run(ctx, func(ctx context.Context) error {
			username := "coinlistofficialchannel" // Replace with the actual username of the channel
			channel, err := resolveChannelUsername(ctx, client, username)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to resolve channel username: %v\n", err)
				os.Exit(1)
			}
			if err := fetchChannelMessages(ctx, client, channel.ChannelID, channel.AccessHash); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to fetch channel messages: %v\n", err)
				os.Exit(1)
			}
			return nil // No operation, just start the client
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start client with existing session: %v\n", err)
			os.Exit(1)
		}
	}

}
