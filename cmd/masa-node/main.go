package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/anthdm/hollywood/actor"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/staking"
	"github.com/sirupsen/logrus"
)

func main() {
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

	// start actor worker listener and subscribe if scraper and isStaked
	// go node.ActorEngine.Spawn(workers.NewWorker, "peer_worker", actor.WithID("peer"))
	if cfg.TwitterScraper || cfg.WebScraper {
		if isStaked {
			node.ActorEngine.Subscribe(actor.NewPID("0.0.0.0:4001", fmt.Sprintf("%s/%s", "peer_worker", "peer")))
		}
	}

	// tests SendWorkToPeers i.e. twitter, web
	// d, _ := json.Marshal(map[string]string{"request": "web", "url": "https://www.masa.finance", "depth": "2"})
	// d, _ := json.Marshal(map[string]string{"request": "twitter", "query": "$MASA token price", "count": "5"})
	// go workers.SendWorkToPeers(node, d)

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
	config.DisplayWelcomeMessage(multiAddr, ipAddr, keyManager.EthAddress, isStaked, isWriterNode, cfg.TwitterScraper, cfg.WebScraper)

	<-ctx.Done()
}
