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

// VerifyStakingEvent checks if the given user address has staked tokens by
// calling the stakes() view function on the ProtocolStaking contract.
// It connects to an Ethereum node, encodes the stakes call, calls the contract,
// unpacks the result, and returns true if the stakes amount is > 0.
func VerifyStakingEvent(rpcUrl string, userAddress string) (bool, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return false, fmt.Errorf("[-] Failed to connect to the Ethereum client: %v", err)
	}

	parsedABI, err := GetABI(ProtocolStakingABIPath) // Use the GetABI function from abi.go
	if err != nil {
		return false, err
	}

	address := common.HexToAddress(userAddress)
	stake, err := parsedABI.Pack("stakes", address)
	if err != nil {
		return false, fmt.Errorf("[-] Failed to pack data for stakes call: %v", err)
	}

	addresses, err := LoadContractAddresses()
	if err != nil {
		return false, fmt.Errorf("[-] Failed to load contract addresses: %v", err)
	}
	contractAddr := common.HexToAddress(addresses.Sepolia.ProtocolStaking)

	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: stake,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return false, fmt.Errorf("[-] Failed to call stakes function: %v", err)
	}

	stakesAmountInterfaces, err := parsedABI.Unpack("stakes", result)
	if err != nil {
		return false, fmt.Errorf("[-] Failed to unpack stakes: %v", err)
	}

	stakesAmount, ok := stakesAmountInterfaces[0].(*big.Int)
	if !ok {
		return false, errors.New("[-] Failed to assert type: stakesAmount is not *big.Int")
	}
	return stakesAmount.Cmp(big.NewInt(0)) > 0, nil
}
