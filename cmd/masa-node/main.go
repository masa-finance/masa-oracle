package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/masa-finance/masa-oracle/pkg/consensus"
	"github.com/masa-finance/masa-oracle/pkg/db"

	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/api"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	myNetwork "github.com/masa-finance/masa-oracle/pkg/network"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

type NodeStatus struct {
	PeerID        string        `json:"peerId"`
	IsStaked      bool          `json:"isStaked"`
	TotalUpTime   time.Duration `json:"totalUpTime"`
	FirstLaunched time.Time     `json:"firstLaunched"`
	LastLaunched  time.Time     `json:"lastLaunched"`
}

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

	peers, _ := myNetwork.GetBootNodesMultiAddress(config.GetInstance().Bootnodes)
	for _, b := range peers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(b)
		logrus.Println(peerInfo)
		_ = node.Host.Connect(ctx, *peerInfo)
	}

	data := []byte(node.Host.ID().String())
	signature, err := consensus.SignData(keyManager.Libp2pPrivKey, data)
	if err != nil {
		logrus.Errorf("%v", err)
	}

	// *** initialize dht database example for review ***

	_ = db.Verifier(node.Host, data, signature)

	up := node.NodeTracker.GetNodeData(node.Host.ID().String())
	totalUpTime := up.GetAccumulatedUptime()
	status := NodeStatus{
		PeerID:        node.Host.ID().String(),
		IsStaked:      isStaked,
		TotalUpTime:   totalUpTime,
		FirstLaunched: time.Now().Add(-totalUpTime),
		LastLaunched:  time.Now(),
	}
	jsonData, _ := json.Marshal(status)

	success, _ := db.WriteData(node, "/db/"+node.Host.ID().String(), jsonData, node.Host)
	logrus.Printf("writeResult %+v", success)

	nodeVal := db.ReadData(node, "/db/"+node.Host.ID().String(), node.Host)
	_ = json.Unmarshal(nodeVal, &status)
	logrus.Printf("readResult: %+v\n", status)

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
