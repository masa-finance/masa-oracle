package staking

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (sc *Client) Stake(amount *big.Int) (string, error) {
	chainID, err := sc.EthClient.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get network ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(sc.PrivateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create keyed transactor: %v", err)
	}

	parsedABI, err := GetABI(OracleNodeStakingABIPath)
	if err != nil {
		return "", err
	}

	stakingContract := bind.NewBoundContract(OracleNodeStakingContractAddress, parsedABI, sc.EthClient, sc.EthClient, sc.EthClient)
	if err != nil {
		return "", fmt.Errorf("failed to bind staking contract instance: %v", err)
	}

	tx, err := stakingContract.Transact(auth, "stake", amount)
	if err != nil {
		return "", fmt.Errorf("failed to send stake transaction: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), sc.EthClient, tx)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction receipt: %v", err)
	}
	if receipt.Status != 1 {
		return "", fmt.Errorf("transaction failed: %v", receipt)
	}

	return tx.Hash().Hex(), nil
}
