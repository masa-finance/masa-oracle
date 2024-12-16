package utils

import "encoding/json"

// BytesToMap converts a byte slice to a map[string]interface{}
func BytesToMap(jsonBytes []byte) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}
