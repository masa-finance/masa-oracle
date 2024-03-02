package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgraph-io/badger/v4"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/badgerdb"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/crypto"
	"github.com/masa-finance/masa-oracle/pkg/routes"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

func main() {
	cfg := config.GetInstance()
	cfg.LogConfig()

	//@dBP added initialization for badgerdb, we need to add the DB_PATH to the viper configs in /cmd/masa-node/db.go
	// Initialize the database
	db, err := badgerdb.InitializeDB()
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer func(db *badger.DB) {
		err := db.Close()
		if err != nil {
			logrus.Errorf("Failed to close the database: %v", err)
		}
	}(db) // This ensures the database is properly closed on application exit

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	privKey, ecdsaPrivKey, ethAddress, err := crypto.GetOrCreatePrivateKey(cfg.PrivateKeyFile)

	if err != nil {
		logrus.Fatal(err)
	}
	if cfg.StakeAmount != "" {
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
	node, err := masa.NewOracleNode(ctx, privKey, isStaked)
	if err != nil {
		logrus.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		logrus.Fatal(err)
	}

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

	// Get the multiaddress and IP address of the node
	multiAddr := node.GetMultiAddrs().String() // Get the multiaddress
	ipAddr := node.Host.Addrs()[0].String()    // Get the IP address
	publicKeyHex, _ := crypto.GetPublicKeyForHost(node.Host)

	// Display the welcome message with the multiaddress and IP address
	config.DisplayWelcomeMessage(multiAddr, ipAddr, publicKeyHex, isStaked)

	<-ctx.Done()

}

// Add node type for startup notification of what kind of node you are running and what that means
