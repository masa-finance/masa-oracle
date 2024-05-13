package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/masa-finance/masa-oracle/pkg/workers"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/staking"
	"github.com/sirupsen/logrus"
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

	// Start monitoring actor workers
	if cfg.TwitterScraper || cfg.WebScraper {
		if isStaked {
			go workers.MonitorWorkers(ctx, node)
		}
	}

	// WIP
	//models, _ := config.GetCloudflareModels()
	//logrus.Info(models)

	if os.Getenv("PG_URL") != "" {
		type Work struct {
			id      int64
			payload json.RawMessage
			raw     json.RawMessage
		}

		// IMPORTANT migrations true will drop all
		database, err := db.ConnectToPostgres(false)
		if err != nil {
			logrus.Error(err)
		}
		defer database.Close()

		insertQuery := `INSERT INTO "public"."work" ("payload", "raw") VALUES ($1, $2)`
		payloadJSON := json.RawMessage(`{"request":"twitter", "query":"$MASA", "count":5}`)
		rawJSON := json.RawMessage(`{"tweets": []}`)
		_, err = database.Exec(insertQuery, payloadJSON, rawJSON)
		if err != nil {
			logrus.Error(err)
		}

		data := []Work{}
		query := `SELECT "id", "payload", "raw" FROM "public"."work"`
		rows, err := database.Query(query)
		if err != nil {
			logrus.Error(err)
		}
		defer rows.Close()

		var (
			id      int64
			payload json.RawMessage
			raw     json.RawMessage
		)

		for rows.Next() {
			if err = rows.Scan(&id, &payload, &raw); err != nil {
				log.Fatal(err)
			}
			data = append(data, Work{id, payload, raw})
		}
		logrus.Infof("record from pg %s", data[0].payload)
	}

	// JWT
	// jwtToken, err := consensus.GenerateJWTToken(node.Host.ID().String())
	// if err != nil {
	// 	logrus.Error(err)
	// }
	// logrus.Infof("jwt token: %s", jwtToken)
	// JWT

	// PoW
	// apiKey := consensus.GeneratePoW(node.Host.ID().String())
	// logrus.Infof("api key: %s", apiKey)
	// PoW

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
	config.DisplayWelcomeMessage(multiAddr, ipAddr, keyManager.EthAddress, isStaked, isWriterNode, cfg.TwitterScraper, cfg.WebScraper, config.Version)

	<-ctx.Done()
}
