package staking

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var MasaTokenAddress common.Address
var ProtocolStakingContractAddress common.Address

type Client struct {
	EthClient  *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
}

// NewClient initializes a new ethClient.Client instance with the provided private key.
// It loads the contract addresses, initializes an Ethereum client, and returns
// a Client instance.
func NewClient(rpcUrl string, privateKey *ecdsa.PrivateKey) (*Client, error) {
	addresses, err := LoadContractAddresses()
	if err != nil {
		return nil, fmt.Errorf("[-] Failed to load contract addresses: %v", err)
	}

	MasaTokenAddress = common.HexToAddress(addresses.Sepolia.MasaToken)
	ProtocolStakingContractAddress = common.HexToAddress(addresses.Sepolia.ProtocolStaking)

	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	return &Client{
		EthClient:  client,
		PrivateKey: privateKey,
	}, nil
}
