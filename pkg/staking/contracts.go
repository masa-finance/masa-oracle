package staking

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getRPCURL() string {
	return os.Getenv("RPC_URL")
}

type ContractAddresses struct {
	Sepolia struct {
		MasaToken         string `json:"MasaToken"`
		OracleNodeStaking string `json:"OracleNodeStaking"`
		StakingMasaToken  string `json:"StakingMasaToken"`
	} `json:"sepolia"`
}

func LoadContractAddresses() (*ContractAddresses, error) {
	path := filepath.Join("contracts", "node_modules", "@masa-finance", "masa-contracts-oracle", "addresses.json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var addresses ContractAddresses
	err = json.Unmarshal(data, &addresses)
	if err != nil {
		return nil, err
	}

	return &addresses, nil
}

var MasaTokenAddress common.Address
var OracleNodeStakingContractAddress common.Address

type Client struct {
	EthClient  *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
}

func getABI(jsonPath string) (abi.ABI, error) {
	jsonFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI: %v", err)
	}

	var contract struct {
		ABI json.RawMessage `json:"abi"`
	}
	err = json.Unmarshal(jsonFile, &contract)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to unmarshal ABI JSON: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(contract.ABI)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %v", err)
	}

	return parsedABI, nil
}

func NewClient(privateKey *ecdsa.PrivateKey) (*Client, error) {
	addresses, err := LoadContractAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to load contract addresses: %v", err)
	}

	MasaTokenAddress = common.HexToAddress(addresses.Sepolia.MasaToken)
	OracleNodeStakingContractAddress = common.HexToAddress(addresses.Sepolia.OracleNodeStaking)

	rpcURL := getRPCURL() // Use the getRPCURL function to get the environment variable
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		EthClient:  client,
		PrivateKey: privateKey,
	}, nil
}

func (sc *Client) Approve(amount *big.Int) (string, error) {
	parsedABI, err := getABI("contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/MasaToken.sol/MasaToken.json")
	if err != nil {
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(sc.PrivateKey.PublicKey)
	nonce, err := sc.EthClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	value := big.NewInt(0)
	data, err := parsedABI.Pack("approve", OracleNodeStakingContractAddress, amount)
	if err != nil {
		return "", fmt.Errorf("failed to pack data for approve: %v", err)
	}

	gasPrice, err := sc.EthClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %v", err)
	}

	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &MasaTokenAddress,
		Data: data,
	}
	gasLimit, err := sc.EthClient.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %v", err)
	}

	tx := types.NewTransaction(nonce, MasaTokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := sc.EthClient.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get network ID: %v", err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), sc.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = sc.EthClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), sc.EthClient, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction receipt: %v", err)
	}
	if receipt.Status != 1 {
		return "", fmt.Errorf("transaction failed: %v", receipt)
	}

	return signedTx.Hash().Hex(), nil
}

func (sc *Client) Stake(amount *big.Int) (string, error) {
	chainID, err := sc.EthClient.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get network ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(sc.PrivateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create keyed transactor: %v", err)
	}

	parsedABI, err := getABI("contracts/node_modules/@masa-finance/masa-contracts-oracle/artifacts/contracts/OracleNodeStaking.sol/OracleNodeStaking.json")
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
