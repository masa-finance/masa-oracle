package staking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	MasaFaucetABIPath        = `{"abi":[{"inputs":[{"internalType":"address","name":"_tokenAddress","type":"address"},{"internalType":"uint256","name":"_tokenAmount","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"sender","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"EtherReceived","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"admin","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"EtherWithdrawn","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"admin","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"TokensWithdrawn","type":"event"},{"inputs":[],"name":"dailyLimit","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"dailyTokenReceived","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"faucet","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"nextRequestAt","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"_dailyLimit","type":"uint256"}],"name":"setDailyLimit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_tokenAmount","type":"uint256"}],"name":"setTokenAmount","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_waitTime","type":"uint256"}],"name":"setWaitTime","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"tokenAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"tokenAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"waitTime","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"withdrawEther","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"withdrawTokens","outputs":[],"stateMutability":"nonpayable","type":"function"},{"stateMutability":"payable","type":"receive"}]}`
	MasaTokenABIPath         = "contracts/node_modules/@masa-finance/masa-token/deployments/sepolia/MasaToken.json"
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
