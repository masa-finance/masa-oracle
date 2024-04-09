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
	"github.com/masa-finance/masa-oracle/pkg/proto/msg"
	"github.com/masa-finance/masa-oracle/pkg/staking"
	"github.com/sirupsen/logrus"
)

type foo struct{}

func newFoo() actor.Receiver {
	return &foo{}
}

func (f *foo) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		fmt.Println("actor started")
	case actor.Stopped:
		fmt.Println("actor stopped")
	case *msg.Message:
		fmt.Println("actor has received", msg.Data)
	}
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

	// var isStaked bool
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

	// WIP use actor engine from global level node object
	go func() {
		pid := node.ActorEngine.Spawn(newFoo, "my_foo_actor", actor.WithID(node.Host.ID().String()))
		for i := 0; i < 3; i++ {
			node.ActorEngine.Send(pid, &msg.Message{Data: "hello world!"})
			getpid := node.ActorEngine.Registry.GetPID("my_foo_actor", node.Host.ID().String())
			fmt.Println(getpid)
			peerPID := actor.NewPID("192.168.4.164:4001", "my_foo_actor/peer")
			node.ActorEngine.Send(peerPID, &msg.Message{Data: "hello 164!"})
		}
		// node.ActorEngine.Poison(pid).Wait() // use this where we want to stop a actor listener
	}()
	// WIP use actor engine on node

	// WIP web scrape
	//go func() {
	//	res, err := scraper.ScrapeWebDataUsingActors([]string{"https://en.wikipedia.org/wiki/Badger"}, 5)
	//	if err != nil {
	//		logrus.Errorf("Error collecting data: %s", err.Error())
	//		return
	//	}
	//	logrus.Infof("%+v", res)
	//}()
	// WIP web

	// WIP testing db
	// type Sentiment struct {
	// 	ConversationId int64
	// 	Tweet          string
	// 	PromptId       int64
	// }

	// // IMPORTANT migrations true will drop all
	// database, err := db.ConnectToPostgres(false)
	// if err != nil {
	// 	logrus.Errorf(err)
	// }
	// defer database.Close()

	// data := []Sentiment{}
	// query := `SELECT "conversation_id", "tweet", "prompt_id" FROM sentiment`
	// rows, err := database.Query(query)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// var (
	// 	conversationId int64
	// 	tweet          string
	// 	promptId       int64
	// )

	// for rows.Next() {
	// 	if err = rows.Scan(&conversationId, &tweet, &promptId); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	data = append(data, Sentiment{conversationId, tweet, promptId})
	// }
	// fmt.Println(data)
	// WIP testing

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
	config.DisplayWelcomeMessage(multiAddr, ipAddr, keyManager.EthAddress, isStaked, isWriterNode)

	<-ctx.Done()
}
