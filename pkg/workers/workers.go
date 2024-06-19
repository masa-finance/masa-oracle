package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"

	"github.com/multiformats/go-multiaddr"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/ipfs/go-cid"

	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

type WorkerType string

const (
	Discord                WorkerType = "discord"
	DiscordProfile         WorkerType = "discord-profile"
	DiscordChannelMessages WorkerType = "discord-channel-messages"
	DiscordSentiment       WorkerType = "discord-sentiment"
	DiscordGuildChannels   WorkerType = "discord-guild-channels"
	DiscordUserGuilds      WorkerType = "discord-user-guilds"
	LLMChat                WorkerType = "llm-chat"
	Twitter                WorkerType = "twitter"
	TwitterFollowers       WorkerType = "twitter-followers"
	TwitterProfile         WorkerType = "twitter-profile"
	TwitterSentiment       WorkerType = "twitter-sentiment"
	TwitterTrends          WorkerType = "twitter-trends"
	Web                    WorkerType = "web"
	WebSentiment           WorkerType = "web-sentiment"
	Test                   WorkerType = "test"
)

var WORKER = struct {
	Discord, DiscordProfile, DiscordChannelMessages, DiscordSentiment, DiscordGuildChannels, DiscordUserGuilds, LLMChat, Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends, Web, WebSentiment, Test WorkerType
}{
	Discord:                Discord,
	DiscordProfile:         DiscordProfile,
	DiscordChannelMessages: DiscordChannelMessages,
	DiscordSentiment:       DiscordSentiment,
	DiscordGuildChannels:   DiscordGuildChannels,
	DiscordUserGuilds:      DiscordUserGuilds,
	LLMChat:                LLMChat,
	Twitter:                Twitter,
	TwitterFollowers:       TwitterFollowers,
	TwitterProfile:         TwitterProfile,
	TwitterSentiment:       TwitterSentiment,
	TwitterTrends:          TwitterTrends,
	Web:                    Web,
	WebSentiment:           WebSentiment,
	Test:                   Test,
}

var (
	clients        = actor.NewPIDSet()
	workerStatusCh = make(chan *pubsub2.Message)
	workerDoneCh   = make(chan *pubsub2.Message)
)

type CID struct {
	Duration  float64 `json:"duration"`
	RecordId  string  `json:"cid"`
	Timestamp int64   `json:"timestamp"`
}

type Record struct {
	PeerId string `json:"peerid"`
	CIDs   []CID  `json:"cids"`
}

type ChanResponse struct {
	Response  map[string]interface{}
	ChannelId string
}

type Worker struct {
	Node *masa.OracleNode
}

// NewWorker creates a new instance of the Worker actor.
// It implements the actor.Receiver interface, allowing it to receive and handle messages.
//
// Returns:
//   - An instance of the Worker struct that implements the actor.Receiver interface.
func NewWorker(node *masa.OracleNode) actor.Producer {
	return func() actor.Actor {
		return &Worker{Node: node}
	}
}

// Receive is the message handling method for the Worker actor.
// It receives messages through the actor context and processes them based on their type.
func (a *Worker) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *messages.Connect:
		a.HandleConnect(ctx, m)
	case *actor.Started:
		a.HandleLog(ctx, "[+] Actor started")
	case *actor.Stopping:
		a.HandleLog(ctx, "[+] Actor stopping")
	case *actor.Stopped:
		a.HandleLog(ctx, "[+] Actor stopped")
	case *messages.Work:
		if a.Node.IsActor() {
			a.HandleWork(ctx, m, a.Node)
		}
	case *messages.Response:
		msg := &pubsub2.Message{}
		err := json.Unmarshal([]byte(m.Value), msg)
		if err != nil {
			msg, err = getResponseMessage(m)
			if err != nil {
				logrus.Errorf("Error getting response message: %v", err)
				return
			}
		}
		workerDoneCh <- msg
		ctx.Poison(ctx.Self())
	default:
		logrus.Warningf("[+] Received unknown message: %T", m)
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
		if bn == "" {
			return false
		}
		bootNodeAddr := strings.Split(bn, "/")[2]
		if bootNodeAddr == ipAddr {
			return true
		}
	}
	return false
}

