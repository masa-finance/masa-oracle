package main

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	masa "github.com/masa-finance/masa-oracle/pkg"
	util "github.com/masa-finance/masa-oracle/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// Set up masa file path based on current user and config settings
	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}

	// Set default values and use constants for values used elsewhere in the application
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILEPATH", "masa_oracle_node.log")
	viper.SetDefault(masa.PrivKeyFile, "masa_oracle_key")
	viper.SetDefault(masa.MasaDir, filepath.Join(usr.HomeDir, ".masa"))
	viper.SetDefault(masa.RpcUrl, "https://ethereum-sepolia.publicnode.com")
	viper.SetDefault(masa.BootNodes, "")
	viper.SetDefault("PORT_NBR", "4001")
	viper.SetDefault("UDP", true)
	viper.SetDefault("TCP", false)
	viper.SetDefault("STAKE_AMOUNT", "1000")
	viper.SetDefault("allowedPeer", false)

	// Load or generate the private key using utility function
	privKey, err := util.LoadPrivateKey()
	if err != nil {
		log.Fatalf("Failed to load or create private key: %v", err)
	}

	// Derive the hex-encoded public key using utility function
	allowedPeerPubKeyHex, err := util.GetHexEncodedPublicKey(privKey)
	if err != nil {
		log.Fatalf("Failed to get hex-encoded public key: %v", err)
	}

	// Get the PeerID using utility function
	peerIDString, err := util.GetPeerIDFromPubKey(privKey.GetPublic())
	if err != nil {
		log.Fatalf("Failed to convert public key to peer ID: %v", err)
	}

	if viper.GetBool("allowedPeer") {
		viper.Set("ALLOWED_PEER_ID", peerIDString)
		viper.Set("ALLOWED_PEER_PUBKEY", allowedPeerPubKeyHex)
		logrus.Infof("This node is set as the allowed peer with ID: %s and PubKey: %s", peerIDString, allowedPeerPubKeyHex)
	} else {
		logrus.Info("This node is not set as the allowed peer. Skipping setting ALLOWED_PEER_ID and ALLOWED_PEER_PUBKEY.")
	}

	// Log the flags
	bootnodesList := strings.Split(viper.GetString(masa.BootNodes), ",")
	logrus.Infof("Bootnodes: %v", bootnodesList)
	logrus.Infof("Port number: %d", viper.GetInt("PORT_NBR"))
	logrus.Infof("UDP: %v", viper.GetBool("UDP"))
	logrus.Infof("TCP: %v", viper.GetBool("TCP"))

	// Check for env vars, config files, in order to override above defaults
	viper.AutomaticEnv() // Read from environment variables
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // Optionally: add other paths, e.g., home directory or etc

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %s", err)
		}
	}

	if _, err := os.Stat(filepath.Dir(viper.GetString(masa.MasaDir))); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(viper.GetString(masa.MasaDir)), 0755)
		if err != nil {
			logrus.Error("could not create directory:", err)
		}
	}

	// Open output file for logging
	f, err := os.OpenFile(viper.GetString("LOG_FILEPATH"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)

	if viper.GetString("LOG_LEVEL") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Infof("2 Bootnodes: %v", bootnodesList)
	logrus.Infof("2 Port number: %d", viper.GetInt("PORT_NBR"))
	logrus.Infof("2 UDP: %v", viper.GetBool("UDP"))
	logrus.Infof("2 TCP: %v", viper.GetBool("TCP"))
}
