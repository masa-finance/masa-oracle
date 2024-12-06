package handlers

import (
	"encoding/json"
	"os"

	worker "github.com/masa-finance/tee-worker/pkg/client"
)

func JsonBytesToMap(jsonBytes []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func getSealedData(result *worker.JobResult) (resData []byte) {
	if os.Getenv("KEEP_SEALED_DATA") == "true" {
		res, err := result.Get()
		if err == nil {
			resData = []byte(res)
		}
	} else {
		res, err := result.GetDecrypted()
		if err == nil {
			resData = []byte(res)
		}
	}

	return
}
