// staking/contracts.go
package staking

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Assuming you have the ABI of the MasaToken and OracleNodeStakingContract
const MasaTokenABI = `...`
const OracleNodeStakingContractABI = `...`

// Addresses of the deployed contracts (replace with actual addresses)
var MasaTokenAddress = common.HexToAddress("...")
var OracleNodeStakingContractAddress = common.HexToAddress("...")

// StakingClient holds the necessary details to interact with the Ethereum contracts
type StakingClient struct {
	EthClient  *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
}

// NewStakingClient creates a new StakingClient
func NewStakingClient(ethEndpoint string, privateKey *ecdsa.PrivateKey) (*StakingClient, error) {
	client, err := ethclient.Dial(ethEndpoint)
	if err != nil {
		return nil, err
	}
	return &StakingClient{
		EthClient:  client,
		PrivateKey: privateKey,
	}, nil
}

// Approve allows the staking contract to spend tokens on behalf of the user
func (sc *StakingClient) Approve(amount *big.Int) error {
	// Create an authenticated session
	auth, err := bind.NewKeyedTransactorWithChainID(sc.PrivateKey, big.NewInt(1)) // Replace with actual chain ID
	if err != nil {
		return err
	}

	// Create an instance of the MasaToken contract
	masaToken, err := NewMasaToken(MasaTokenAddress, sc.EthClient)
	if err != nil {
		return err
	}

	// Call the approve function of the MasaToken contract
	tx, err := masaToken.Approve(auth, OracleNodeStakingContractAddress, amount)
	if err != nil {
		return err
	}

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(context.Background(), sc.EthClient, tx)
	return err
}

// Stake allows the user to stake tokens
func (sc *StakingClient) Stake(amount *big.Int) error {
	// Create an authenticated session
	auth, err := bind.NewKeyedTransactorWithChainID(sc.PrivateKey, big.NewInt(1)) // Replace with actual chain ID
	if err != nil {
		return err
	}

	// Create an instance of the OracleNodeStakingContract
	stakingContract, err := NewOracleNodeStakingContract(OracleNodeStakingContractAddress, sc.EthClient)
	if err != nil {
		return err
	}

	// Call the stake function of the OracleNodeStakingContract
	tx, err := stakingContract.Stake(auth, amount)
	if err != nil {
		return err
	}

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(context.Background(), sc.EthClient, tx)
	return err
}
