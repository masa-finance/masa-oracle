package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/masa-finance/masa-oracle/pkg/db"

	"github.com/masa-finance/masa-oracle/pkg/llmbridge"
	"github.com/masa-finance/masa-oracle/pkg/scraper"
	"github.com/masa-finance/masa-oracle/pkg/twitter"

	"github.com/multiformats/go-multiaddr"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"

	actor "github.com/asynkron/protoactor-go/actor"
	messages "github.com/masa-finance/masa-oracle/pkg/workers/messages"

	"github.com/ipfs/go-cid"

	masa "github.com/masa-finance/masa-oracle/pkg"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

var (
	clients        = actor.NewPIDSet()
	dataLengthCh   = make(chan int)
	workerStatusCh = make(chan []byte)
	workerDoneCh   = make(chan []byte)
)

// @todo to send to api gateway

type CID struct {
	RecordId  string    `json:"recordid"`
	Timestamp time.Time `json:"timestamp"`
}

type Record struct {
	PeerId string `json:"peerid"`
	CIDs   []CID  `json:"cids"`
}

type OracleData struct {
	Id        string `json:"id"`
	PeerId    string `json:"peer_id"`
	Request   string `json:"request"`
	Domain    string `json:"domain"`
	ModelType string `json:"model_type"`
	ModelName string `json:"model_name"`
	Steps     []struct {
		Idx               int    `json:"idx"`
		StructuredPrompt  string `json:"structured_prompt,omitempty"`
		Timestamp         string `json:"timestamp"`
		UserPrompt        string `json:"user_prompt,omitempty"`
		RawContent        string `json:"raw_content,omitempty"`
		StructuredContent string `json:"structured_content,omitempty"`
	} `json:"steps"`
}

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
	switch m := ctx.Message().(type) {
	case *messages.Connect:
		logrus.Infof("Worker %v connected", m.Sender)
		clients.Add(m.Sender)
	case *actor.Started:
		logrus.Info("actor started")
	case *actor.Stopping:
		logrus.Info("actor stopping")
	case *actor.Stopped:
		logrus.Info("actor stopped")
	case *messages.Work:
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
			if webData != "" {
				collectedData := llmbridge.SanitizeResponse(webData)
				jsonData, _ := json.Marshal(collectedData)
				workerDoneCh <- jsonData
			}
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
			if tweets != nil {
				collectedData := llmbridge.ConcatenateTweets(tweets)
				jsonData, _ := json.Marshal(collectedData)
				workerDoneCh <- jsonData
			}
		case "twitter-sentiment":
			count, err := strconv.Atoi(workData["count"])
			if err != nil {
				logrus.Errorf("Error converting count to int: %v", err)
				return
			}
			_, sentimentSummary, _ := twitter.ScrapeTweetsForSentiment(workData["query"], count, workData["model"])
			if sentimentSummary != "" {
				jsonData, _ := json.Marshal(sentimentSummary)
				workerDoneCh <- jsonData
			}
		case "web-sentiment":
			depth, err := strconv.Atoi(workData["depth"])
			if err != nil {
				logrus.Errorf("Error converting depth to int: %v", err)
				return
			}
			_, sentimentSummary, _ := scraper.ScrapeWebDataForSentiment([]string{workData["url"]}, depth, workData["model"])
			if sentimentSummary != "" {
				jsonData, _ := json.Marshal(sentimentSummary)
				workerDoneCh <- jsonData
			}
		}
		ctx.Poison(ctx.Self())
	default:
		logrus.Warningf("Received unknown message: %T", m)
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

// isBootnode checks if the given IP address belongs to a bootnode.
// It takes an IP address as a string and returns a boolean value.
// If the IP address is found in the list of bootnodes, it returns true.
// Otherwise, it returns false.
func isBootnode(ipAddr string) bool {
	for _, bn := range config.GetInstance().Bootnodes {
		bootNodeAddr := strings.Split(bn, "/")[2]
		if bootNodeAddr == ipAddr {
			return true
		}
	}
	return false
}

// updateParticipation updates the participation metrics for a given node based on the total bytes processed.
//
// Parameters:
//   - node: A pointer to the OracleNode instance whose participation metrics need to be updated.
//   - totalBytes: The total number of bytes processed by the node which will be added to its current metrics.
func updateParticipation(node *masa.OracleNode, totalBytes int) {
	if totalBytes == 0 {
		return
	}
	nodeData := node.NodeTracker.GetNodeData(node.Host.ID().String())
	sharedData := db.SharedData{}
	nodeVal := db.ReadData(node, nodeData.PeerId.String())
	_ = json.Unmarshal(nodeVal, &sharedData)
	bytesScraped, _ := strconv.Atoi(fmt.Sprintf("%v", sharedData["bytesScraped"]))
	nodeData.BytesScraped += bytesScraped + totalBytes
	err := node.NodeTracker.AddOrUpdateNodeData(nodeData, true)
	if err != nil {
		logrus.Error(err)
	}
	jsonData, _ := json.Marshal(nodeData)
	go db.WriteData(node, node.Host.ID().String(), jsonData)
}

