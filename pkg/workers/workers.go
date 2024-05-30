package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
	Discord          WorkerType = "discord"
	LLMChat          WorkerType = "llm-chat"
	Twitter          WorkerType = "twitter"
	TwitterFollowers WorkerType = "twitter-followers"
	TwitterProfile   WorkerType = "twitter-profile"
	TwitterSentiment WorkerType = "twitter-sentiment"
	TwitterTrends    WorkerType = "twitter-trends"
	Web              WorkerType = "web"
	WebSentiment     WorkerType = "web-sentiment"
)

var WORKER = struct {
	Discord, LLMChat, Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends, Web, WebSentiment WorkerType
}{
	Discord:          Discord,
	LLMChat:          LLMChat,
	Twitter:          Twitter,
	TwitterFollowers: TwitterFollowers,
	TwitterProfile:   TwitterProfile,
	TwitterSentiment: TwitterSentiment,
	TwitterTrends:    TwitterTrends,
	Web:              Web,
	WebSentiment:     WebSentiment,
}

var (
	clients        = actor.NewPIDSet()
	dataLengthCh   = make(chan int)
	workerStatusCh = make(chan *pubsub2.Message)
	workerDoneCh   = make(chan *pubsub2.Message)
)

type CID struct {
	RecordId  string    `json:"cid"`
	Timestamp time.Time `json:"timestamp"`
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
}

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
		a.HandleConnect(ctx, m)
	case *actor.Started:
		a.HandleLog(ctx, "[+] Actor started")
	case *actor.Stopping:
		a.HandleLog(ctx, "[+] Actor stopping")
	case *actor.Stopped:
		a.HandleLog(ctx, "[+] Actor stopped")
	case *messages.Work:
		a.HandleWork(ctx, m)
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
			return true
		}
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
func updateParticipation(node *masa.OracleNode, totalBytes int, peerId string) {
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
	_ = db.WriteData(node, peerId, jsonData)
}

// updateRecords updates the records for a given node and key with the provided data.
//
// Parameters:
//   - node: A pointer to the OracleNode instance whose records need to be updated.
//   - data: The data to be written for the specified key.
//   - key: The key under which the data should be stored.
func updateRecords(node *masa.OracleNode, data []byte, key string, peerId string) {
	_ = db.WriteData(node, key, data)

	var nodeData pubsub.NodeData
	nodeDataBytes, err := db.GetCache(context.Background(), peerId)
	if err != nil {
		nd := node.NodeTracker.GetNodeData(peerId)
		nodeData = *nd
	} else {
		err = json.Unmarshal(nodeDataBytes, &nodeData)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	newCID := CID{
		RecordId:  key,
		Timestamp: time.Now(),
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
		records = append(nodeData.Records.([]interface{}), newCID)
		nodeData.Records = records
		err = node.NodeTracker.AddOrUpdateNodeData(&nodeData, true)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		logrus.Error(err)
		return
	}
	_ = db.WriteData(node, peerId, jsonData)
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

func SendWork(node *masa.OracleNode, m *pubsub2.Message) {
	var wg sync.WaitGroup
	props := actor.PropsFromProducer(NewWorker())
	pid := node.ActorEngine.Spawn(props)
	message := &messages.Work{Data: string(m.Data), Sender: pid, Id: m.ReceivedFrom.String()}
	if node.IsStaked {
		wg.Add(1)
		go func() {
			defer wg.Done()
			future := node.ActorEngine.RequestFuture(pid, message, 30*time.Second)
			result, err := future.Result()
			if err != nil {
				logrus.Errorf("Error receiving response: %v", err)
				return
			}
			response := result.(*messages.Response)
			msg := &pubsub2.Message{}
			err = json.Unmarshal([]byte(response.Value), msg)
			if err != nil {
				msg, err = getResponseMessage(result.(*messages.Response))
				if err != nil {
					logrus.Errorf("Error getting response message: %v", err)
					return
				}
			}
			workerDoneCh <- msg
		}()
	}
	peers := node.Host.Network().Peers()
	for _, p := range peers {
		conns := node.Host.Network().ConnsToPeer(p)
		for _, conn := range conns {
			addr := conn.RemoteMultiaddr()
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			if !isBootnode(ipAddr) {
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
							logrus.Errorf("Error receiving response: %v", err)
							return
						}
						response := result.(*messages.Response)
						msg := &pubsub2.Message{}
						err = json.Unmarshal([]byte(response.Value), msg)
						if err != nil {
							msg, err = getResponseMessage(result.(*messages.Response))
							if err != nil {
								logrus.Errorf("Error getting response message: %v", err)
								return
							}
						}
						workerDoneCh <- msg
					}
				}()
			}
		}
	}
	wg.Wait()
}

// SubscribeToWorkers subscribes the given OracleNode to worker events.
// @note need to add this even if not participating to send messages ....
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
	node.ActorRemote.Register("peer", actor.PropsFromProducer(NewWorker()))

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	rcm := pubsub.GetResponseChannelMap()
	var validatorData []byte

	for {
		select {
		case <-ticker.C:
			logrus.Debug("tick")
		case totalBytes := <-dataLengthCh:
			go updateParticipation(node, totalBytes, node.Host.ID().String())
		case work := <-node.WorkerTracker.WorkerStatusCh:
			logrus.Info("[+] Sending work to network")
			go SendWork(node, work)
		case data := <-workerDoneCh:
			if _, ok := data.ValidatorData.([]byte); ok {
				validatorData = data.ValidatorData.([]byte)
			} else if _, ok := data.ValidatorData.(string); ok {
				validatorData = []byte(data.ValidatorData.(string))
			} else {
				switch reflect.TypeOf(data.ValidatorData).Kind() {
				case reflect.Map:
					validatorDataMap, ok := data.ValidatorData.(map[string]interface{})
					if !ok {
						logrus.Errorf("Error asserting type: %v", ok)
					}

					if ch, ok := rcm.Get(validatorDataMap["ChannelId"].(string)); ok {
						validatorData, err := json.Marshal(validatorDataMap["Response"])
						if err != nil {
							logrus.Errorf("Error marshalling data.ValidatorData: %v", err)
						}
						ch <- validatorData
						close(ch)
					}
				default:
					logrus.Errorf("Error processing data.ValidatorData: %v", data.ValidatorData)
				}
			}
			if validatorDataMap, ok := data.ValidatorData.(map[string]interface{}); ok {
				if response, ok := validatorDataMap["Response"].(map[string]interface{}); ok {
					if errorMessage, ok := response["error"].(string); ok {
						if errorMessage == "there was an error authenticating with your Twitter credentials" {
							logrus.Infof("[+] Work failed %s", errorMessage)
						}
					} else {
						key, _ := computeCid(string(validatorData))
						logrus.Infof("[+] Work done %s", key)
						updateRecords(node, validatorData, key, data.ID)
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
