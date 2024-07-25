package staking

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Approve approves the specified amount of MASA tokens for transfer from the
// caller's account to the ProtocolStakingContractAddress. It constructs an
// Ethereum transaction with the approve call, signs it, sends it to the network,
// waits for confirmation, and returns the transaction hash if successful.
func (sc *Client) Approve(amount *big.Int) (string, error) {
	parsedABI, err := GetABI(MasaTokenABIPath)
	if err != nil {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(sc.PrivateKey.PublicKey)
	nonce, err := sc.EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("[-] Failed to get nonce: %v", err)
	}

	value := big.NewInt(0)
	data, err := parsedABI.Pack("approve", ProtocolStakingContractAddress, amount)
	if err != nil {
		return "", fmt.Errorf("[-] Failed to pack data for approve: %v", err)
	}

	gasPrice, err := sc.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("[-] Failed to suggest gas price: %v", err)
	}

	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &MasaTokenAddress,
		Data: data,
	}
	gasLimit, err := sc.EthClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", fmt.Errorf("[-] Failed to estimate gas: %v", err)
	}

	tx := types.NewTransaction(nonce, MasaTokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := sc.EthClient.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("[-] Failed to get network ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), sc.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("[-] Failed to sign transaction: %v", err)
	}

	err = sc.EthClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("[-] Failed to send transaction: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), sc.EthClient, signedTx)
	if err != nil {
		return "", fmt.Errorf("[-] Failed to get transaction receipt: %v", err)
	}
	if receipt.Status != 1 {
		return "", fmt.Errorf("[-] Transaction failed: %v", receipt)
	}

	return signedTx.Hash().Hex(), nil
}
