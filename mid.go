package main

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// MasaIdentity represents a connection to the Masa Identity contract on the Ethereum blockchain.
type MasaIdentity struct {
	address  common.Address
	abi      abi.ABI
	contract *bind.BoundContract
}

// NewMasaIdentity initializes a new instance of MasaIdentity.
func NewMasaIdentity(client *ethclient.Client, contractAddress, abiJSON string) (*MasaIdentity, error) {
	address := common.HexToAddress(contractAddress)

	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	contract := bind.NewBoundContract(address, parsedABI, client, client, client)

	return &MasaIdentity{
		address:  address,
		abi:      parsedABI,
		contract: contract,
	}, nil
}

// Mint calls the "mint" function on the Masa Identity contract.
func (mi *MasaIdentity) Mint(opts *bind.TransactOpts, to common.Address, tokenId int) (*types.Transaction, error) {
	return mi.contract.Mint(opts, to, big.NewInt(int64(tokenId)))
}
