package staking

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (sc *Client) RunFaucet() (string, error) {
	chainID, err := sc.EthClient.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get network ID: %v", err)
	}

	// Load the contract addresses
	addresses, err := LoadContractAddresses()
	if err != nil {
		return "", fmt.Errorf("failed to load contract addresses: %v", err)
	}

	parsedABI, err := GetABI(MasaFaucetABIPath)
	if err != nil {
		return "", fmt.Errorf("failed to get ABI: %v", err)
	}

	// Create a new transactor
	auth, err := bind.NewKeyedTransactorWithChainID(sc.PrivateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create keyed transactor: %v", err)
	}

	// Bind the contract
	faucetContract := bind.NewBoundContract(common.HexToAddress(addresses.Sepolia.MasaFaucet), parsedABI, sc.EthClient, sc.EthClient, sc.EthClient)

	// Call the faucet function
	tx, err := faucetContract.Transact(auth, "faucet")
	if err != nil {
		return "", fmt.Errorf("failed to call faucet function: %v", err)
	}

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), sc.EthClient, tx)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// Check the status of the transaction
	if receipt.Status != 1 {
		return "", fmt.Errorf("transaction failed: %v", receipt)
	}

	return tx.Hash().Hex(), nil
}
