// staking/contracts.go
package staking

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
func (sc *StakingClient) Approve(amount *big.Int) (*types.Receipt, error) {
	// Read the ABI from a JSON file
	abiJSON, err := ioutil.ReadFile("path/to/MasaTokenABI.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI: %v", err)
	}

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Retrieve the sender's address from the private key
	fromAddress := crypto.PubkeyToAddress(sc.PrivateKey.PublicKey)

	// Get the nonce for the sender's address
	nonce, err := sc.EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Define the value to send with the transaction, which is 0 for a token approve
	value := big.NewInt(0)

	// Pack the data to send with the transaction
	data, err := parsedABI.Pack("approve", OracleNodeStakingContractAddress, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data for approve: %v", err)
	}

	// Estimate gas limit and gas price dynamically based on the current network conditions
	gasPrice, err := sc.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// Estimate the gas limit for the approve function call
	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &MasaTokenAddress,
		Data: data,
	}
	gasLimit, err := sc.EthClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %v", err)
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, MasaTokenAddress, value, gasLimit, gasPrice, data)

	// Sign the transaction
	chainID, err := sc.EthClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), sc.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = sc.EthClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), sc.EthClient, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction to be mined: %v", err)
	}

	return receipt, nil
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
