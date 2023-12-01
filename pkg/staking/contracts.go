// staking/contracts.go
package staking

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"strings"

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
func (sc *StakingClient) Approve(amount *big.Int) error {
	// Read the ABI from a JSON file
	abiJSON, err := ioutil.ReadFile("path/to/MasaTokenABI.json")
	if err != nil {
		return err
	}

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		return err
	}

	// Retrieve the sender's address from the private key
	fromAddress := crypto.PubkeyToAddress(sc.PrivateKey.PublicKey)

	// Get the nonce for the sender's address
	nonce, err := sc.EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	// Define the value to send with the transaction, which is 0 for a token approve
	value := big.NewInt(0)

	// Define gas limit and gas price; these should be determined based on the current network conditions
	// Placeholder values are used here and should be replaced with actual estimates
	gasLimit := uint64(21000) // This is just a placeholder value
	gasPrice, err := sc.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// Pack the data to send with the transaction
	data, err := parsedABI.Pack("approve", OracleNodeStakingContractAddress, amount)
	if err != nil {
		return err
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, MasaTokenAddress, value, gasLimit, gasPrice, data)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, sc.PrivateKey)
	if err != nil {
		return err
	}

	// Send the transaction
	err = sc.EthClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	// Transaction sent, you can now wait for it to be mined using the transaction hash if needed
	_, err = bind.WaitMined(context.Background(), sc.EthClient, signedTx)
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
