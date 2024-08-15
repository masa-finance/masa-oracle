package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/twitter"
	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type TwitterQueryHandler struct{}
type TwitterFollowersHandler struct{}
type TwitterProfileHandler struct{}
type TwitterSentimentHandler struct{}
type TwitterTrendsHandler struct{}

func (h *TwitterQueryHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterQueryHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse twitter query data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)
	resp, err := twitter.ScrapeTweetsByQuery(query, count)
	return data_types.WorkResponse{Data: resp, Error: err}
}

func (h *TwitterFollowersHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterFollowersHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse twitter followers data: %v", err)}
	}
	username := dataMap["username"].(string)
	count := int(dataMap["count"].(float64))
	resp, err := twitter.ScrapeFollowersForProfile(username, count)
	return data_types.WorkResponse{Data: resp, Error: err}
}

func (h *TwitterProfileHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterProfileHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse twitter profile data: %v", err)}
	}
	username := dataMap["username"].(string)
	resp, err := twitter.ScrapeTweetsProfile(username)
	return data_types.WorkResponse{Data: resp, Error: err}
}

func (h *TwitterSentimentHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterSentimentHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse twitter sentiment data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)
	model := dataMap["model"].(string)
	_, resp, err := twitter.ScrapeTweetsForSentiment(query, count, model)
	return data_types.WorkResponse{Data: resp, Error: err}
}

func (h *TwitterTrendsHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterTrendsHandler %s", data)
	resp, err := twitter.ScrapeTweetsByTrends()
	return data_types.WorkResponse{Data: resp, Error: err}
}
