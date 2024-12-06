package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/masa-finance/masa-oracle/internal/versioning"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("Log level set to Debug")

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		logrus.Infof("Masa Oracle Node Version: %s\nMasa Oracle Protocol verison: %s", versioning.ApplicationVersion, versioning.ProtocolVersion)
		os.Exit(0)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("[-] %v", err)
	}

	cfg.SetupLogging()
	cfg.LogConfig()

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	if cfg.Faucet {
		err := handleFaucet(cfg.RpcUrl, cfg.KeyManager.EcdsaPrivKey)
		if err != nil {
			logrus.Errorf("[-] %v", err)
			os.Exit(1)
		} else {
			logrus.Info("[+] Faucet event completed for this address")
			os.Exit(0)
		}
	}

	if cfg.StakeAmount != "" {
		err := handleStaking(cfg.RpcUrl, cfg.KeyManager.EcdsaPrivKey, cfg.StakeAmount)
		if err != nil {
			logrus.Warningf("%v", err)
		} else {
			logrus.Info("[+] Staking event completed for this address")
			os.Exit(0)
		}
	}

	// Verify the staking event
	isStaked, err := staking.VerifyStakingEvent(cfg.RpcUrl, cfg.KeyManager.EthAddress)
	if err != nil {
		logrus.Error(err)
	}

	if !isStaked {
		logrus.Warn("No staking event found for this address")
	}

	masaNodeOptions, workHandlerManager, pubKeySub := config.InitOptions(cfg)
	// Create a new OracleNode
	masaNode, err := node.NewOracleNode(ctx, masaNodeOptions...)

	if err != nil {
		logrus.Fatal(err)
	}

	if err = masaNode.Start(); err != nil {
		logrus.Fatal(err)
	}

	if cfg.AllowedPeer {
		cfg.AllowedPeerId = masaNode.Host.ID().String()
		cfg.AllowedPeerPublicKey = cfg.KeyManager.HexPubKey
		logrus.Infof("[+] Allowed peer with ID: %s and PubKey: %s", cfg.AllowedPeerId, cfg.AllowedPeerPublicKey)
	} else {
		logrus.Warn("[-] This node is not set as the allowed peer")
	}

	// Init cache resolver
	db.InitResolverCache(masaNode, cfg.KeyManager, cfg.AllowedPeerId, cfg.AllowedPeerPublicKey, cfg.Validator)

	// Cancel the context when SIGINT is received
	go handleSignals(cancel, masaNode, cfg)

	if cfg.APIEnabled {
		router := api.SetupRoutes(masaNode, workHandlerManager, pubKeySub)
		go func() {
			if err := router.Run(); err != nil {
				logrus.Fatal(err)
			}
		}()
		logrus.Info("API server started")
	} else {
		logrus.Info("API server is disabled")
	}

	// Get the multiaddress and IP address of the node
	multiAddrs, err := masaNode.GetP2PMultiAddrs()
	if err != nil {
		logrus.Errorf("[-] Error while getting node multiaddrs: %v", err)
	} else {
		config.DisplayWelcomeMessage(multiAddrs, cfg.KeyManager.EthAddress, isStaked, cfg.Validator, cfg.TwitterScraper, cfg.TelegramScraper, cfg.DiscordScraper, cfg.WebScraper, versioning.ApplicationVersion, versioning.ProtocolVersion)
	}

	<-ctx.Done()
}

func handleSignals(cancel context.CancelFunc, masaNode *node.OracleNode, cfg *config.AppConfig) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	nodeData := masaNode.NodeTracker.GetNodeData(masaNode.Host.ID().String())
	if nodeData != nil {
		nodeData.Left()
	}
	cancel()
	if cfg.TelegramStop != nil {
		if err := cfg.TelegramStop(); err != nil {
			logrus.Errorf("Error stopping the background connection: %v", err)
		}
	}
}
