package workers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/ipfs/go-cid"

	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

// WorkerTypeToCategory maps WorkerType to WorkerCategory
func WorkerTypeToCategory(wt WorkerType) pubsub.WorkerCategory {
	logrus.Infof("Mapping WorkerType %s to WorkerCategory", wt)
	switch wt {
	case Discord, DiscordProfile, DiscordChannelMessages, DiscordSentiment, DiscordGuildChannels, DiscordUserGuilds:
		logrus.Info("WorkerType is related to Discord")
		return pubsub.CategoryDiscord
	case TelegramSentiment, TelegramChannelMessages:
		logrus.Info("WorkerType is related to Telegram")
		return pubsub.CategoryTelegram
	case Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends:
		logrus.Info("WorkerType is related to Twitter")
		return pubsub.CategoryTwitter
	case Web, WebSentiment:
		logrus.Info("WorkerType is related to Web")
		return pubsub.CategoryWeb
	default:
		logrus.Warn("WorkerType is invalid or not recognized")
		return -1 // Invalid category
	}
}

// NewWorker creates a new instance of the Worker actor.
// It implements the actor.Receiver interface, allowing it to receive and handle messages.
//
// Returns:
//   - An instance of the Worker struct that implements the actor.Receiver interface.
func NewWorker(node *masa.OracleNode) actor.Producer {
	logrus.Info("Creating a new Worker actor")
	return func() actor.Actor {
		return &Worker{Node: node}
	}
}