// updateRecords updates the records for a given node and key with the provided data.
//
// Parameters:
//   - node: A pointer to the OracleNode instance whose records need to be updated.
//   - data: The data to be written for the specified key.
//   - key: The key under which the data should be stored.
func updateRecords(node *masa.OracleNode, workEvent db.WorkEvent) {

	exists := db.ReadData(node, workEvent.CID)
	if exists == nil {
		_ = db.WriteData(node, workEvent.CID, workEvent.Payload)
	}

	var nodeData pubsub.NodeData
	nodeDataBytes, err := db.GetCache(context.Background(), workEvent.PeerId)
	if err != nil {
		nodeData = *node.NodeTracker.GetNodeData(workEvent.PeerId)
	} else {
		err = json.Unmarshal(nodeDataBytes, &nodeData)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	nodeData.BytesScraped = nodeData.BytesScraped + len(workEvent.Payload)

	newCID := CID{
		RecordId:  workEvent.CID,
		Duration:  workEvent.Duration,
		Timestamp: time.Now().Unix(),
	}

	records := nodeData.Records

	if records == nil {
		recordsSlice, ok := records.([]CID)
		if !ok {
			recordsSlice = []CID{}
		}
		recordsSlice = append(recordsSlice, newCID)
		nodeData.Records = recordsSlice
		err = node.NodeTracker.AddOrUpdateNodeData(&nodeData, true)
		if err != nil {
			logrus.Error(err)
			return
		}
	} else {
		if exists == nil {
			records = append(nodeData.Records.([]interface{}), newCID)
			nodeData.Records = records
			err = node.NodeTracker.AddOrUpdateNodeData(&nodeData, true)
			if err != nil {
				logrus.Error(err)
				return
			}
		}
	}

	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		logrus.Error(err)
		return
	}
	_ = db.WriteData(node, workEvent.PeerId, jsonData)
	logrus.Infof("[+] Updated records key %s for node %s", workEvent.CID, workEvent.PeerId)
}

// getResponseMessage converts a messages.Response object into a pubsub2.Message object.
// It unmarshals the JSON-encoded response value into a map and then constructs a new pubsub2.Message
// using the extracted data.
//
// Parameters:
//   - response: A pointer to a messages.Response object containing the JSON-encoded response data.
//
// Returns:
//   - A pointer to a pubsub2.Message object constructed from the response data.
//   - An error if there is an issue with unmarshalling the response data.
func getResponseMessage(response *messages.Response) (*pubsub2.Message, error) {
	responseData := map[string]interface{}{}

	err := json.Unmarshal([]byte(response.Value), &responseData)
	if err != nil {
		return nil, err
	}
	msg := &pubsub2.Message{
		ID:            responseData["ID"].(string),
		ReceivedFrom:  peer.ID(responseData["ReceivedFrom"].(string)),
		ValidatorData: responseData["ValidatorData"],
		Local:         responseData["Local"].(bool),
	}
	return msg, nil
}

// SendWork is a function that sends work to a node. It takes two parameters:
// node: A pointer to a masa.OracleNode object. This is the node to which the work will be sent.
// m: A pointer to a pubsub2.Message object. This is the message that contains the work to be sent.
func SendWork(node *masa.OracleNode, m *pubsub2.Message) {
	var wg sync.WaitGroup
	props := actor.PropsFromProducer(NewWorker(node))
	pid := node.ActorEngine.Spawn(props)
	message := &messages.Work{Data: string(m.Data), Sender: pid, Id: m.ReceivedFrom.String()}
	// local
	if node.IsStaked && node.IsActor() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			future := node.ActorEngine.RequestFuture(pid, message, 30*time.Second)
			result, err := future.Result()
			if err != nil {
				logrus.Debugf("Error receiving response: %v", err)
				return
			}
			response := result.(*messages.Response)
			msg := &pubsub2.Message{}
			err = json.Unmarshal([]byte(response.Value), msg)
			if err != nil {
				msg, err = getResponseMessage(result.(*messages.Response))
				if err != nil {
					logrus.Debugf("Error getting response message: %v", err)
					return
				}
			}
			workerDoneCh <- msg
		}()
	}
	// remote
	peers := node.NodeTracker.GetAllNodeData()
	for _, p := range peers {
		for _, addr := range p.Multiaddrs {
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			logrus.Infof("[+] Worker Address: %s", ipAddr)
			if !isBootnode(ipAddr) && p.IsTwitterScraper || p.IsWebScraper || p.IsDiscordScraper {
				wg.Add(1)
				go func() {
					defer wg.Done()
					spawned, err := node.ActorRemote.SpawnNamed(fmt.Sprintf("%s:4001", ipAddr), "worker", "peer", -1)
					if err != nil {
						logrus.Debugf("Spawned error %v", err)
					} else {
						spawnedPID := spawned.Pid
						client := node.ActorEngine.Spawn(props)
						node.ActorEngine.Send(spawnedPID, &messages.Connect{
							Sender: client,
						})
						future := node.ActorEngine.RequestFuture(spawnedPID, message, 30*time.Second)
						result, err := future.Result()
						if err != nil {
							logrus.Debugf("Error receiving response: %v", err)
							return
						}
						response := result.(*messages.Response)
						node.ActorEngine.Send(spawnedPID, response)
					}
				}()
			}
		}
	}
	wg.Wait()
}

