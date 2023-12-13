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

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/cicd_helpers"
	"github.com/masa-finance/masa-oracle/pkg/crypto"
	"github.com/masa-finance/masa-oracle/pkg/routes"
	"github.com/masa-finance/masa-oracle/pkg/staking"
	"github.com/masa-finance/masa-oracle/pkg/welcome"
)

func init() {
	f, err := os.OpenFile("masa_oracle_node.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)
	if os.Getenv("debug") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}
	envFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_node.env")
	keyFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_key")
	err = setUpFiles(envFilePath, keyFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	err = godotenv.Load(envFilePath)
	if err != nil {
		logrus.Error("Error loading .env file")
	}
	os.Setenv(masa.NodeBackupPath, filepath.Join(usr.HomeDir, ".masa", masa.NodeBackupFileName))
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

	privKey, ecdsaPrivKey, ethAddress, err := crypto.GetOrCreatePrivateKey(os.Getenv(masa.KeyFileKey))
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

	// Display the welcome message with the multiaddress and IP address
	welcome.DisplayWelcomeMessage(multiAddr, ipAddr)

	// Set env variables for CI/CD pipelines
	cicd_helpers.SetEnvVariablesForPipeline(multiAddr)

	// Listen for SIGINT (CTRL+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Cancel the context when SIGINT is received
	go func() {
		<-c
		node.NodeTracker.DumpNodeData()
		cancel()
	}()

	// BP: Add gin router to get peers (multiaddress) and get peer addresses
	// @Bob - I am not sure if this is the right place for this to live if we end up building out more endpoints
	router := routes.SetupRoutes(node)
	go router.Run()

	<-ctx.Done()
}

func setUpFiles(envFilePath, keyFilePath string) error {
	// Create the directories if they don't already exist
	if _, err := os.Stat(filepath.Dir(envFilePath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(envFilePath), 0755)
		if err != nil {
			logrus.Error("could not create directory:")
			return err
		}
	}
	// Check if the .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// If not, create it with default values
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("%s=%s\n", masa.KeyFileKey, keyFilePath))
		err = os.WriteFile(envFilePath, []byte(builder.String()), 0644)
		if err != nil {
			logrus.Error("could not write to .env file:")
			return err
		}
	}
	return nil
}
