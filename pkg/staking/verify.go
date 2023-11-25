package staking

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL       = "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"
	contractAddress = "0xYourContractAddress"
	contractABI     = `[{"constant":true,"inputs":[],"name":"stakes","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"user","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"Staked","type":"event"}]`
)

func VerifyStakingEvent(userAddress string) bool {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	address := common.HexToAddress(userAddress)
	callOpts := &bind.CallOpts{}
	stake, err := contractAbi.Pack("stakes", address)
	if err != nil {
		log.Fatalf("Failed to pack data for stakes call: %v", err)
	}

	result, err := client.CallContract(context.Background(), callOpts, contractAddress, stake)
	if err != nil {
		log.Fatalf("Failed to call stakes function: %v", err)
	}

	stakes := new(big.Int)
	err = contractAbi.Unpack(stakes, "stakes", result)
	if err != nil {
		log.Fatalf("Failed to unpack stakes: %v", err)
	}

	return stakes.Cmp(big.NewInt(0)) > 0
}
