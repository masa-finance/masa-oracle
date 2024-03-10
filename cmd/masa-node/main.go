package main

import (
	"context"
	"encoding/json"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/db"

	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

func main() {
	cfg := config.GetInstance()
	cfg.LogConfig()
	cfg.SetupLogging()
	keyManager := masacrypto.KeyManagerInstance()

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	if cfg.StakeAmount != "" {
		// Exit after staking, do not proceed to start the node
		err := handleStaking(keyManager.EcdsaPrivKey)
		if err != nil {
			logrus.Fatal(err)
		}
		os.Exit(0)
	}

	var isStaked bool
	// Verify the staking event
	isStaked, err := staking.VerifyStakingEvent(keyManager.EthAddress)
	if err != nil {
		logrus.Error(err)
	}
	if !isStaked {
		logrus.Warn("No staking event found for this address")
	}

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

	go db.InitResolverCache(node)

	// *** Store NodeStatus ***
	data := []byte(node.Host.ID().String())
	signature, err := consensus.SignData(keyManager.Libp2pPrivKey, data)
	if err != nil {
		logrus.Errorf("%v", err)
	}
	_ = db.Verifier(node.Host, data, signature)

	up := node.NodeTracker.GetNodeData(node.Host.ID().String())
	if up != nil {
		totalUpTime := up.GetAccumulatedUptime()
		status := db.NodeStatus{
			PeerID:        node.Host.ID().String(),
			IsStaked:      isStaked,
			TotalUpTime:   totalUpTime,
			FirstLaunched: time.Now().Add(-totalUpTime),
			LastLaunched:  time.Now(),
		}
		jsonData, _ := json.Marshal(status)

		keyStr := node.Host.ID().String() // user ID for this nodes status key
		success, _ := db.WriteData(node, "/db/"+keyStr, jsonData)
		logrus.Infof("Store NodeStatus %+v", success)
	}
	// *** Store NodeStatus ***

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
		err := router.Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	// Get the multiaddress and IP address of the node
	multiAddr := node.GetMultiAddrs().String() // Get the multiaddress
	ipAddr := node.Host.Addrs()[0].String()    // Get the IP address
	// Display the welcome message with the multiaddress and IP address
	config.DisplayWelcomeMessage(multiAddr, ipAddr, keyManager.EthAddress, isStaked)

	<-ctx.Done()
}
