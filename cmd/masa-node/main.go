package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/consensus"

	"github.com/masa-finance/masa-oracle/pkg/db"

	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

type NodeStatus struct {
	PeerID        string        `json:"peerId"`
	IsStaked      bool          `json:"isStaked"`
	TotalUpTime   time.Duration `json:"totalUpTime"`
	FirstLaunched time.Time     `json:"firstLaunched"`
	LastLaunched  time.Time     `json:"lastLaunched"`
}

type SharedData map[string]interface{}

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

	// WIP
	go db.InitResolverCache()

	data := []byte(node.Host.ID().String())
	signature, err := consensus.SignData(keyManager.Libp2pPrivKey, data)
	if err != nil {
		logrus.Errorf("%v", err)
	}
	_ = db.Verifier(node.Host, data, signature)

	// *** initialize dht database example for review ***

	//up := node.NodeTracker.GetNodeData(node.Host.ID().String())
	//totalUpTime := up.GetAccumulatedUptime()
	//status := NodeStatus{
	//	PeerID:        node.Host.ID().String(),
	//	IsStaked:      isStaked,
	//	TotalUpTime:   totalUpTime,
	//	FirstLaunched: time.Now().Add(-totalUpTime),
	//	LastLaunched:  time.Now(),
	//}
	//jsonData, _ := json.Marshal(status)
	//
	//keyStr := node.Host.ID().String() // for this nodes status data
	//
	//success, _ := db.WriteData(node, "/db/"+keyStr, jsonData, node.Host)
	//logrus.Printf("writeResult %+v", success)
	//
	//nodeVal := db.ReadData(node, "/db/"+keyStr, node.Host)
	//_ = json.Unmarshal(nodeVal, &status)
	//logrus.Printf("readResult: %+v\n", status)

	// example key for public shared data
	// requires its own struct if data is specific
	// or an empty SharedData struct for any data

	sharedData := SharedData{}
	sharedData["name"] = "John Doe"
	sharedData["age"] = 30
	sharedData["metadata"] = map[string]string{"twitter_handle": "@john"}

	jsonData, _ := json.Marshal(sharedData)

	keyStr := "twitter_data"

	success, _ := db.WriteData(node, "/db/"+keyStr, jsonData, node.Host)
	logrus.Printf("writeResult %+v", success)

	nodeVal := db.ReadData(node, "/db/"+keyStr, node.Host)
	_ = json.Unmarshal(nodeVal, &sharedData)
	logrus.Printf("readResult: %+v\n", sharedData)

	// *** initialize dht database example for review ***

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