// SubscribeToWorkers subscribes the given OracleNode to worker events.
//
// Parameters:
//   - node: A pointer to the OracleNode instance that will be subscribed to worker events.
//
// The function initializes the WorkerEventTracker for the node and adds a subscription
// to the worker topic using the PubSubManager. If an error occurs during the subscription,
// it logs the error.
func SubscribeToWorkers(node *masa.OracleNode) {
	node.WorkerTracker = &pubsub.WorkerEventTracker{WorkerStatusCh: workerStatusCh}
	err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.WorkerTopic), node.WorkerTracker, true)
	if err != nil {
		logrus.Errorf("Subscribe error %v", err)
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
	node.ActorRemote.Register("peer", actor.PropsFromProducer(NewWorker(node)))

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	rcm := pubsub.GetResponseChannelMap()
	var startTime time.Time

	for {
		select {
		case <-ticker.C:
			logrus.Debug("tick")
		case work := <-node.WorkerTracker.WorkerStatusCh:
			logrus.Info("[+] Sending work to network")
			var workData map[string]string
			err := json.Unmarshal(work.Data, &workData)
			if err != nil {
				logrus.Error(err)
				continue
			}
			startTime = time.Now()
			go SendWork(node, work)
		case data := <-workerDoneCh:
			validatorDataMap, ok := data.ValidatorData.(map[string]interface{})
			if !ok {
				logrus.Errorf("Error asserting type: %v", ok)
				continue
			}

			if ch, ok := rcm.Get(validatorDataMap["ChannelId"].(string)); ok {
				validatorData, err := json.Marshal(validatorDataMap["Response"])
				if err != nil {
					logrus.Errorf("Error marshalling data.ValidatorData: %v", err)
					continue
				}
				ch <- validatorData
				close(ch)
			} else {
				logrus.Debugf("Error processing data.ValidatorData: %v", data.ValidatorData)
			}

			if validatorDataMap, ok := data.ValidatorData.(map[string]interface{}); ok {
				if response, ok := validatorDataMap["Response"].(map[string]interface{}); ok {
					if _, ok := response["error"].(string); ok {
						logrus.Infof("[+] Work failed %s", response["error"])
					} else if work, ok := response["data"].(string); ok {
						key, _ := computeCid(work)
						logrus.Infof("[+] Work done %s", key)

						endTime := time.Now()
						duration := endTime.Sub(startTime)

						workEvent := db.WorkEvent{
							CID:       key,
							PeerId:    data.ID,
							Payload:   []byte(work),
							Duration:  duration.Seconds(),
							Timestamp: time.Now().Unix(),
						}

						updateRecords(node, workEvent)
					} else if w, ok := response["data"].(map[string]interface{}); ok {
						work, err := json.Marshal(w)
						if err != nil {
							logrus.Errorf("Error marshalling data.ValidatorData: %v", err)
							continue
						}
						key, _ := computeCid(string(work))
						logrus.Infof("[+] Work done %s", key)

						endTime := time.Now()
						duration := endTime.Sub(startTime)

						workEvent := db.WorkEvent{
							CID:       key,
							PeerId:    data.ID,
							Payload:   work,
							Duration:  duration.Seconds(),
							Timestamp: time.Now().Unix(),
						}

						updateRecords(node, workEvent)
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
