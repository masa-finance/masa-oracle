package staking

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL       = "HTTP://127.0.0.1:7545"
	contractAddress = "0x767b636793c3399B7a517A6487974Bc474db1e7B"
)

type Contract struct {
	ABI string `json:"abi"`
}

func getContractABI() string {
	jsonFile, err := ioutil.ReadFile("contracts/build/contracts/OracleNodeStakingContract.json")
	if err != nil {
		log.Fatalf("Failed to read contract JSON: %v", err)
	}

	var contract Contract
	json.Unmarshal(jsonFile, &contract)

	return contract.ABI
}

func VerifyStakingEvent(userAddress string) bool {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractABI := getContractABI()
	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	address := common.HexToAddress(userAddress)
	stake, err := parsedABI.Pack("stakes", address)
	if err != nil {
		log.Fatalf("Failed to pack data for stakes call: %v", err)
	}

	// Address correction
	contractAddr := common.HexToAddress(contractAddress)

	callMsg := ethereum.CallMsg{
		To:   &contractAddr, // Adjusted to use a pointer
		Data: stake,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Failed to call stakes function: %v", err)
	}

	stakesAmountInterfaces, err := parsedABI.Unpack("stakes", result)
	if err != nil {
		log.Fatalf("Failed to unpack stakes: %v", err)
	}

	stakesAmount, ok := stakesAmountInterfaces[0].(*big.Int)
	if !ok {
		log.Fatalf("Failed to assert type: stakesAmount is not *big.Int")
	}

	return stakesAmount.Cmp(big.NewInt(0)) > 0

}
