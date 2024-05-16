package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/google/uuid"

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

	if cfg.StakeAmount != "" {
		err := handleStaking(keyManager.EcdsaPrivKey)
		if err != nil {
			logrus.Warningf("%v", err)
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

	var isWriterNode bool
	isWriterNode, _ = strconv.ParseBool(cfg.WriterNode)

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
		logrus.Infof("This node is set as the allowed peer with ID: %s and PubKey: %s", cfg.AllowedPeerId, cfg.AllowedPeerPublicKey)
	} else {
		logrus.Info("This node is not set as the allowed peer")
	}

	go db.InitResolverCache(node, keyManager)

	// Subscribe and if actor start monitoring actor workers
	go workers.SubscribeToWorkers(node)
	if node.IsActor() && isStaked {
		go workers.MonitorWorkers(ctx, node)
	}

	// WIP
	if os.Getenv("PG_URL") != "" {
		// run migrations
		_, err = db.ConnectToPostgres(true)
		if err != nil {
			logrus.Error(err)
		} else {
			uid := uuid.New().String()
			err := db.FireData(uid, []byte(`{"request":"twitter", "query":"$MASA", "count":5, "model": "gpt-4"}`), []byte(`{"tweets": ["twit", "twit"]}`))
			if err != nil {
				logrus.Error(err)
			}
			fErr := db.FireEvent(uid, []byte(`{"event":"Actor Started"}`))
			if fErr != nil {
				logrus.Error(fErr)
			}

			work, gErr := db.GetData(uid)
			if gErr != nil {
				logrus.Error(gErr)
			}

			for _, w := range work {
				jsonData, err := json.Marshal(w)
				if err != nil {
					logrus.Error("Failed to parse work into JSON: ", err)
				} else {
					logrus.Info(string(jsonData))
				}
			}
		}
	}
	// WIP

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
	config.DisplayWelcomeMessage(multiAddr, ipAddr, keyManager.EthAddress, isStaked, isWriterNode, cfg.TwitterScraper, cfg.DiscordScraper, cfg.WebScraper, config.Version)

	<-ctx.Done()
}
