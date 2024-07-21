package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/masa-finance/masa-oracle/pkg/workers"

	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("Masa Oracle Node Version: %s\n", config.Version)
		os.Exit(0)
	}

	cfg := config.GetInstance()
	cfg.LogConfig()
	cfg.SetupLogging()
	keyManager := masacrypto.KeyManagerInstance()

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	if cfg.Faucet {
		err := handleFaucet(keyManager.EcdsaPrivKey)
		if err != nil {
			logrus.Errorf("%v", err)
			os.Exit(1)
		} else {
			logrus.Info("[+] Faucet event completed for this address")
			os.Exit(0)
		}
	}

	if cfg.StakeAmount != "" {
		err := handleStaking(keyManager.EcdsaPrivKey)
		if err != nil {
			logrus.Warningf("%v", err)
		} else {
			logrus.Info("[+] Staking event completed for this address")
			os.Exit(0)
		}
	}

	// Verify the staking event
	isStaked, err := staking.VerifyStakingEvent(keyManager.EthAddress)
	if err != nil {
		logrus.Error(err)
	}

	if !isStaked {
		logrus.Warn("No staking event found for this address")
	}

	isValidator := cfg.Validator

	// Create a new OracleNode
	node, err := masa.NewOracleNode(ctx, isStaked)
	if err != nil {
		logrus.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		logrus.Fatal(err)
	}

	if cfg.AllowedPeer {
		cfg.AllowedPeerId = node.Host.ID().String()
		cfg.AllowedPeerPublicKey = keyManager.HexPubKey
		logrus.Infof("[+] Allowed peer with ID: %s and PubKey: %s", cfg.AllowedPeerId, cfg.AllowedPeerPublicKey)
	} else {
		logrus.Warn("[-] This node is not set as the allowed peer")
	}

	// Init cache resolver
	db.InitResolverCache(node, keyManager)

	// Subscribe and if actor start monitoring actor workers
	// considering all that matters is if the node is staked
	// and other peers can do work we only need to check this here
	// if this peer can or cannot scrape or write that is checked in other places
	if node.IsStaked {
		go workers.SubscribeToWorkers(node)
		go workers.MonitorWorkers(ctx, node)
		go masa.SubscribeToBlocks(ctx, node)
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
		cancel()
	}()

	router := api.SetupRoutes(node)
	go func() {
		err = router.Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	// Get the multiaddress and IP address of the node
	multiAddr := node.GetMultiAddrs().String() // Get the multiaddress
	ipAddr := node.Host.Addrs()[0].String()    // Get the IP address
	// Display the welcome message with the multiaddress and IP address
	config.DisplayWelcomeMessage(multiAddr, ipAddr, keyManager.EthAddress, isStaked, isValidator, cfg.TwitterScraper, cfg.DiscordScraper, cfg.WebScraper, config.Version)

	<-ctx.Done()
}
