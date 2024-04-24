package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	actor "github.com/asynkron/protoactor-go/actor"
	messages "github.com/masa-finance/masa-oracle/pkg/workers/messages"

	"github.com/masa-finance/masa-oracle/pkg/pubsub"

	"github.com/ipfs/go-cid"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

var (
	clients        = actor.NewPIDSet()
	workerStatusCh = make(chan []byte)
)

type Worker struct{}

// NewWorker creates a new instance of the Worker actor.
// It implements the actor.Receiver interface, allowing it to receive and handle messages.
//
// Returns:
//   - An instance of the Worker struct that implements the actor.Receiver interface.
func NewWorker() actor.Producer {
	return func() actor.Actor {
		return &Worker{}
	}
}

// Receive is the message handling method for the Worker actor.
// It receives messages through the actor context and processes them based on their type.
func (a *Worker) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.Connect:
		log.Printf("Worker %v connected", msg.Sender)
		clients.Add(msg.Sender)
	case *actor.Started:
		log.Println("actor started")
	case *actor.Stopped:
		log.Println("actor stopped")
	case *messages.Work:
		fmt.Printf("%v", msg.Data)
		// broadcast(ctx, clients, &messages.Response{Value: "some value"})
	}
}

// computeCid calculates the CID (Content Identifier) for a given string.
//
// Parameters:
//   - str: The input string for which to compute the CID.
//
// Returns:
//   - string: The computed CID as a string.
//   - error: An error, if any occurred during the CID computation.
//
// The function uses the multihash package to create a SHA2-256 hash of the input string.
// It then creates a CID (version 1) from the multihash and returns the CID as a string.
// If an error occurs during the multihash computation or CID creation, it is returned.
func computeCid(str string) (string, error) {
	// Create a multihash from the string
	mhHash, err := mh.Sum([]byte(str), mh.SHA2_256, -1)
	if err != nil {
		return "", err
	}
	// Create a CID from the multihash
	cidKey := cid.NewCidV1(cid.Raw, mhHash).String()
	return cidKey, nil
}

// SendWorkToPeers sends work data to peer nodes in the network.
// It subscribes to the local actor engine and the actor engines of peer nodes.
// The work data is then broadcast as an event to all subscribed nodes.
//
// Parameters:
//   - node: A pointer to the OracleNode instance.
//   - data: The work data to be sent, as a byte slice.
// func SendWorkToPeers(node *masa.OracleNode, data []byte) {
// peers := node.Host.Network().Peers()
// for _, peer := range peers {
// 	conns := node.Host.Network().ConnsToPeer(peer)
// 	for _, conn := range conns {
// 		addr := conn.RemoteMultiaddr()
// 		ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
// 		peerPID := actor.NewPID(fmt.Sprintf("%s:4001", ipAddr), fmt.Sprintf("%s/%s", "peer_worker", "peer"))
// 		fmt.Println(peerPID)
// 		node.ActorEngine.Subscribe(peerPID)
// 	}
// }
//node.ActorEngine.BroadcastEvent(&msg.Message{Data: string(data)})
// go MonitorWorkers(context.Background(), node)
// }

// StartWorkers starts the worker goroutines for the given OracleNode.
// It spawns a new actor using the node's ActorEngine and sets up message handling
// for the actor. The actor listens for Started and Message events.
// When a Message event is received, it parses the work data and performs the
// corresponding task (web scraping or Twitter scraping) based on the request type.
// The collected data is then sent to the workerStatusCh channel.
func StartWork(node *masa.OracleNode) {

	// props := actor.PropsFromProducer(NewWorker())
	// pid := node.ActorEngine.Spawn(props)
	// message := &messages.Work{Data: "hi-1", Sender: pid}
	// fmt.Println("debug ", message)

	// spawnResponse, err := node.ActorRemote.SpawnNamed("127.0.0.1:4002", "worker", "peer", time.Second)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {

	// 	// get spawned PID
	// 	spawnedPID := spawnResponse.Pid
	// 	client := node.ActorEngine.Spawn(props)
	// 	node.ActorEngine.Send(spawnedPID, &messages.Connect{
	// 		Sender: client,
	// 	})

	// 	for i := 0; i < 10; i++ {
	// 		node.ActorEngine.Send(spawnedPID, message)
	// 	}
	// }
}

// func StartWorkersOLD(node *masa.OracleNode) {
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)

// 	node.ActorEngine.SpawnFunc(func(c *actor.Context) {
// 		switch m := c.Message().(type) {
// 		case actor.Started:
// 			// fmt.Println(c.PID(), c.PID().Address)
// 			// c.Engine().Subscribe(c.PID())
// 			//peerPID := actor.NewPID(fmt.Sprintf("%s", c.PID().Address), fmt.Sprintf("%s/%s", "peer_worker", "peer"))
// 			//fmt.Println(peerPID)
// 			node.ActorEngine.Subscribe(c.PID())
// 		case *msg.Message:
// 			fmt.Println("actor worker received event ", m.Data)

