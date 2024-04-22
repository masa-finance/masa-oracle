package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/scraper"

	"github.com/ipfs/go-cid"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/llmbridge"

	"github.com/anthdm/hollywood/actor"
	masa "github.com/masa-finance/masa-oracle/pkg"
	msg "github.com/masa-finance/masa-oracle/pkg/proto/msg"
	"github.com/masa-finance/masa-oracle/pkg/twitter"
	"github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

var workerStatusCh = make(chan []byte)

type Worker struct{}

// NewWorker creates a new instance of the Worker actor.
// It implements the actor.Receiver interface, allowing it to receive and handle messages.
//
// Returns:
//   - An instance of the Worker struct that implements the actor.Receiver interface.
func NewWorker() actor.Receiver {
	return &Worker{}
}

// Receive is the message handling method for the Worker actor.
// It receives messages through the actor context and processes them based on their type.
func (w *Worker) Receive(ctx *actor.Context) {
	switch m := ctx.Message().(type) {
	case *msg.Message:
		logrus.Info("Actor worker initialized")
		var workData map[string]string
		err := json.Unmarshal([]byte(m.Data), &workData)
		if err != nil {
			logrus.Errorf("Error parsing work data: %v", err)
			return
		}
		switch workData["request"] {
		case "web":
			depth, err := strconv.Atoi(workData["depth"])
			if err != nil {
				logrus.Errorf("Error converting depth to int: %v", err)
				return
			}
			webData, err := scraper.ScrapeWebData([]string{workData["url"]}, depth)
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			collectedData := llmbridge.SanitizeResponse(webData)
			jsonData, _ := json.Marshal(collectedData)
			workerStatusCh <- jsonData
		case "twitter":
			count, err := strconv.Atoi(workData["count"])
			if err != nil {
				logrus.Errorf("Error converting count to int: %v", err)
				return
			}
			tweets, err := twitter.ScrapeTweetsByQuery(workData["query"], count)
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			collectedData := llmbridge.ConcatenateTweets(tweets)
			logrus.Info("Actor worker stopped")

			jsonData, _ := json.Marshal(collectedData)
			workerStatusCh <- jsonData
		}
	default:
		break
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
//
// Examples:
//
//	d, _ := json.Marshal(map[string]string{"request": "web", "url": "https://en.wikipedia.org/wiki/Maize", "depth": "2"})
//	d, _ := json.Marshal(map[string]string{"request": "twitter", "query": "$MASA", "count": "5"})
//	go workers.SendWorkToPeers(node, d)
func SendWorkToPeers(node *masa.OracleNode, data []byte) {
	peers := node.Host.Network().Peers()
	for _, peer := range peers {
		conns := node.Host.Network().ConnsToPeer(peer)
		for _, conn := range conns {
			addr := conn.RemoteMultiaddr()
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			peerPID := actor.NewPID(fmt.Sprintf("%s:4001", ipAddr), fmt.Sprintf("%s/%s", "peer_worker", "peer"))
			node.ActorEngine.Subscribe(peerPID)
		}
	}
	node.ActorEngine.BroadcastEvent(&msg.Message{Data: string(data)})
	// ctx.Send(ctx.PID(), &msg.Message{Data: string(jsonData)})
	go monitorWorkerData(context.Background(), node)
}

// monitorWorkerData monitors worker data by subscribing to the completed work topic,
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
func monitorWorkerData(ctx context.Context, node *masa.OracleNode) {
	syncInterval := time.Second * 60
	workerStatusHandler := &pubsub.WorkerStatusHandler{WorkerStatusCh: workerStatusCh}
	err := node.PubSubManager.Subscribe(config.TopicWithVersion(config.CompletedWorkTopic), workerStatusHandler)
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
				nodeData.BytesScraped += len(data)
				err = node.NodeTracker.AddOrUpdateNodeData(nodeData, true)
				if err != nil {
					logrus.Error(err)
				}
				jsonData, _ := json.Marshal(nodeData)
				go db.WriteData(node, node.Host.ID().String(), jsonData)
				err = node.PubSubManager.Publish(config.TopicWithVersion(config.NodeDataSyncProtocol), jsonData)
				if err != nil {
					logrus.Errorf("Error publishing node data: %v", err)
				}
			}

			// @todo add list of keys to nodeData ie Records: []string ?
		case <-ctx.Done():
			return
		}
	}
}
