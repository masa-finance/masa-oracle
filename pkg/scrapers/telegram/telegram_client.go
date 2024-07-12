package telegram

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
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
func StartAuthentication(ctx context.Context, phoneNumber string) (string, error) {
	// Initialize the Telegram client (if not already initialized)
	client, err := InitializeClient()
	if err != nil {
		log.Printf("Failed to initialize Telegram client: %v", err)
		return "", err
	}

	// Define a variable to hold the phoneCodeHash
	var phoneCodeHash string

	// Use client.Run to start the client and execute the SendCode method
	err = client.Run(ctx, func(ctx context.Context) error {
		// Call the SendCode method of the client to send the code to the user's Telegram app
		sentCode, err := client.Auth().SendCode(ctx, phoneNumber, auth.SendCodeOptions{
			AllowFlashCall: true,
			CurrentNumber:  true,
		})
		if err != nil {
			log.Printf("Error sending code: %v", err)
			return err
		}

		log.Printf("Code sent successfully to: %s", phoneNumber)

		// Extract the phoneCodeHash from the sentCode object
		switch code := sentCode.(type) {
		case *tg.AuthSentCode:
			phoneCodeHash = code.PhoneCodeHash
		default:
			return errors.New("unexpected type of AuthSentCode")
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to run client or send code: %v", err)
		return "", err
	}

	// Return the phoneCodeHash to be used in the next step
	log.Printf("Authentication process started successfully for: %s", phoneNumber)
	return phoneCodeHash, nil
}

// CompleteAuthentication uses the provided code to authenticate with Telegram.
func CompleteAuthentication(ctx context.Context, phoneNumber, code, phoneCodeHash string) (*tg.AuthAuthorization, error) {
	// Initialize the Telegram client (if not already initialized)
	client, err := InitializeClient()
	if err != nil {
		log.Printf("Failed to initialize Telegram client: %v", err)
		return nil, err
	}

	// Define a variable to hold the authentication result
	var authResult *tg.AuthAuthorization

	// Use client.Run to start the client and execute the SignIn method
	err = client.Run(ctx, func(ctx context.Context) error {
		// Use the provided code and phoneCodeHash to authenticate
		auth, err := client.Auth().SignIn(ctx, phoneNumber, code, phoneCodeHash)
		if err != nil {
			log.Printf("Error during SignIn: %v", err)
			return err
		}

		// At this point, authentication was successful, and you have the user's Telegram auth data.
		authResult = auth
		return nil
	})

	if err != nil {
		log.Printf("Failed to run client or sign in: %v", err)
		return nil, err
	}

	// You can now create a session for the user or perform other post-authentication tasks.
	log.Printf("Authentication successful for: %s", phoneNumber)
	return authResult, nil
}

func GetClient() *telegram.Client {
	return client
}
