package workers

import (
	"context"
	"encoding/json"
	"fmt"
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
	Discord                 WorkerType = "discord"
	DiscordProfile          WorkerType = "discord-profile"
	DiscordChannelMessages  WorkerType = "discord-channel-messages"
	DiscordSentiment        WorkerType = "discord-sentiment"
	TelegramSentiment       WorkerType = "telegram-sentiment"
	TelegramChannelMessages WorkerType = "telegram-channel-messages"
	DiscordGuildChannels    WorkerType = "discord-guild-channels"
	DiscordUserGuilds       WorkerType = "discord-user-guilds"
	LLMChat                 WorkerType = "llm-chat"
	Twitter                 WorkerType = "twitter"
	TwitterFollowers        WorkerType = "twitter-followers"
	TwitterProfile          WorkerType = "twitter-profile"
	TwitterSentiment        WorkerType = "twitter-sentiment"
	TwitterTrends           WorkerType = "twitter-trends"
	Web                     WorkerType = "web"
	WebSentiment            WorkerType = "web-sentiment"
	Test                    WorkerType = "test"
)

var WORKER = struct {
	Discord, DiscordProfile, DiscordChannelMessages, DiscordSentiment, TelegramSentiment, TelegramChannelMessages, DiscordGuildChannels, DiscordUserGuilds, LLMChat, Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends, Web, WebSentiment, Test WorkerType
}{
	Discord:                 Discord,
	DiscordProfile:          DiscordProfile,
	DiscordChannelMessages:  DiscordChannelMessages,
	DiscordSentiment:        DiscordSentiment,
	TelegramSentiment:       TelegramSentiment,
	TelegramChannelMessages: TelegramChannelMessages,
	DiscordGuildChannels:    DiscordGuildChannels,
	DiscordUserGuilds:       DiscordUserGuilds,
	LLMChat:                 LLMChat,
	Twitter:                 Twitter,
	TwitterFollowers:        TwitterFollowers,
	TwitterProfile:          TwitterProfile,
	TwitterSentiment:        TwitterSentiment,
	TwitterTrends:           TwitterTrends,
	Web:                     Web,
	WebSentiment:            WebSentiment,
	Test:                    Test,
}

var (
	clients        = actor.NewPIDSet()
	workerStatusCh = make(chan *pubsub2.Message)
	workerDoneCh   = make(chan *pubsub2.Message)
)

