package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/db"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
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
	LLMChat, Twitter, TwitterFollowers, TwitterProfile, TwitterSentiment, TwitterTrends, Web, WebSentiment WorkerType
}{
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

type OracleData struct {
	Id        string `json:"id"`
	PeerId    string `json:"peer_id"`
	Request   string `json:"request"`
	ModelName string `json:"model_name,omitempty"`
	Steps     []struct {
		Idx               int    `json:"idx"`
		RawContent        string `json:"raw_content,omitempty"`
		StructuredContent string `json:"structured_content,omitempty"`
		SystemPrompt      string `json:"system_prompt,omitempty"`
		Timestamp         string `json:"timestamp"`
		UserPrompt        string `json:"user_prompt,omitempty"`
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
		logrus.Infof("[+] Worker %v connected", m.Sender)
		clients.Add(m.Sender)
	case *actor.Started:
		logrus.Info("[+] Actor started")
	case *actor.Stopping:
		logrus.Info("[+] Actor stopping")
	case *actor.Stopped:
		logrus.Info("[+] Actor stopped")
	case *messages.Work:
		var workData map[string]string
		err := json.Unmarshal([]byte(m.Data), &workData)
		if err != nil {
			logrus.Errorf("Error parsing work data: %v", err)
			return
		}

		// WIP oracle work data
		//id := uuid.New().String()
		//oracleData := OracleData{
		//	Id:        id,
		//	PeerId:    m.Id,
		//	Request:   workData["request"],
		//	ModelName: workData["model"],
		//	Steps: []struct {
		//		Idx               int    `json:"idx"`
		//		RawContent        string `json:"raw_content,omitempty"`
		//		StructuredContent string `json:"structured_content,omitempty"`
		//		SystemPrompt      string `json:"system_prompt,omitempty"`
		//		Timestamp         string `json:"timestamp"`
		//		UserPrompt        string `json:"user_prompt,omitempty"`
		//	}{
		//		{
		//			Idx:        0,
		//			Timestamp:  time.Now().String(),
		//			RawContent: `Actor Started`,
		//		},
		//		{
		//			Idx:        1,
		//			RawContent: workData["request"],
		//			Timestamp:  time.Now().String(),
		//		},
		//		{
		//			Idx:        2,
		//			Timestamp:  time.Now().String(),
		//			RawContent: `Actor Stopped`,
		//		},
		//	},
		//}
		//jsonOD, _ := json.Marshal(oracleData)
		//jsonPayload := map[string]string{
		//	"raw": string(jsonOD),
		//}
		//err = db.SendToS3(id, jsonPayload)
		//if err != nil {
		//	logrus.Errorf("[-] Failed to send oracle data: %v", err)
		//}

		switch workData["request"] {
		case string(WORKER.LLMChat):
			logrus.Infof("[+] LLM Chat %s %s", m.Data, m.Sender)
			uri := config.GetInstance().LLMChatUrl
			if uri == "" {
				logrus.Error("missing env variable LLM_CHAT_URL")
				return
			}
			resp, err := Post(uri, []byte(workData["body"]), nil)
			if err != nil {
				return
			}
			var temp map[string]interface{}
			err = json.Unmarshal(resp, &temp)
			if err != nil {
				logrus.Errorf("Error unmarshalling response: %v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  temp,
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		case string(WORKER.Twitter):
			var bodyData map[string]interface{}
			if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
				logrus.Errorf("Error unmarshalling body: %v", err)
				return
			}
			count := int(bodyData["count"].(float64))
			resp, err := twitter.ScrapeTweetsByQuery(bodyData["query"].(string), count)
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			var temp []map[string]interface{}
			respBytes, err := json.Marshal(resp)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			err = json.Unmarshal(respBytes, &temp)
			if err != nil {
				logrus.Errorf("Error unmarshalling response: %v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  map[string]interface{}{"data": temp},
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		case string(WORKER.TwitterFollowers):
			var bodyData map[string]interface{}
			if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
				logrus.Errorf("Error unmarshalling body: %v", err)
				return
			}
			username := bodyData["username"].(string)
			count := int(bodyData["count"].(float64))

			resp, err := twitter.ScrapeFollowersForProfile(username, count)
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  map[string]interface{}{"data": resp},
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		case string(WORKER.TwitterProfile):
			var bodyData map[string]interface{}
			if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
				logrus.Errorf("Error unmarshalling body: %v", err)
				return
			}
			username := bodyData["username"].(string)

			resp, err := twitter.ScrapeTweetsProfile(username)
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  map[string]interface{}{"data": resp},
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		case string(WORKER.TwitterSentiment):
			var bodyData map[string]interface{}
			if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
				logrus.Errorf("Error unmarshalling body: %v", err)
				return
			}
			count := int(bodyData["count"].(float64))
			_, resp, err := twitter.ScrapeTweetsForSentiment(bodyData["query"].(string), count, bodyData["model"].(string))
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  map[string]interface{}{"data": resp},
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		case string(WORKER.TwitterTrends):
			resp, err := twitter.ScrapeTweetsByTrends()
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  map[string]interface{}{"data": resp},
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		case string(WORKER.Web):
			var bodyData map[string]interface{}
			if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
				logrus.Errorf("Error unmarshalling body: %v", err)
				return
			}
			depth := int(bodyData["depth"].(float64))
			resp, err := web.ScrapeWebData([]string{bodyData["url"].(string)}, depth)
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}

			var temp map[string]interface{}
			err = json.Unmarshal([]byte(resp), &temp)
			if err != nil {
				logrus.Errorf("Error unmarshalling response: %v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  temp,
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		case string(WORKER.WebSentiment):
			var bodyData map[string]interface{}
			if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
				logrus.Errorf("Error unmarshalling body: %v", err)
				return
			}
			depth := int(bodyData["depth"].(float64))
			_, resp, err := web.ScrapeWebDataForSentiment([]string{workData["url"]}, depth, workData["model"])
			if err != nil {
				logrus.Errorf("%v", err)
				return
			}
			var temp map[string]interface{}
			err = json.Unmarshal([]byte(resp), &temp)
			if err != nil {
				logrus.Errorf("Error unmarshalling response: %v", err)
				return
			}
			chanResponse := ChanResponse{
				Response:  temp,
				ChannelId: workData["request_id"],
			}
			val := &pubsub2.Message{
				ValidatorData: chanResponse,
				ID:            m.Id,
			}
			jsn, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Error marshalling response: %v", err)
				return
			}
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		}
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
func SendWork(node *masa.OracleNode, m *pubsub2.Message) {
	props := actor.PropsFromProducer(NewWorker())
	pid := node.ActorEngine.Spawn(props)
	message := &messages.Work{Data: string(m.Data), Sender: pid, Id: m.ReceivedFrom.String()}
	if node.IsActor() {
		// node.ActorEngine.Send(pid, message)
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
	}
	peers := node.Host.Network().Peers()
	for _, p := range peers {
		conns := node.Host.Network().ConnsToPeer(p)
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
			}
		}
	}
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
			key, _ := computeCid(string(validatorData))
			logrus.Infof("[+] Work done %s", key)
			// @todo to use FireData ...
			updateRecords(node, validatorData, key, data.ID)
		case <-ctx.Done():
			return
		}
	}
}
