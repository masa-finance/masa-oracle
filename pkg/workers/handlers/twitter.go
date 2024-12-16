package handlers

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/tee"
	"github.com/masa-finance/masa-oracle/pkg/utils"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
	types "github.com/masa-finance/tee-worker/api/types"
)

type TwitterQueryHandler struct{ MasaDir string }
type TwitterFollowersHandler struct{ MasaDir string }
type TwitterProfileHandler struct{ MasaDir string }

func (h *TwitterQueryHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterQueryHandler input: %s", data)
	dataMap, err := utils.BytesToMap(data)
	if err != nil {
		logrus.Errorf("[+] TwitterQueryHandler error parsing data: %v", err)
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}
	count := int(dataMap["count"].(float64))
	query := dataMap["query"].(string)

	logrus.Infof("[+] Scraping tweets for query: %s, count: %d", query, count)

	client := tee.NewClient()
	res, err := client.SubmitJob(types.Job{
		Type: "twitter-scraper",
		Arguments: map[string]interface{}{
			"type":  "searchbyquery",
			"query": query,
			"count": count,
		},
	})
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}

	result, err := res.Get()
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}

	logrus.Infof("[+] TwitterQueryHandler Work response for %s: %v", data_types.Twitter, result)
	return data_types.WorkResponse{Data: result}
}

func (h *TwitterFollowersHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterFollowersHandler %s", data)
	dataMap, err := utils.BytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter followers data: %v", err)}
	}
	username := dataMap["username"].(string)
	count := int(dataMap["count"].(float64))

	client := tee.NewClient()
	res, err := client.SubmitJob(types.Job{
		Type: "twitter-scraper",
		Arguments: map[string]interface{}{
			"type":  "searchfollowers",
			"query": username,
			"count": count,
		},
	})
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter followers data: %v", err)}
	}

	result, err := res.Get()
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}

	logrus.Infof("[+] TwitterQueryHandler Work response for %s: %v", data_types.Twitter, result)
	return data_types.WorkResponse{Data: result}
}

func (h *TwitterProfileHandler) HandleWork(data []byte) data_types.WorkResponse {
	logrus.Infof("[+] TwitterProfileHandler %s", data)
	dataMap, err := utils.BytesToMap(data)
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter profile data: %v", err)}
	}
	username := dataMap["username"].(string)

	client := tee.NewClient()
	res, err := client.SubmitJob(types.Job{
		Type: "twitter-scraper",
		Arguments: map[string]interface{}{
			"type":  "searchbyprofile",
			"query": username,
		},
	})
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}

	result, err := res.Get()
	if err != nil {
		return data_types.WorkResponse{Error: fmt.Sprintf("unable to parse twitter query data: %v", err)}
	}

	logrus.Infof("[+] TwitterQueryHandler Work response for %s: %v", data_types.Twitter, result)
	return data_types.WorkResponse{Data: result}
}