// WorkerTypeToCategory maps WorkerType to WorkerCategory
func WorkerTypeToCategory(wt WorkerType) pubsub.WorkerCategory {
	switch wt {
	case Discord, DiscordProfile, DiscordChannelMessages, DiscordSentiment, DiscordGuildChannels, DiscordUserGuilds:
		return pubsub.CategoryDiscord
	case TelegramSentiment, TelegramChannelMessages:
		return pubsub.CategoryTelegram
	case Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends:
		return pubsub.CategoryTwitter
	case Web, WebSentiment:
		return pubsub.CategoryWeb
	default:
		return -1 // Invalid category
	}
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
		if a.Node.IsWorker() {
			a.HandleLog(ctx, "Actor started")
		}
	case *actor.Stopping:
		if a.Node.IsWorker() {
			a.HandleLog(ctx, "Actor stopping")
		}
	case *actor.Stopped:
		if a.Node.IsWorker() {
			a.HandleLog(ctx, "Actor stopped")
		}
	case *messages.Work:
		if a.Node.IsWorker() {
			logrus.Infof("Received Work")
			a.HandleWork(ctx, m, a.Node)
		}
	case *messages.Response:
		logrus.Infof("Received Response")
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
		logrus.Warningf("Received unknown message in workers: %T, message: %+v", m, m)
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
	message := &messages.Work{Data: string(m.Data), Sender: pid, Id: m.ReceivedFrom.String(), Type: int64(pubsub.CategoryTwitter)}
	n := 0

	responseCollector := make(chan *pubsub2.Message, 100) // Buffered channel to collect responses
	timeout := time.After(8 * time.Second)

	// Local worker
	if node.IsStaked && node.IsWorker() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			future := node.ActorEngine.RequestFuture(pid, message, 60*time.Second) // Increase timeout from 30 to 60 seconds
			result, err := future.Result()
			if err != nil {
				logrus.Errorf("Error receiving response from local worker: %v", err)
				responseCollector <- &pubsub2.Message{
					ValidatorData: map[string]interface{}{"error": err.Error()},
				}
				return
			}
			response := result.(*messages.Response)
			msg := &pubsub2.Message{}
			rErr := json.Unmarshal([]byte(response.Value), msg)
			if rErr != nil {
				gMsg, gErr := getResponseMessage(result.(*messages.Response))
				if gErr != nil {
					logrus.Errorf("Error getting response message: %v", gErr)
					responseCollector <- &pubsub2.Message{
						ValidatorData: map[string]interface{}{"error": gErr.Error()},
					}
					return
				}
				msg = gMsg
			}
			responseCollector <- msg
			n++
		}()
	}

	// Remote workers
	peers := node.NodeTracker.GetAllNodeData()
	for _, p := range peers {
		for _, addr := range p.Multiaddrs {
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			if (p.PeerId.String() != node.Host.ID().String()) &&
				p.IsStaked &&
				node.NodeTracker.GetNodeData(p.PeerId.String()).CanDoWork(pubsub.WorkerCategory(message.Type)) {
				logrus.Infof("Worker Address: %s", ipAddr)
				wg.Add(1)
				go func(p pubsub.NodeData) {
					defer wg.Done()
					spawned, err := node.ActorRemote.SpawnNamed(fmt.Sprintf("%s:4001", ipAddr), "worker", "peer", -1)
					if err != nil {
						logrus.Debugf("Error spawning remote worker: %v", err)
						responseCollector <- &pubsub2.Message{
							ValidatorData: map[string]interface{}{"error": err.Error()},
						}
						return
					}
					spawnedPID := spawned.Pid
					logrus.Infof("Worker Address: %s", spawnedPID)
					if spawnedPID == nil {
						logrus.Errorf("Spawned PID is nil for IP: %s", ipAddr)
						responseCollector <- &pubsub2.Message{
							ValidatorData: map[string]interface{}{"error": "Spawned PID is nil"},
						}
						return
					}
					client := node.ActorEngine.Spawn(props)
					node.ActorEngine.Send(spawnedPID, &messages.Connect{Sender: client})
					future := node.ActorEngine.RequestFuture(spawnedPID, message, 30*time.Second)
					result, fErr := future.Result()
					if fErr != nil {
						logrus.Debugf("Error receiving response from remote worker: %v", fErr)
						responseCollector <- &pubsub2.Message{
							ValidatorData: map[string]interface{}{"error": fErr.Error()},
						}
						return
					}
					response := result.(*messages.Response)
					msg := &pubsub2.Message{}
					rErr := json.Unmarshal([]byte(response.Value), &msg)
					if rErr != nil {
						gMsg, gErr := getResponseMessage(response)
						if gErr != nil {
							logrus.Errorf("Error getting response message: %v", gErr)
							responseCollector <- &pubsub2.Message{
								ValidatorData: map[string]interface{}{"error": gErr.Error()},
							}
							return
						}
						if gMsg != nil {
							msg = gMsg
						}
					}
					responseCollector <- msg
					n++
					// cap at 3 for performance
					if n == len(peers) || n == 3 {
						logrus.Info("All workers have responded")
						responseCollector <- msg
					}
				}(p)
			}
		}
	}

	// Queue responses and send to workerDoneCh
	go func() {
		var responses []*pubsub2.Message
		for {
			select {
			case response := <-responseCollector:
				responses = append(responses, response)
			case <-timeout:
				for _, resp := range responses {
					workerDoneCh <- resp
				}
				return
			}
		}
	}()

	wg.Wait()
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
	node.WorkerTracker = &pubsub.WorkerEventTracker{WorkerStatusCh: workerStatusCh}
	err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.WorkerTopic), node.WorkerTracker, true)
	if err != nil {
		logrus.Errorf("Subscribe error %v", err)
	}

	// Register self as a remote node for the network
	node.ActorRemote.Register("peer", actor.PropsFromProducer(NewWorker(node)))

	if node.WorkerTracker == nil {
		logrus.Error("MonitorWorkers: WorkerTracker is nil")
		return
	}

	if node.WorkerTracker.WorkerStatusCh == nil {
		logrus.Error("MonitorWorkers: WorkerStatusCh is nil")
		return
	}

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()
	rcm := pubsub.GetResponseChannelMap()
	var startTime time.Time

	for {
		select {
		case work := <-node.WorkerTracker.WorkerStatusCh:
			logrus.Info("Sending work to network")
			var workData map[string]string

			err := json.Unmarshal(work.Data, &workData)
			if err != nil {
				logrus.Error("Error unmarshalling work: ", err)
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
				defer close(ch)
			} else {
				logrus.Debugf("Error processing data.ValidatorData: %v", data.ValidatorData)
				continue
			}

			processValidatorData(data, validatorDataMap, &startTime, node)

		case <-ticker.C:
			logrus.Info("worker tick")

		case <-ctx.Done():
			return
		}
	}
}

