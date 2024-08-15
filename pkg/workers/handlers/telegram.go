package handlers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/telegram"
	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type TelegramSentimentHandler struct{}
type TelegramChannelHandler struct{}

// HandleWork implements the WorkHandler interface for TelegramSentimentHandler.
func (h *TelegramSentimentHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TelegramSentimentHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse telegram json data: %v", err)}
	}
	userName := dataMap["username"].(string)
	model := dataMap["model"].(string)
	prompt := dataMap["prompt"].(string)
	_, resp, err := telegram.ScrapeTelegramMessagesForSentiment(context.Background(), userName, model, prompt)
	return data_types.WorkResponse{Data: resp, Error: err}
}

// HandleWork implements the WorkHandler interface for TelegramChannelHandler.
func (h *TelegramChannelHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TelegramChannelHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse telegram json data: %v", err)}
	}
	userName := dataMap["username"].(string)
	resp, err := telegram.FetchChannelMessages(context.Background(), userName)
	return data_types.WorkResponse{Data: resp, Error: err}
}
