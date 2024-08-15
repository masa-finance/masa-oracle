package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/web"
	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

// WebHandler - All the web handlers implement the WorkHandler interface.
type WebHandler struct{}
type WebSentimentHandler struct{}

func (h *WebHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] WebHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse web data: %v", err)}
	}
	depth := int(dataMap["depth"].(float64))
	urls := []string{dataMap["url"].(string)}
	resp, err := web.ScrapeWebData(urls, depth)
	return data_types.WorkResponse{Data: resp, Error: err}
}

func (h *WebSentimentHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] WebSentimentHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse web sentiment data: %v", err)}
	}
	depth := int(dataMap["depth"].(float64))
	urls := []string{dataMap["url"].(string)}
	model := dataMap["model"].(string)
	_, resp, err := web.ScrapeWebDataForSentiment(urls, depth, model)
	rslt := make(map[string]interface{})
	rslt["sentiment"] = resp
	return data_types.WorkResponse{Data: rslt, Error: err}
}
