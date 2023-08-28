package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	f, err := os.OpenFile("masa_oracle_node.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)
	logrus.SetLevel(logrus.DebugLevel)

	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}
	envFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_node.env")
	keyFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_key")

	// Create the directories if they don't already exist
	if _, err := os.Stat(filepath.Dir(envFilePath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(envFilePath), 0755)
		if err != nil {
			logrus.Fatal("could not create directory:", err)
		}
	}
	// Check if the .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// If not, create it with default values
		defaultEnv := fmt.Sprintf("%s=%s\n", keyFileKey, keyFilePath)
		err = os.WriteFile(envFilePath, []byte(defaultEnv), 0644)
		if err != nil {
			logrus.Fatal("could not write to .env file:", err)
		}
	}
	err = godotenv.Load(envFilePath)
	if err != nil {
		logrus.Error("Error loading .env file")
	}
}

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Listen for SIGINT (CTRL+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Cancel the context when SIGINT is received
	go func() {
		<-c
		cancel()
	}()

	privKey, err := getOrCreatePrivateKey(os.Getenv(keyFileKey))
	if err != nil {
		logrus.Fatal(err)
	}
	certFile := os.Getenv(cert)
	keyFile := os.Getenv(certPem)

	// Check if the certificate and key files already exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			// If not, generate them
			if err := generateSelfSignedCert(certFile, keyFile); err != nil {
				logrus.Fatal("Failed to generate self-signed certificate:", err)
			}
		}
	}

	node, err := NewOracleNode(privKey)
	if err != nil {
		logrus.Fatal(err)
	}
	err = node.Start(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	<-ctx.Done()
}
