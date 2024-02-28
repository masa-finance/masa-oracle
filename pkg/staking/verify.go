package staking

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func VerifyStakingEvent(userAddress string) (bool, error) {
	rpcURL, err := GetRPCURL()
	if err != nil {
		return false, err
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return false, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	parsedABI, err := GetABI(OracleNodeStakingABIPath) // Use the GetABI function from abi.go
	if err != nil {
		return false, err
	}

	address := common.HexToAddress(userAddress)
	stake, err := parsedABI.Pack("stakes", address)
	if err != nil {
		return false, fmt.Errorf("failed to pack data for stakes call: %v", err)
	}

	addresses, err := LoadContractAddresses()
	if err != nil {
		return false, fmt.Errorf("failed to load contract addresses: %v", err)
	}
	contractAddr := common.HexToAddress(addresses.Sepolia.OracleNodeStaking)

	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: stake,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return false, fmt.Errorf("failed to call stakes function: %v", err)
	}

	stakesAmountInterfaces, err := parsedABI.Unpack("stakes", result)
	if err != nil {
		return false, fmt.Errorf("failed to unpack stakes: %v", err)
	}

	stakesAmount, ok := stakesAmountInterfaces[0].(*big.Int)
	if !ok {
		return false, errors.New("failed to assert type: stakesAmount is not *big.Int")
	}
	return stakesAmount.Cmp(big.NewInt(0)) > 0, nil
}