// SendWork is responsible for handling work messages and processing them based on the request type.
// It supports the following request types:
// - "web": Scrapes web data from the specified URL with the given depth.
// - "twitter": Scrapes tweets based on the provided query and count.
//
// The Worker actor receives messages through its Receive method, which is called by the actor system when a message is sent to the actor.
// It handles the following message types:
//   - *messages.Connect: Indicates that a client has connected to the worker. The client's sender information is added to the clients set.
//   - *actor.Started: Indicates that the actor has started. It prints a debug message.
//   - *messages.Work: Contains the work data to be processed. The work data is parsed based on the request type, and the corresponding scraping function is called.
//     The scraped data is then sent to the workerStatusCh channel for further processing.
//
// The Worker actor is responsible for the NewWorker function, which returns an actor.Producer that can be used to spawn new instances of the Worker actor.
// @note we can use the WorkerTopic to SendWork anywhere in the service
// Usage with the Worker Gossip Topic
//
//	if err := node.PubSubManager.Publish(config.TopicWithVersion(config.WorkerTopic), data); err != nil {
//	    logrus.Errorf("%v", err)
//	}
func SendWork(node *masa.OracleNode, data []byte) {
	props := actor.PropsFromProducer(NewWorker())
	pid := node.ActorEngine.Spawn(props)
	id := uuid.New().String()
	message := &messages.Work{Data: string(data), Sender: pid, Id: id}
	node.ActorEngine.Send(pid, message)
	peers := node.Host.Network().Peers()
	for _, peer := range peers {
		conns := node.Host.Network().ConnsToPeer(peer)
		for _, conn := range conns {
			addr := conn.RemoteMultiaddr()
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			if !isBootnode(ipAddr) {
				spawned, err := node.ActorRemote.SpawnNamed(fmt.Sprintf("%s:4001", ipAddr), "worker", "peer", -1)
				if err != nil {
					logrus.Debugf("Spawned error %v", err)
				} else {
					spawnedPID := spawned.Pid
					client := node.ActorEngine.Spawn(props)
					node.ActorEngine.Send(spawnedPID, &messages.Connect{
						Sender: client,
					})
					node.ActorEngine.Send(spawnedPID, message)
				}
			}
		}
	}
}

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
	// Register self as a remote node for the network
	node.ActorRemote.Register("peer", actor.PropsFromProducer(NewWorker()))

	var err error
	workerEventTracker := &pubsub.WorkerEventTracker{WorkerStatusCh: workerStatusCh}
	err = node.PubSubManager.Subscribe(config.TopicWithVersion(config.WorkerTopic), workerEventTracker)
	if err != nil {
		logrus.Errorf("Subscribe error %v", err)
	}

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			logrus.Debug("tick")
		case totalBytes := <-dataLengthCh:
			go updateParticipation(node, totalBytes)
		case work := <-workerStatusCh:
			logrus.Infof("work ==> %s", string(work))
			// If participating node
			if node.IsTwitterScraper || node.IsWebScraper {
				pid := node.ActorEngine.Spawn(actor.PropsFromProducer(NewWorker()))
				id := uuid.New().String()
				message := &messages.Work{Data: string(work), Sender: pid, Id: id}
				logrus.Infof("work ==> %s %s", string(work), pid)
				node.ActorEngine.Send(pid, message)
			} else {
				// Send to participating nodes
				go SendWork(node, work)
			}
		case data := <-workerDoneCh:
			key, _ := computeCid(string(data))
			logrus.Infof("work done ==> %s %s", key, string(data))
			val := db.ReadData(node, key)
			if val == nil {
				// store raw data
				_ = db.WriteData(node, key, data)

				// store record cids for the peer
				existingRecord := Record{}
				existingData := db.ReadData(node, key)
				if existingData != nil {
					_ = json.Unmarshal(existingData, &existingRecord)
				}

				newCID := CID{
					RecordId:  key,
					Timestamp: time.Now(),
				}

				record := Record{
					PeerId: node.Host.ID().String(),
					CIDs:   append(existingRecord.CIDs, newCID),
				}
				jsonData, _ := json.Marshal(record)
				_ = db.WriteData(node, key, jsonData)
			}
		case <-ctx.Done():
			return
		}
	}
}
