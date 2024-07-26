package telegram

import (
	"context"
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
)

var (
	client     *telegram.Client
	once       sync.Once
	appID      = 28423325                           // Your actual app ID
	appHash    = "c60c0a268973ea3f7d52e16e4ab2a0d3" // Your actual app hash
	sessionDir = filepath.Join(os.Getenv("HOME"), ".telegram-sessions")
)

func GetClient() (*telegram.Client, error) {
	var err error
	once.Do(func() {
		// Ensure the session directory exists
		if err = os.MkdirAll(sessionDir, 0700); err != nil {
			logrus.Error(err)
			return
		}

		// Create a session storage
		storage := &session.FileStorage{
			Path: filepath.Join(sessionDir, "session.json"),
		}

		// Create a random seed for the client
		seed := make([]byte, 32)
		if _, err = rand.Read(seed); err != nil {
			logrus.Error(err)
			return
		}

		// Initialize the Telegram client
		client = telegram.NewClient(appID, appHash, telegram.Options{
			SessionStorage: storage,
		})
	})
	return client, err
}

// StartAuthentication sends the phone number to Telegram and requests a code.
func StartAuthentication(ctx context.Context, phoneNumber string) (string, error) {
	// Initialize the Telegram client (if not already initialized)
	client, err := GetClient()
	if err != nil {
		logrus.Errorf("Failed to initialize Telegram client: %v", err)
		return "", err
	}
	sentCode, err := client.Auth().SendCode(ctx, phoneNumber, auth.SendCodeOptions{
		AllowFlashCall: true,
		CurrentNumber:  true,
	})
	if err != nil {
		logrus.Errorf("Error sending code: %v", err)
		return "", err
	}
	logrus.Debugf("Code sent successfully to: %s", phoneNumber)

	// Extract the phoneCodeHash from the sentCode object
	var phoneCodeHash string
	if code, ok := sentCode.(*tg.AuthSentCode); ok {
		phoneCodeHash = code.PhoneCodeHash
	} else {
		return "", errors.New("unexpected type of AuthSentCode")
	}

	logrus.Infof("Authentication process started successfully for: %s", phoneNumber)
	return phoneCodeHash, nil
}

// CompleteAuthentication uses the provided code to authenticate with Telegram.
func CompleteAuthentication(ctx context.Context, phoneNumber, code, phoneCodeHash string) (*tg.AuthAuthorization, error) {
	// Initialize the Telegram client (if not already initialized)
	client, err := GetClient()
	if err != nil {
		logrus.Printf("Failed to initialize Telegram client: %v", err)
		return nil, err // Edit: Added nil as the first return value
	}

	authResult, err := client.Auth().SignIn(ctx, phoneNumber, code, phoneCodeHash)
	if err != nil {
		logrus.Printf("Error during SignIn: %v", err)
		return nil, err
	}

	cfg := config.GetInstance()
	cfg.TelegramStop, err = bg.Connect(client)
	if err != nil {
		return nil, err
	}

	// Now you can use client.
	if _, err := client.Auth().Status(ctx); err != nil {
		logrus.Printf("Failed to run client or sign in: %v", err)
	}
	return authResult, nil
}
