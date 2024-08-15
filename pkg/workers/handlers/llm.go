package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

// TODO: LLMChatBody isn't used anywhere in the codebase. Remove after testing
type LLMChatBody struct {
	Model    string `json:"model,omitempty"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages,omitempty"`
	Stream bool `json:"stream"`
}

type LLMChatHandler struct{}

// HandleWork implements the WorkHandler interface for LLMChatHandler.
// It contains the logic for processing LLMChat work.
func (h *LLMChatHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] LLM Chat %s", data)
	uri := config.GetInstance().LLMChatUrl
	if uri == "" {
		return data_types.WorkResponse{Error: errors.New("missing env variable LLM_CHAT_URL")}
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return data_types.WorkResponse{Error: fmt.Errorf("unable to parse LLM chat data: %v", err)}
	}

	jsnBytes, err := json.Marshal(dataMap)
	if err != nil {
		return data_types.WorkResponse{Error: err}
	}
	resp, err := Post(uri, jsnBytes, nil)
	return data_types.WorkResponse{Data: resp, Error: err}
}
