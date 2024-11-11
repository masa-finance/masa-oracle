package handlers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/scrapers/telegram"
	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

type TelegramChannelHandler struct{}

// HandleWork implements the WorkHandler interface for TelegramChannelHandler.
func (h *TelegramChannelHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TelegramChannelHandler %s", data)
	dataMap, err := JsonBytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse telegram json data: %v", err)}
	}
	userName := dataMap["username"].(string)
	resp, err := telegram.FetchChannelMessages(context.Background(), userName)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to get telegram channel messages: %v", err)}
	}
	logrus.Infof("[+] TelegramChannelHandler Work response for %s: %d records returned", data_types.TelegramChannelMessages, len(resp))
	return data_types.WorkResponse{Data: resp, RecordCount: len(resp)}
}
