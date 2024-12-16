package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/tee"
	"github.com/masa-finance/masa-oracle/pkg/utils"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
	types "github.com/masa-finance/tee-worker/api/types"
)

// WebHandler - All the web handlers implement the WorkHandler interface.
type WebHandler struct{}

func (h *WebHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] WebHandler %s", data)
	client := tee.NewClient()

	dataMap, err := utils.BytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse web data: %v", err)}
	}
	depth := int(dataMap["depth"].(float64))

	res, err := client.SubmitJob(types.Job{
		Type: "web-scraper",
		Arguments: map[string]interface{}{
			"url":   dataMap["url"].(string),
			"depth": depth,
		},
	})
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse web query data: %v", err)}
	}

	result, err := res.Get()
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}

	logrus.Infof("[+] WebHandler Work response for %s: %v returned", data_types.Web, result)
	return data_types.WorkResponse{Data: result}
}
