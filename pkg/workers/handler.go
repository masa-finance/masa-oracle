package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/telegram"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

type LLMChatBody struct {
	Model    string `json:"model,omitempty"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages,omitempty"`
	Stream bool `json:"stream"`
}

// getPeers is a function that takes an OracleNode as an argument and returns a slice of actor.PID pointers.
// These actor.PID pointers represent the peers of the given OracleNode in the network.
func getPeers(node *masa.OracleNode) []*actor.PID {
	var actors []*actor.PID
	peers := node.Host.Network().Peers()
	for _, p := range peers {
		conns := node.Host.Network().ConnsToPeer(p)
		for _, conn := range conns {
			addr := conn.RemoteMultiaddr()
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			if p.String() != node.Host.ID().String() {
				spawned, err := node.ActorRemote.SpawnNamed(fmt.Sprintf("%s:4001", ipAddr), "worker", "peer", -1)
				if err != nil {
					if strings.Contains(err.Error(), "future: dead letter") {
						logrus.Debugf("Ignoring dead letter error for peer %s: %v", p.String(), err)
						continue
					}
					logrus.Debugf("Spawned error %v", err)
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
func (a *Worker) HandleConnect(ctx actor.Context, m *messages.Connect) {
	logrus.Infof("[+] Worker %v connected", m.Sender)
	clients.Add(m.Sender)
}

// HandleLog is a method of the Worker struct that handles logging.
// It takes in an actor context and a string message as parameters.
func (a *Worker) HandleLog(ctx actor.Context, l string) {
	logrus.Info(l)
}

// HandleWork is a method of the Worker struct that handles the work assigned to a worker.
// It takes in an actor context and a Work message as parameters.
func (a *Worker) HandleWork(ctx actor.Context, m *messages.Work, node *masa.OracleNode) {
	var resp interface{}
	var err error

	var workData map[string]string
	err = json.Unmarshal([]byte(m.Data), &workData)
	if err != nil {
		logrus.Errorf("[-] Error parsing work data: %v", err)
		return
	}

	var bodyData map[string]interface{}
	if workData["body"] != "" {
		if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
			logrus.Errorf("[-] Error unmarshalling body: %v", err)
			return
		}
	}

	// Log the workData and bodyData before the switch
	logrus.Infof("[+] Message: %+v", m)
	logrus.Infof("[+] Worker: %+v", WORKER.DiscordChannelMessages)
	logrus.Infof("[+] WorkData: %+v", workData)
	logrus.Infof("[+] BodyData: %+v", bodyData)

	switch workData["request"] {
	case string(WORKER.DiscordProfile):
		userID := bodyData["userID"].(string)
		resp, err = discord.GetUserProfile(userID)
	case string(WORKER.DiscordChannelMessages):
		channelID := bodyData["channelID"].(string)
		resp, err = discord.GetChannelMessages(channelID, bodyData["limit"].(string), bodyData["before"].(string))
	case string(WORKER.DiscordSentiment):
		logrus.Infof("[+] Discord Channel Messages %s %s", m.Data, m.Sender)
		channelID := bodyData["channelID"].(string)
		_, resp, err = discord.ScrapeDiscordMessagesForSentiment(channelID, bodyData["model"].(string), bodyData["prompt"].(string))
	case string(WORKER.TelegramChannelMessages):
		logrus.Infof("[+] Telegram Channel Messages %s %s", m.Data, m.Sender)
		username := bodyData["username"].(string)
		resp, err = telegram.FetchChannelMessages(context.Background(), username) // Removed the underscore placeholder
	case string(WORKER.TelegramSentiment):
		logrus.Infof("[+] Telegram Channel Messages %s %s", m.Data, m.Sender)
		username := bodyData["username"].(string)
		_, resp, err = telegram.ScrapeTelegramMessagesForSentiment(context.Background(), username, bodyData["model"].(string), bodyData["prompt"].(string))
	case string(WORKER.DiscordGuildChannels):
		guildID := bodyData["guildID"].(string)
		resp, err = discord.GetGuildChannels(guildID)
	case string(WORKER.DiscordUserGuilds):
		resp, err = discord.GetUserGuilds()
	case string(WORKER.LLMChat):
		uri := config.GetInstance().LLMChatUrl
		if uri == "" {
			logrus.Error("[-] Missing env variable LLM_CHAT_URL")
			return
		}
		bodyBytes, _ := json.Marshal(bodyData)
		headers := map[string]string{
			"Content-Type": "application/json",
		}
		resp, _ = Post(uri, bodyBytes, headers)
	case string(WORKER.Twitter):
		query := bodyData["query"].(string)
		count := int(bodyData["count"].(float64))
		resp, err = twitter.ScrapeTweetsByQuery(query, count)
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
		host, _, err := net.SplitHostPort(m.Sender.Address)
		addrs := node.Host.Addrs()
		isLocalHost := false
		for _, addr := range addrs {
			addrStr := addr.String()
			if strings.HasPrefix(addrStr, "/ip4/") {
				ipStr := strings.Split(strings.Split(addrStr, "/")[2], "/")[0]
				if host == ipStr {
					isLocalHost = true
					break
				}
			}
		}

		if isLocalHost {
			logrus.Errorf("[-] Local node: Error processing request: %s", err.Error())
		} else {
			logrus.Errorf("[-] Remote node %s: Error processing request: %s", m.Sender, err.Error())
		}

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
			logrus.Errorf("[-] Error marshalling response: %v", err)
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
			logrus.Errorf("[-] Error marshalling response: %v", err)
			return
		}
		cfg := config.GetInstance()

		if cfg.TwitterScraper || cfg.DiscordScraper || cfg.TelegramScraper || cfg.WebScraper {
			ctx.Respond(&messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		}
		for _, pid := range getPeers(node) {
			ctx.Send(pid, &messages.Response{RequestId: workData["request_id"], Value: string(jsn)})
		}
	}
	ctx.Poison(ctx.Self())
}
