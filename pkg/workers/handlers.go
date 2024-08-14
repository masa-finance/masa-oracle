package workers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/discord"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
)

//TODO: reduce all info level logs to debug level after testing

type LLMChatHandler struct{}

// HandleWork implements the WorkHandler interface for LLMChatHandler.
// It contains the logic for processing LLMChat work.
func (h *LLMChatHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] LLM Chat %s", data)
	uri := config.GetInstance().LLMChatUrl
	if uri == "" {
		return WorkResponse{Error: errors.New("missing env variable LLM_CHAT_URL")}
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse LLM chat data: %v", err)}
	}

	jsnBytes, err := json.Marshal(dataMap)
	if err != nil {
		return WorkResponse{Error: err}
	}
	resp, err := Post(uri, jsnBytes, nil)
	return WorkResponse{Data: resp, Error: err}
}

// DiscordHandler is a struct that implements the WorkHandler interface for Discord work.
type DiscordHandler struct{}

// HandleWork implements the WorkHandler interface for DiscordHandler.
func (h *DiscordHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] Discord %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse discord json data: %v", err)}
	}
	userID := dataMap["userID"].(string)
	resp, err := discord.GetUserProfile(userID)
	return WorkResponse{Data: resp, Error: err}
}

// TwitterQueryHandler All the twitter handlers implement the WorkHandler interface.
type TwitterQueryHandler struct{}
type TwitterFollowersHandler struct{}
type TwitterProfileHandler struct{}
type TwitterSentimentHandler struct{}
type TwitterTrendsHandler struct{}

func (h *TwitterQueryHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] TwitterQueryHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse twitter query data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)
	resp, err := twitter.ScrapeTweetsByQuery(query, count)
	return WorkResponse{Data: resp, Error: err}
}

func (h *TwitterFollowersHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] TwitterFollowersHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse twitter followers data: %v", err)}
	}
	username := dataMap["username"].(string)
	count := int(dataMap["count"].(float64))
	resp, err := twitter.ScrapeFollowersForProfile(username, count)
	return WorkResponse{Data: resp, Error: err}
}

func (h *TwitterProfileHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] TwitterProfileHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse twitter profile data: %v", err)}
	}
	username := dataMap["username"].(string)
	resp, err := twitter.ScrapeTweetsProfile(username)
	return WorkResponse{Data: resp, Error: err}
}

func (h *TwitterSentimentHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] TwitterSentimentHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse twitter sentiment data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)
	model := dataMap["model"].(string)
	_, resp, err := twitter.ScrapeTweetsForSentiment(query, count, model)
	return WorkResponse{Data: resp, Error: err}
}

func (h *TwitterTrendsHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] TwitterTrendsHandler %s", data)
	resp, err := twitter.ScrapeTweetsByTrends()
	return WorkResponse{Data: resp, Error: err}
}

// WebHandler - All the web handlers implement the WorkHandler interface.
type WebHandler struct{}
type WebSentimentHandler struct{}

func (h *WebHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] WebHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse web data: %v", err)}
	}
	depth := int(dataMap["depth"].(float64))
	urls := []string{dataMap["url"].(string)}
	resp, err := web.ScrapeWebData(urls, depth)
	return WorkResponse{Data: resp, Error: err}
}

func (h *WebSentimentHandler) HandleWork(data []byte) WorkResponse {
	logrus.Infof("[+] WebSentimentHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return WorkResponse{Error: fmt.Errorf("unable to parse web sentiment data: %v", err)}
	}
	depth := int(dataMap["depth"].(float64))
	urls := []string{dataMap["url"].(string)}
	model := dataMap["model"].(string)
	_, resp, err := web.ScrapeWebDataForSentiment(urls, depth, model)
	rslt := make(map[string]interface{})
	rslt["sentiment"] = resp
	return WorkResponse{Data: rslt, Error: err}
}

func JsonBytesToMap(jsonBytes []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}
