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
	infuraURL       = "https://rpc.sepolia.org"                    // update to Sepolia - this should be added as an environment variable sometime
	contractAddress = "0xd925bc5d3eCd899a3F7B8D762397D2DC75E1187b" // this is the sepolia contract address
)

type Contract struct {
	ABI []interface{} `json:"abi"`
}

func getContractABI() []interface{} {
	jsonFile, err := ioutil.ReadFile("contracts/build/contracts/OracleNodeStakingContract.json")
	if err != nil {
		log.Fatalf("Failed to read contract JSON: %v", err)
	}

	var contract Contract
	err = json.Unmarshal(jsonFile, &contract)
	if err != nil {
		log.Fatalf("Failed to unmarshal contract JSON: %v", err)
	}

	return contract.ABI
}

func VerifyStakingEvent(userAddress string) bool {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractABI := getContractABI()
	abiJSON, err := json.Marshal(contractABI)
	if err != nil {
		log.Fatalf("Failed to marshal contract ABI: %v", err)
	}
	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
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
