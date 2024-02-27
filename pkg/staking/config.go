package staking

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetRPCURL() (string, error) {
	url := viper.GetString("RPC_URL")
	if url == "" {
		return "", fmt.Errorf("%s environment variable is not set", "RPC_URL")
	}
	return url, nil
}

func LoadContractAddresses() (*ContractAddresses, error) {
	path := filepath.Join("contracts", "node_modules", "@masa-finance", "masa-contracts-oracle", "addresses.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var addresses ContractAddresses
	err = json.Unmarshal(data, &addresses)
	if err != nil {
		return nil, err
	}

	return &addresses, nil
}
