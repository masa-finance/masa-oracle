package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type TwitterQueryHandler struct{}
type TwitterFollowersHandler struct{}
type TwitterProfileHandler struct{}
type TwitterSentimentHandler struct{}
type TwitterTrendsHandler struct{}

func (h *TwitterQueryHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterQueryHandler input: %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		logrus.Errorf("[+] TwitterQueryHandler error parsing data: %v", err)
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)

	logrus.Infof("[+] Scraping tweets for query: %s, count: %d", query, count)

	resp, err := twitter.ScrapeTweetsByQuery(query, count)
	if err != nil {
		logrus.Errorf("[+] TwitterQueryHandler error scraping tweets: %v", err)
		return data_types.WorkResponse{Error: err.Error()}
	}

	logrus.Infof("[+] TwitterQueryHandler response: %d tweets found", len(resp))

	return data_types.WorkResponse{Data: resp}
}

func (h *TwitterFollowersHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterFollowersHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter followers data: %v", err)}
	}
	username := dataMap["username"].(string)
	count := int(dataMap["count"].(float64))
	resp, err := twitter.ScrapeFollowersForProfile(username, count)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to get twitter followers: %v", err)}
	}
	return data_types.WorkResponse{Data: resp}
}

func (h *TwitterProfileHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterProfileHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter profile data: %v", err)}
	}
	username := dataMap["username"].(string)
	resp, err := twitter.ScrapeTweetsProfile(username)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to get twitter profile: %v", err)}
	}
	return data_types.WorkResponse{Data: resp}
}

func (h *TwitterSentimentHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterSentimentHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter sentiment data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)
	model := dataMap["model"].(string)
	_, resp, err := twitter.ScrapeTweetsForSentiment(query, count, model)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to get twitter sentiment: %v", err)}
	}
	return data_types.WorkResponse{Data: resp}
}

func (h *TwitterTrendsHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterTrendsHandler %s", data)
	resp, err := twitter.ScrapeTweetsByTrends()
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to get twitter trends: %v", err)}
	}
	return data_types.WorkResponse{Data: resp}
}
