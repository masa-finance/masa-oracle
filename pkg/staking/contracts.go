package staking

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var MasaTokenAddress common.Address
var OracleNodeStakingContractAddress common.Address

type Client struct {
	EthClient  *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
}

func NewClient(privateKey *ecdsa.PrivateKey) (*Client, error) {
	addresses, err := LoadContractAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to load contract addresses: %v", err)
	}

	MasaTokenAddress = common.HexToAddress(addresses.Sepolia.MasaToken)
	OracleNodeStakingContractAddress = common.HexToAddress(addresses.Sepolia.OracleNodeStaking)

	rpcURL, err := GetRPCURL()
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		EthClient:  client,
		PrivateKey: privateKey,
	}, nil
}
