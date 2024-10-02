package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/multiformats/go-multiaddr"

	"github.com/masa-finance/masa-oracle/internal/versioning"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("Log level set to Debug")

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		logrus.Infof("Masa Oracle Node Version: %s\nMasa Oracle Protocol verison: %s", versioning.ApplicationVersion, versioning.ProtocolVersion)
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
			logrus.Errorf("[-] %v", err)
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

	masaNodeOptions, workHandlerManager, pubKeySub := initOptions(cfg)
	// Create a new OracleNode
	masaNode, err := node.NewOracleNode(ctx, masaNodeOptions...)

	if err != nil {
		logrus.Fatal(err)
	}

	err = masaNode.Start()
	if err != nil {
		logrus.Fatal(err)
	}

	if cfg.TwitterScraper && cfg.DiscordScraper && cfg.WebScraper {
		logrus.Warn("[+] Node is set as all types of scrapers. This may not be intended behavior.")
	}

	if cfg.AllowedPeer {
		cfg.AllowedPeerId = masaNode.Host.ID().String()
		cfg.AllowedPeerPublicKey = keyManager.HexPubKey
		logrus.Infof("[+] Allowed peer with ID: %s and PubKey: %s", cfg.AllowedPeerId, cfg.AllowedPeerPublicKey)
	} else {
		logrus.Warn("[-] This node is not set as the allowed peer")
	}

	// Init cache resolver
	db.InitResolverCache(masaNode, keyManager)

	// Listen for SIGINT (CTRL+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Cancel the context when SIGINT is received
	go func() {
		<-c
		nodeData := masaNode.NodeTracker.GetNodeData(masaNode.Host.ID().String())
		if nodeData != nil {
			nodeData.Left()
		}
		cancel()
		// Call the global StopFunc to stop the Telegram background connection
		cfg := config.GetInstance()
		if cfg.TelegramStop != nil {
			if err := cfg.TelegramStop(); err != nil {
				logrus.Errorf("Error stopping the background connection: %v", err)
			}
		}
	}()

	router := api.SetupRoutes(masaNode, workHandlerManager, pubKeySub)
	go func() {
		err = router.Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	// Get the multiaddress and IP address of the node
	multiAddr := masaNode.GetMultiAddrs()                      // Get the multiaddress
	ipAddr, err := multiAddr.ValueForProtocol(multiaddr.P_IP4) // Get the IP address
	// Display the welcome message with the multiaddress and IP address
	config.DisplayWelcomeMessage(multiAddr.String(), ipAddr, keyManager.EthAddress, isStaked, isValidator, cfg.TwitterScraper, cfg.TelegramScraper, cfg.DiscordScraper, cfg.WebScraper, versioning.ApplicationVersion, versioning.ProtocolVersion)

	<-ctx.Done()
}
