package staking

import (
	"encoding/json"
	"path/filepath"

	"github.com/masa-finance/masa-oracle/contracts"
)

// LoadContractAddresses loads the contract addresses from the addresses.json file.
// It returns a ContractAddresses struct containing the loaded addresses.
func LoadContractAddresses() (*ContractAddresses, error) {
	masaTokenPath := filepath.Join("node_modules", "@masa-finance", "masa-token", "addresses.json")

	masaTokenData, err := contracts.EmbeddedContracts.ReadFile(masaTokenPath)
	if err != nil {
		return nil, err
	}

	masaOracleTokensPath := filepath.Join("node_modules", "@masa-finance", "masa-contracts-oracle", "addresses.json")

	masaOracleTokensData, err := contracts.EmbeddedContracts.ReadFile(masaOracleTokensPath)
	if err != nil {
		return nil, err
	}

	var tokenAddresses map[string]map[string]string
	var addresses ContractAddresses
	err = json.Unmarshal(masaTokenData, &tokenAddresses)
	if err != nil {
		return nil, err
	}
	addresses.Sepolia.MasaToken = tokenAddresses["sepolia"]["MasaToken"]
	err = json.Unmarshal(masaOracleTokensData, &tokenAddresses)
	if err != nil {
		return nil, err
	}
	addresses.Sepolia.MasaFaucet = tokenAddresses["sepolia"]["MasaFaucet"]
	addresses.Sepolia.ProtocolStaking = tokenAddresses["sepolia"]["OracleNodeStaking"]

	return &addresses, nil
}
