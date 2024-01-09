// pkg/staking/config.go
package staking

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetRPCURL() string {
	return os.Getenv("RPC_URL")
}

func LoadContractAddresses() (*ContractAddresses, error) {
	path := filepath.Join("contracts", "node_modules", "@masa-finance", "masa-contracts-oracle", "addresses.json")
	data, err := ioutil.ReadFile(path)
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
