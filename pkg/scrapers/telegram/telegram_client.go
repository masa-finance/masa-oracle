package telegram

import (
	"context"
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

var (
	client     *telegram.Client
	once       sync.Once
	appID      = 28423325                           // Your actual app ID
	appHash    = "c60c0a268973ea3f7d52e16e4ab2a0d3" // Your actual app hash
	sessionDir = filepath.Join(os.Getenv("HOME"), ".telegram-sessions")
)

func InitializeClient() (*telegram.Client, error) {
	var err error
	once.Do(func() {
		// Ensure the session directory exists
		if err = os.MkdirAll(sessionDir, 0700); err != nil {
			return
		}

		// Create a session storage
		storage := &session.FileStorage{
			Path: filepath.Join(sessionDir, "session.json"),
		}

		// Create a random seed for the client
		seed := make([]byte, 32)
		if _, err = rand.Read(seed); err != nil {
			return
		}

		// Initialize the Telegram client
		client = telegram.NewClient(appID, appHash, telegram.Options{
			SessionStorage: storage,
		})
	})

	return client, err
}

// func AuthenticateUser(ctx context.Context, phoneNumber string) error {
// 	// Check if the session file already exists and has content
// 	sessionExists := false
// 	if info, err := os.Stat(filepath.Join(sessionDir, "session.json")); err == nil && info.Size() > 0 {
// 		sessionExists = true
// 	}

// 	if !sessionExists {
// 		// Define the code prompt function
// 		codePrompt := func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
// 			fmt.Print("Enter code: ")
// 			code, err := bufio.NewReader(os.Stdin).ReadString('\n')
// 			if err != nil {
// 				return "", err
// 			}
// 			return strings.TrimSpace(code), nil
// 		}

// 		// Set up and perform the authentication flow
// 		authFlow := auth.NewFlow(
// 			auth.Constant(phoneNumber, "", auth.CodeAuthenticatorFunc(codePrompt)), // Empty password for simplicity
// 			auth.SendCodeOptions{},
// 		)
// 		if err := client.Run(ctx, func(ctx context.Context) error {
// 			return authFlow.Run(ctx, client.Auth())
// 		}); err != nil {
// 			return fmt.Errorf("authentication failed: %v", err)
// 		}
// 	}

// 	return nil
// }

// StartAuthentication sends the phone number to Telegram and requests a code.
func StartAuthentication(ctx context.Context, client *telegram.Client, phoneNumber string) (string, error) {
	// Call the SendCode method of the client to send the code to the user's Telegram app
	sentCode, err := client.Auth().SendCode(ctx, phoneNumber, auth.SendCodeOptions{})
	if err != nil {
		return "", err
	}

	// Extract the phoneCodeHash from the sentCode object
	var phoneCodeHash string
	switch code := sentCode.(type) {
	case *tg.AuthSentCode:
		phoneCodeHash = code.PhoneCodeHash
	default:
		return "", errors.New("unexpected type of AuthSentCode")
	}

	// Return the phoneCodeHash to be used in the next step
	return phoneCodeHash, nil
}

// CompleteAuthentication uses the provided code to authenticate with Telegram.
func CompleteAuthentication(ctx context.Context, client *telegram.Client, phoneNumber, code, phoneCodeHash string) (*tg.AuthAuthorization, error) {
	// Use the provided code and phoneCodeHash to authenticate
	auth, err := client.Auth().SignIn(ctx, phoneNumber, code, phoneCodeHash)
	if err != nil {
		// Handle the specific error if password authentication is needed
		if err == errors.New("2FA required") {
			// Here you would handle the second factor authentication (2FA)
			// This usually involves prompting the user for their password.
		}
		return nil, err
	}

	// At this point, authentication was successful, and you have the user's Telegram auth data.
	// You can now create a session for the user or perform other post-authentication tasks.

	return auth, nil
}

func GetClient() *telegram.Client {
	return client
}
