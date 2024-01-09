package staking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getContractABI() (abi.ABI, error) {
	jsonFile, err := os.ReadFile("contracts/build/contracts/OracleNodeStakingContract.json")
	if err != nil {
		return abi.ABI{}, fmt.Errorf("Failed to read contract JSON: %v", err)
	}

	var contract struct {
		ABI json.RawMessage `json:"abi"`
	}
	err = json.Unmarshal(jsonFile, &contract)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("Failed to unmarshal contract JSON: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contract.ABI)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("Failed to parse contract ABI: %v", err)
	}

	return parsedABI, nil
}

func VerifyStakingEvent(userAddress string) (bool, error) {
	rpcURL := GetRPCURL()
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false, fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
	}

	parsedABI, err := getContractABI()
	if err != nil {
		return false, err
	}

	address := common.HexToAddress(userAddress)
	stake, err := parsedABI.Pack("stakes", address)
	if err != nil {
		return false, fmt.Errorf("Failed to pack data for stakes call: %v", err)
	}

	addresses, err := LoadContractAddresses()
	if err != nil {
		return false, fmt.Errorf("Failed to load contract addresses: %v", err)
	}
	contractAddr := common.HexToAddress(addresses.Sepolia.OracleNodeStaking)

	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: stake,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return false, fmt.Errorf("Failed to call stakes function: %v", err)
	}

	stakesAmountInterfaces, err := parsedABI.Unpack("stakes", result)
	if err != nil {
		return false, fmt.Errorf("Failed to unpack stakes: %v", err)
	}

	stakesAmount, ok := stakesAmountInterfaces[0].(*big.Int)
	if !ok {
		return false, errors.New("failed to assert type: stakesAmount is not *big.Int")
	}
	return stakesAmount.Cmp(big.NewInt(0)) > 0, nil
}
