package main

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	masa "github.com/masa-finance/masa-oracle/pkg"
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
	viper.SetDefault("BOOTNODES", "")
	viper.SetDefault("PORT_NBR", "4001")
	viper.SetDefault("UDP", true)
	viper.SetDefault("TCP", false)
	viper.SetDefault("STAKE_AMOUNT", "1000")
	// Add other default values as needed
	// log the flags
	bootnodesList := strings.Split(viper.GetString("BOOTNODES"), ",")
	logrus.Infof("1 Bootnodes: %v", bootnodesList)
	logrus.Infof("1 Port number: %d", viper.GetInt("PORT_NBR"))
	logrus.Infof("1 UDP: %v", viper.GetBool("UDP"))
	logrus.Infof("1 TCP: %v", viper.GetBool("TCP"))
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
			logrus.Error("could not create directory:")
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
