package staking

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadContractAddresses loads the contract addresses from the addresses.json file.
// It returns a ContractAddresses struct containing the loaded addresses.
func LoadContractAddresses() (*ContractAddresses, error) {
	masaOracleTokensPath := filepath.Join("contracts", "node_modules", "@masa-finance", "masa-contracts-oracle", "addresses.json")
	masaTokenPath := filepath.Join("contracts", "node_modules", "@masa-finance", "masa-token", "addresses.json")
	masaTokenData, err := os.ReadFile(masaTokenPath)
	if err != nil {
		return nil, err
	}
	masaOracleTokensData, err := os.ReadFile(masaOracleTokensPath)
	if err != nil {
		return nil, err
	}
	var tokenAddresses map[string]map[string]string
	var addresses ContractAddresses
	err = json.Unmarshal(masaTokenData, &tokenAddresses)
	if err != nil {
		return nil, err
	}
	// used until we get someone to make a proper npm package
	addresses.Sepolia.MasaFaucet = "0x244813DaABFd59483fe34d1FdB3598AF06fE6c63"
	// used until we get someone to make a proper npm package
	addresses.Sepolia.MasaToken = tokenAddresses["sepolia"]["MasaToken"]
	err = json.Unmarshal(masaOracleTokensData, &tokenAddresses)
	if err != nil {
		return nil, err
	}
	addresses.Sepolia.NodeDataMetrics = tokenAddresses["sepolia"]["NodeDataMetrics"]
	addresses.Sepolia.NodeRewardPool = tokenAddresses["sepolia"]["NodeRewardPool"]
	addresses.Sepolia.OracleNodeStaking = tokenAddresses["sepolia"]["OracleNodeStaking"]
	addresses.Sepolia.StakedMasaToken = tokenAddresses["sepolia"]["StakedMasaToken"]

	return &addresses, nil
}
