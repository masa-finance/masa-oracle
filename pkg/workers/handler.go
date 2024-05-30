package workers

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
	"github.com/sirupsen/logrus"
)

func (a *Worker) HandleConnect(ctx actor.Context, m *messages.Connect) {
	logrus.Infof("[+] Worker %v connected", m.Sender)
	clients.Add(m.Sender)
}

func (a *Worker) HandleLog(ctx actor.Context, l string) {
	logrus.Info(l)
}

func (a *Worker) HandleWork(ctx actor.Context, m *messages.Work) {
	// @todo fire data to masa sdk
	var workData map[string]string
	err := json.Unmarshal([]byte(m.Data), &workData)
	if err != nil {
		logrus.Errorf("Error parsing work data: %v", err)
		return
	}

	switch workData["request"] {
	case string(WORKER.Discord):
		logrus.Infof("[+] Discord %s %s", m.Data, m.Sender)
		var bodyData map[string]interface{}
		if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
			logrus.Errorf("Error unmarshalling body: %v", err)
			return
		}
		userID := bodyData["userID"].(string)
		botToken := bodyData["botToken"].(string)

		resp, err := discord.GetUserProfile(userID, botToken)
		if err != nil {
			logrus.Errorf("Error getting user profile: %v", err)
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
		return
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
		return
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
		return
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
		return
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
		return
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
		return
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
		return
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
		return
	case string(WORKER.WebSentiment):
		var bodyData map[string]interface{}
		if err := json.Unmarshal([]byte(workData["body"]), &bodyData); err != nil {
			logrus.Errorf("Error unmarshalling body: %v", err)
			return
		}
		depth := int(bodyData["depth"].(float64))
		_, resp, err := web.ScrapeWebDataForSentiment([]string{bodyData["url"].(string)}, depth, bodyData["model"].(string))
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
		return
	}
	ctx.Poison(ctx.Self())
}
