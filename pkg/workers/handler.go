package workers

import (
	"encoding/json"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func getPeers(node *masa.OracleNode) []*actor.PID {
	var actors []*actor.PID
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
					return actors
				} else {
					actors = append(actors, spawned.Pid)
				}
			}
		}
	}
	return actors
}

// HandleConnect is a method of the Worker struct that handles the connection of a worker.
// It takes in an actor context and a Connect message as parameters.
// @todo fire event to masa sdk
func (a *Worker) HandleConnect(ctx actor.Context, m *messages.Connect) {
	logrus.Infof("[+] Worker %v connected", m.Sender)
	clients.Add(m.Sender)
}

// HandleLog is a method of the Worker struct that handles logging.
// It takes in an actor context and a string message as parameters.
// @todo fire event to masa sdk
func (a *Worker) HandleLog(ctx actor.Context, l string) {
	logrus.Info(l)
}

// HandleWork is a method of the Worker struct that handles the work assigned to a worker.
// It takes in an actor context and a Work message as parameters.
// @todo fire data to masa sdk
func (a *Worker) HandleWork(ctx actor.Context, m *messages.Work, node *masa.OracleNode) {
	var resp interface{}
	var err error

	var workData map[string]string
	err = json.Unmarshal([]byte(m.Data), &workData)
	if err != nil {
		logrus.Errorf("Error parsing work data: %v", err)
		return
	}

	var bodyData map[string]interface{}
	if workData["body"] != "" {
		if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
			logrus.Errorf("Error unmarshalling body: %v", err)
			return
		}
	}

	switch workData["request"] {
	case string(WORKER.Discord):
		logrus.Infof("[+] Discord %s %s", m.Data, m.Sender)
		userID := bodyData["userID"].(string)
		botToken := bodyData["botToken"].(string)
		resp, err = discord.GetUserProfile(userID, botToken)
	case string(WORKER.LLMChat):
		logrus.Infof("[+] LLM Chat %s %s", m.Data, m.Sender)
		uri := config.GetInstance().LLMChatUrl
		if uri == "" {
			logrus.Error("missing env variable LLM_CHAT_URL")
			return
		}
		resp, err = Post(uri, []byte(workData["body"]), nil)
	case string(WORKER.Twitter):
		count := int(bodyData["count"].(float64))
		resp, err = twitter.ScrapeTweetsByQuery(bodyData["query"].(string), count)
	case string(WORKER.TwitterFollowers):
		username := bodyData["username"].(string)
		count := int(bodyData["count"].(float64))
		resp, err = twitter.ScrapeFollowersForProfile(username, count)
	case string(WORKER.TwitterProfile):
		username := bodyData["username"].(string)
		resp, err = twitter.ScrapeTweetsProfile(username)
	case string(WORKER.TwitterSentiment):
		count := int(bodyData["count"].(float64))
		_, resp, err = twitter.ScrapeTweetsForSentiment(bodyData["query"].(string), count, bodyData["model"].(string))
	case string(WORKER.TwitterTrends):
		resp, err = twitter.ScrapeTweetsByTrends()
	case string(WORKER.Web):
		depth := int(bodyData["depth"].(float64))
		resp, err = web.ScrapeWebData([]string{bodyData["url"].(string)}, depth)
	case string(WORKER.WebSentiment):
		depth := int(bodyData["depth"].(float64))
		_, resp, err = web.ScrapeWebDataForSentiment([]string{bodyData["url"].(string)}, depth, bodyData["model"].(string))
	case string(WORKER.Test):
		count := int(bodyData["count"].(float64))
		resp, err = func(count int) (interface{}, error) {
			return count, err
		}(count)
	default:
		logrus.Warningf("[+] Received unknown message: %T", m)
		return
	}

	if err != nil {
		logrus.Errorf("Error processing request: %v", err)
		chanResponse := ChanResponse{
			Response:  map[string]interface{}{"error": err.Error()},
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
	} else {
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

		for _, pid := range getPeers(node) {
			ctx.Send(pid, &messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		}
	}
	ctx.Poison(ctx.Self())
}
