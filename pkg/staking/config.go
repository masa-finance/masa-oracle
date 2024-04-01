package staking

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadContractAddresses loads the contract addresses from the addresses.json file.
// It returns a ContractAddresses struct containing the loaded addresses.
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
