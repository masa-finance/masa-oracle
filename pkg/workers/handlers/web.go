package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
	types "github.com/masa-finance/tee-worker/api/types"
	worker "github.com/masa-finance/tee-worker/pkg/client"
)

// WebHandler - All the web handlers implement the WorkHandler interface.
type WebHandler struct{}

func (h *WebHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] WebHandler %s", data)
	client := worker.NewClient(teeWorkerURL)

	dataMap, err := JsonBytesToMap(data)
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
	resData := getSealedData(res)

	logrus.Infof("[+] WebHandler Work response for %s: %v returned", data_types.Web, string(resData))
	return data_types.WorkResponse{Data: resData}
}
