package staking

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	masa "github.com/masa-finance/masa-oracle/pkg"
)

func GetRPCURL() (string, error) {
	url := viper.GetString(masa.RpcUrl)
	if url == "" {
		return "", errors.New(fmt.Sprintf("%s environment variable is not set", masa.RpcUrl))
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