// Receive is the message handling method for the Worker actor.
// It receives messages through the actor context and processes them based on their type.
func (a *Worker) Receive(ctx actor.Context) {
	logrus.Infof("Worker received a message of type %T from %s", ctx.Message(), ctx.Sender())
	switch m := ctx.Message().(type) {
	case *messages.Connect:
		logrus.Info("Handling Connect message")
		a.HandleConnect(ctx, m)
	case *actor.Started:
		logrus.Info("Actor has started")
		if a.Node.IsWorker() {
			a.HandleLog(ctx, "[+] Actor started")
		}
	case *actor.Stopping:
		logrus.Info("Actor is stopping")
		if a.Node.IsWorker() {
			a.HandleLog(ctx, "[+] Actor stopping")
		}
	case *actor.Stopped:
		logrus.Info("Actor has stopped")
		if a.Node.IsWorker() {
			a.HandleLog(ctx, "[+] Actor stopped")
		}
	case *messages.Work:
		logrus.Info("Handling Work message")
		if a.Node.IsWorker() {
			logrus.Infof("[+] Received Work")
			a.HandleWork(ctx, m, a.Node)
		}
	case *messages.Response:
		logrus.Info("Handling Response message")
		msg := &pubsub2.Message{}
		err := json.Unmarshal([]byte(m.Value), msg)
		if err != nil {
			logrus.Info("Unmarshalling Response message failed, attempting to get response message directly")
			msg, err = getResponseMessage(m)
			if err != nil {
				logrus.Errorf("[-] Error getting response message: %v", err)
				return
			}
		}
		logrus.Infof("Successfully handled Response message from %s, sending to workerDoneCh", msg.ReceivedFrom)
		workerDoneCh <- msg
		ctx.Poison(ctx.Self())
	default:
		logrus.Warningf("[+] Received unknown message: %T from %s", m, ctx.Sender())
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
	logrus.Infof("Computing CID for string: %s", str)
	// Create a multihash from the string
	mhHash, err := mh.Sum([]byte(str), mh.SHA2_256, -1)
	if err != nil {
		logrus.Errorf("Error computing multihash for string: %s, error: %v", str, err)
		return "", err
	}
	// Create a CID from the multihash
	cidKey := cid.NewCidV1(cid.Raw, mhHash).String()
	logrus.Infof("Computed CID: %s", cidKey)
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
	logrus.Info("Converting Response message to pubsub2.Message")
	responseData := map[string]interface{}{}

	err := json.Unmarshal([]byte(response.Value), &responseData)
	if err != nil {
		logrus.Errorf("Error unmarshalling Response message: %v", err)
		return nil, err
	}
	msg := &pubsub2.Message{
		ID:            responseData["ID"].(string),
		ReceivedFrom:  peer.ID(responseData["ReceivedFrom"].(string)),
		ValidatorData: responseData["ValidatorData"],
		Local:         responseData["Local"].(bool),
	}
	logrus.Infof("Successfully converted Response message to pubsub2.Message from %s", msg.ReceivedFrom)
	return msg, nil
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
	logrus.Info("Starting MonitorWorkers to monitor worker data")
	node.WorkerTracker = &pubsub.WorkerEventTracker{WorkerStatusCh: workerStatusCh}
	err := node.PubSubManager.AddSubscription(config.TopicWithVersion(config.WorkerTopic), node.WorkerTracker, true)
	if err != nil {
		logrus.Errorf("[-] Subscribe error %v", err)
	}

	// Register self as a remote node for the network
	logrus.Info("Registering self as a remote node for the network")
	node.ActorRemote.Register("peer", actor.PropsFromProducer(NewWorker(node)))

	if node.WorkerTracker == nil {
		logrus.Error("[-] MonitorWorkers: WorkerTracker is nil")
		return
	}

	if node.WorkerTracker.WorkerStatusCh == nil {
		logrus.Error("[-] MonitorWorkers: WorkerStatusCh is nil")
		return
	}

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()
	rcm := pubsub.GetResponseChannelMap()
	var startTime time.Time

	for {
		select {
		case work := <-node.WorkerTracker.WorkerStatusCh:
			logrus.Info("[+] Sending work to network")
			var workData map[string]string

			err := json.Unmarshal(work.Data, &workData)
			if err != nil {
				logrus.Error("[-] Error unmarshalling work: ", err)
				continue
			}
			startTime = time.Now()
			go SendWork(node, work)
		case data := <-workerDoneCh:
			logrus.Infof("Processing data from workerDoneCh, received from %s", data.ReceivedFrom)
			validatorDataMap, ok := data.ValidatorData.(map[string]interface{})
			if !ok {
				logrus.Errorf("[-] Error asserting type: %v", ok)
				continue
			}

			if validatorDataMap["ChannelId"] != nil {
				if ch, ok := rcm.Get(validatorDataMap["ChannelId"].(string)); ok {
					validatorData, err := json.Marshal(validatorDataMap["Response"])
					if err != nil {
						logrus.Errorf("[-] Error marshalling data.ValidatorData: %v", err)
						continue
					}
					ch <- validatorData
					defer close(ch)
				} else {
					logrus.Debugf("Channel not found for ChannelId: %v", validatorDataMap["ChannelId"])
					continue
				}
			} else {
				logrus.Debugf("ChannelId is nil in validatorDataMap: %v", validatorDataMap)
				continue
			}

			processValidatorData(data, validatorDataMap, &startTime, node)

		case <-ticker.C:
			logrus.Info("[+] Worker tick")

		case <-ctx.Done():
			logrus.Info("Context done, stopping MonitorWorkers")
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
	logrus.Infof("[+] Processing validator data from %s: %s", data.ReceivedFrom, validatorDataMap)
	logrus.Infof("[+] Processing validator data: %s", validatorDataMap)
	if response, ok := validatorDataMap["Response"].(map[string]interface{}); ok {
		if _, ok := response["error"].(string); ok {
			logrus.Infof("[+] Work failed %s", response["error"])

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
				logrus.Errorf("[-] Error marshalling data.ValidatorData: %v", err)
				return
			}

			processWork(data, string(work), startTime, node)
		} else {
			work, err := json.Marshal(response["data"])
			if err != nil {
				logrus.Errorf("[-] Error marshalling data.ValidatorData: %v", err)
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
	logrus.Infof("[+] Work done %s", key)

	endTime := time.Now()
	duration := endTime.Sub(*startTime)

	workEvent := db.WorkEvent{
		CID:       key,
		PeerId:    data.ID,
		Payload:   []byte(work),
		Duration:  duration.Seconds(),
		Timestamp: time.Now().Unix(),
	}
	logrus.Infof("[+] Publishing work event : %v for Peer %s", workEvent.CID, workEvent.PeerId)
	logrus.Debugf("[+] Publishing work event : %v", workEvent)

	_ = node.PubSubManager.Publish(config.TopicWithVersion(config.BlockTopic), workEvent.Payload)
}
