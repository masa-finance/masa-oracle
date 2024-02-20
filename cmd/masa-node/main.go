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
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/cicd_helpers"
	"github.com/masa-finance/masa-oracle/pkg/crypto"
	"github.com/masa-finance/masa-oracle/pkg/routes"
	"github.com/masa-finance/masa-oracle/pkg/staking"
	"github.com/masa-finance/masa-oracle/pkg/welcome"
)

func init() {

	// Initialize Viper
	viper.AutomaticEnv() // Read from environment variables
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // Optionally: add other paths, e.g., home directory or etc

	// Set default values
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILEPATH", "masa_oracle_node.log")
	viper.SetDefault("RPC_URL", "https://ethereum-sepolia.publicnode.com")
	viper.SetDefault("MASA_KEY_FILE_KEY", "private.key")
	viper.SetDefault("MASA_DIR", ".masa")
	viper.SetDefault("STAKE_AMOUNT", "1000")
	viper.SetDefault("BOOTNODES", "")
	viper.SetDefault("PORT_NBR", 4001)
	viper.SetDefault("UDP", true)
	viper.SetDefault("TCP", true)

	// Add other default values as needed

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %s", err)
		}
	}

	// Attempt to read in environment variables
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error loading .env file")
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

	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}
	keyFilePath := filepath.Join(usr.HomeDir, viper.GetString("MASA_DIR"), viper.GetString("MASA_KEY_FILE_KEY"))
	err = setUpFiles(keyFilePath)
	if err != nil {
		logrus.Fatal(err)
	}

	backupFileName := fmt.Sprintf("%s_%s", masa.Version, masa.NodeBackupFileName)
	err = os.Setenv(masa.NodeBackupPath, filepath.Join(usr.HomeDir, viper.GetString("MASA_DIR"), backupFileName))
	if err != nil {
		logrus.Error(err)
	}
}

func main() {

	// log the flags
	bootnodesList := strings.Split(bootnodes, ",")
	logrus.Infof("Bootnodes: %v", bootnodesList)
	logrus.Infof("Port number: %d", portNbr)
	logrus.Infof("UDP: %v", udp)
	logrus.Infof("TCP: %v", tcp)

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	privKey, ecdsaPrivKey, ethAddress, err := crypto.GetOrCreatePrivateKey(viper.GetString("MASA_KEY_FILE_KEY"))
	if err != nil {
		logrus.Fatal(err)
	}
	if stakeAmount != "" {
		// Exit after staking, do not proceed to start the node
		err = handleStaking(ecdsaPrivKey)
		if err != nil {
			logrus.Fatal(err)
		}
		os.Exit(0)
	}

	var isStaked bool
	// Verify the staking event
	isStaked, err = staking.VerifyStakingEvent(ethAddress)
	if err != nil {
		logrus.Error(err)
	}
	if !isStaked {
		logrus.Warn("No staking event found for this address")
	}

	// Pass the isStaked flag to the NewOracleNode function
	node, err := masa.NewOracleNode(ctx, privKey, portNbr, udp, tcp, isStaked)
	if err != nil {
		logrus.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		logrus.Fatal(err)
	}

	// Get the multiaddress and IP address of the node
	multiAddr := node.GetMultiAddrs().String() // Get the multiaddress
	ipAddr := node.Host.Addrs()[0].String()    // Get the IP address
	publicKeyHex, _ := crypto.GetPublicKeyForHost(node.Host)

	// Display the welcome message with the multiaddress and IP address
	welcome.DisplayWelcomeMessage(multiAddr, ipAddr, publicKeyHex, isStaked)

	// Set env variables for CI/CD pipelines
	cicd_helpers.SetEnvVariablesForPipeline(multiAddr)

	// Listen for SIGINT (CTRL+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Cancel the context when SIGINT is received
	go func() {
		<-c
		nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
		if nodeData != nil {
			nodeData.Left()
		}
		node.NodeTracker.DumpNodeData()
		cancel()
	}()

	// BP: Add gin router to get peers (multiaddress) and get peer addresses
	// @Bob - I am not sure if this is the right place for this to live if we end up building out more endpoints
	router := routes.SetupRoutes(node)
	go func() {
		err := router.Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	<-ctx.Done()
}

func setUpFiles(keyFilePath string) error {
	// Create the directories if they don't already exist
	if _, err := os.Stat(filepath.Dir(viper.GetString("MASA_KEY_FILE_PATH"))); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(viper.GetString("MASA_KEY_FILE_PATH")), 0755)
		if err != nil {
			logrus.Error("could not create directory:")
			return err
		}
	}

	return nil
}

// Add node type for startup notification of what kind of node you are running and what that means
