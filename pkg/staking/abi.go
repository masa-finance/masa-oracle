package staking

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/masa-finance/masa-oracle/contracts"
)

const (
	MasaTokenABIPath       = "node_modules/@masa-finance/masa-token/deployments/sepolia/MasaToken.json"
	MasaFaucetABIPath      = "node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/MasaFaucet.sol/MasaFaucet.json"
	ProtocolStakingABIPath = "node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/ProtocolStaking.sol/ProtocolStaking.json"
)

type ContractAddresses struct {
	Sepolia struct {
		MasaFaucet      string `json:"MasaFaucet"`
		MasaToken       string `json:"MasaToken"`
		ProtocolStaking string `json:"ProtocolStaking"`
	} `json:"sepolia"`
}

// GetABI parses the ABI from the given JSON file path.
// It returns the parsed ABI, or an error if reading or parsing fails.
func GetABI(jsonPath string) (abi.ABI, error) {
	jsonFile, err := contracts.EmbeddedContracts.ReadFile(jsonPath)
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