/**
 * Processes the validator data received from the network.
 *
 * @param {pubsub2.Message} data - The message data received from the network.
 * @param {map[string]interface{}} validatorDataMap - The map containing validator data.
 * @param {time.Time} startTime - The start time of the work.
 * @param {masa.OracleNode} node - The OracleNode instance.
 */
func processValidatorData(data *pubsub2.Message, validatorDataMap map[string]interface{}, startTime *time.Time, node *masa.OracleNode) {
	//logrus.Infof("Work validatorDataMap %s", validatorDataMap)
	if response, ok := validatorDataMap["Response"].(map[string]interface{}); ok {
		if _, ok := response["error"].(string); ok {
			logrus.Infof("Work failed %s", response["error"])

			// Set WorkerTimeout for the node
			nodeData := node.NodeTracker.GetNodeData(data.ReceivedFrom.String())
			if nodeData != nil {
				nodeData.WorkerTimeout = time.Now()
				node.NodeTracker.AddOrUpdateNodeData(nodeData, true)
			}

		} else if work, ok := response["data"].(string); ok {
			processWork(data, work, startTime, node)

		} else if w, ok := response["data"].(map[string]interface{}); ok {
			work, err := json.Marshal(w)

			if err != nil {
				logrus.Errorf("Error marshalling data.ValidatorData: %v", err)
				return
			}

			processWork(data, string(work), startTime, node)
		} else {
			work, err := json.Marshal(response["data"])
			if err != nil {
				logrus.Errorf("Error marshalling data.ValidatorData: %v", err)
				return
			}

			processWork(data, string(work), startTime, node)

		}
	}
}

/**
 * Processes the work received from the network.
 *
 * @param {pubsub2.Message} data - The message data received from the network.
 * @param {string} work - The work data as a string.
 * @param {time.Time} startTime - The start time of the work.
 * @param {masa.OracleNode} node - The OracleNode instance.
 */
func processWork(data *pubsub2.Message, work string, startTime *time.Time, node *masa.OracleNode) {
	key, _ := computeCid(work)
	logrus.Infof("Work done %s", key)

	endTime := time.Now()
	duration := endTime.Sub(*startTime)

	workEvent := db.WorkEvent{
		CID:       key,
		PeerId:    data.ID,
		Payload:   []byte(work),
		Duration:  duration.Seconds(),
		Timestamp: time.Now().Unix(),
	}
	logrus.Infof("Publishing work event : %v for Peer %s", workEvent.CID, workEvent.PeerId)
	logrus.Debugf("Publishing work event : %v", workEvent)

	_ = node.PubSubManager.Publish(config.TopicWithVersion(config.BlockTopic), workEvent.Payload)
}
