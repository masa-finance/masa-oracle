package ethereum

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func AddUser(privateKey *ecdsa.PrivateKey, chainId int64, userId, reputationScore string) (string, error) {
	// Connect to an ethereum node  running locally
	ethNodeUrl := os.Getenv("eth.node.url")
	if ethNodeUrl == "" {
		return "", errors.New("eth.node.url is not set in the environment")
	}
	client, err := ethclient.Dial(ethNodeUrl)
	if err != nil {
		return "", err
	}

	// Initialize transactor
	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainId))

	// Set up gas price and gas limit
	price, err := strconv.ParseInt(os.Getenv("gas.price.wei"), 10, 64)
	if err != nil {
		return "", err
	}
	limit, err := strconv.ParseInt(os.Getenv("gas.limit.wei"), 10, 64)
	if err != nil {
		return "", err
	}
	transactor.GasPrice = big.NewInt(price)
	transactor.GasLimit = uint64(limit)

	address := os.Getenv("reputation.voter.contract")
	if address == "" {
		return "", errors.New("reputation.voter.contract is not set in the environment")
	}
	// Address of the deployed contract
	contractAddress := common.HexToAddress(address)

	// Initialize a new instance of the contract bound to a specific deployed contract
	contract, err := NewPackageName(contractAddress, client)
	if err != nil {
		return "", err
	}

	// Call the contract's AddUser method
	tx, err := contract.AddUser(transactor, userId, reputationScore)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}