// 			var workData map[string]string
// 			err := json.Unmarshal([]byte(m.Data), &workData)
// 			if err != nil {
// 				logrus.Errorf("Error parsing work data: %v", err)
// 				return
// 			}
// 			switch workData["request"] {
// 			case "web":
// 				depth, err := strconv.Atoi(workData["depth"])
// 				if err != nil {
// 					logrus.Errorf("Error converting depth to int: %v", err)
// 					return
// 				}
// 				webData, err := scraper.ScrapeWebData([]string{workData["url"]}, depth)
// 				if err != nil {
// 					logrus.Errorf("%v", err)
// 					return
// 				}
// 				collectedData := llmbridge.SanitizeResponse(webData)
// 				jsonData, _ := json.Marshal(collectedData)
// 				workerStatusCh <- jsonData
// 			case "twitter":
// 				count, err := strconv.Atoi(workData["count"])
// 				if err != nil {
// 					logrus.Errorf("Error converting count to int: %v", err)
// 					return
// 				}
// 				tweets, err := twitter.ScrapeTweetsByQuery(workData["query"], count)
// 				if err != nil {
// 					logrus.Errorf("%v", err)
// 					return
// 				}
// 				collectedData := llmbridge.ConcatenateTweets(tweets)
// 				logrus.Info("Actor worker stopped")

// 				jsonData, _ := json.Marshal(collectedData)
// 				workerStatusCh <- jsonData
// 			}

// 			wg.Done()
// 		}
// 	}, "peer_worker")
// 	// go MonitorWorkers(context.Background(), node)
// }

// func SendWork(node *masa.OracleNode, data []byte) {
// 	node.ActorEngine.BroadcastEvent(&msg.Message{Data: string(data)})
// }

// MonitorWorkers monitors worker data by subscribing to the completed work topic,
// computing a CID for each received data, and writing the data to the database.
//
// Parameters:
//   - ctx: The context for the monitoring operation.
//   - node: A pointer to the OracleNode instance.
//
// The function uses a ticker to periodically log a debug message every 60 seconds.
// It subscribes to the completed work topic using the PubSubManager and handles the received data.
// For each received data, it computes a CID using the computeCid function, logs the CID,
// marshals the data to JSON, and writes it to the database using the WriteData function.
// The monitoring continues until the context is done.
func MonitorWorkers(ctx context.Context, node *masa.OracleNode) {

	// new below
	var err error

	node.ActorRemote.Register("peer", actor.PropsFromProducer(NewWorker()))

	props := actor.PropsFromProducer(NewWorker())
	pid := node.ActorEngine.Spawn(props)
	message := &messages.Work{Data: "hi-1", Sender: pid}
	fmt.Println("debug ", message)

	spawnResponse, err := node.ActorRemote.SpawnNamed("192.168.4.164:4001", "worker", "peer", time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		spawnedPID := spawnResponse.Pid
		client := node.ActorEngine.Spawn(props)
		node.ActorEngine.Send(spawnedPID, &messages.Connect{
			Sender: client,
		})
		for i := 0; i < 10; i++ {
			node.ActorEngine.Send(spawnedPID, message)
		}
	}

	peers := node.Host.Network().Peers()
	for _, peer := range peers {
		conns := node.Host.Network().ConnsToPeer(peer)
		for _, conn := range conns {
			addr := conn.RemoteMultiaddr()
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			logrus.Info(fmt.Sprintf("%s:4001", ipAddr))
			//spawnResponse, err := node.ActorRemote.SpawnNamed(fmt.Sprintf("%s:4001", ipAddr), "worker", "peer", time.Second)
			//if err != nil {
			//	logrus.Errorf("spawn error %v", err)
			//} else {
			//	spawnedPID := spawnResponse.Pid
			//	client := node.ActorEngine.Spawn(props)
			//	node.ActorEngine.Send(spawnedPID, &messages.Connect{
			//		Sender: client,
			//	})
			//}
		}
	}

	syncInterval := time.Second * 60
	workerStatusHandler := &pubsub.WorkerStatusHandler{WorkerStatusCh: workerStatusCh}
	err = node.PubSubManager.Subscribe(config.TopicWithVersion(config.CompletedWorkTopic), workerStatusHandler)
	if err != nil {
		logrus.Errorf("%v", err)
	}

	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			logrus.Debug("tick")
		case data := <-workerStatusCh:
			key, _ := computeCid(string(data))
			val := db.ReadData(node, key)
			// double spend check
			if val == nil {
				go db.WriteData(node, key, data)
				nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
				sharedData := db.SharedData{}
				nodeVal := db.ReadData(node, nodeData.PeerId.String())
				_ = json.Unmarshal(nodeVal, &sharedData)
				bytesScraped, _ := strconv.Atoi(fmt.Sprintf("%v", sharedData["bytesScraped"]))
				nodeData.BytesScraped += bytesScraped
				nodeData.BytesScraped += len(data)
				err = node.NodeTracker.AddOrUpdateNodeData(nodeData, true)
				if err != nil {
					logrus.Error(err)
				}
				jsonData, _ := json.Marshal(nodeData)
				go db.WriteData(node, node.Host.ID().String(), jsonData)
			}
		case <-ctx.Done():
			return
		}
	}
}
