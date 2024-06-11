package staking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	MasaTokenABIPath         = "contracts/node_modules/@masa-finance/masa-token/deployments/sepolia/MasaToken.json"
	MasaFaucetABIPath        = "contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/MasaFaucet.sol/MasaFaucet.json"
	NodeDataMetricsABIPath   = "contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/NodeDataMetrics.sol/NodeDataMetrics.json"
	NodeRewardPoolABIPath    = "contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/NodeRewardPool.sol/NodeRewardPool.json"
	OracleNodeStakingABIPath = "contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/OracleNodeStaking.sol/OracleNodeStaking.json"
	StakedMasaTokenABIPath   = "contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/StakedMasaToken.sol/StakedMasaToken.json"
)

type ContractAddresses struct {
	Sepolia struct {
		MasaFaucet        string `json:"MasaFaucet"`
		MasaToken         string `json:"MasaToken"`
		NodeDataMetrics   string `json:"NodeDataMetrics"`
		NodeRewardPool    string `json:"NodeRewardPool"`
		OracleNodeStaking string `json:"OracleNodeStaking"`
		StakedMasaToken   string `json:"StakedMasaToken"`
	} `json:"sepolia"`
}

// GetABI parses the ABI from the given JSON file path.
// It returns the parsed ABI, or an error if reading or parsing fails.
func GetABI(jsonPath string) (abi.ABI, error) {
	jsonFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI: %v", err)
	}

	var contract struct {
		ABI json.RawMessage `json:"abi"`
	}
	err = json.Unmarshal(jsonFile, &contract)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to unmarshal ABI JSON: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contract.ABI)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return parsedABI, nil
}
