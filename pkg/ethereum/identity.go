package ethereum

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	libp2pCrypto "github.com/libp2p/go-libp2p/core/crypto"

	"github.com/masa-finance/masa-oracle/pkg/ethereum/contracts"
)

const (
	ChainId         = 5611
	RpcEndpoint     = "https://opbnb-testnet.nodereal.io/v1/99613329b67d43e3a52f5ebe7c666efc"
	ContractAddress = "0xb6f59e114f2bF57B1891f08fC23B3C696b7D3b16"
	PaymentMethod   = "0x0000000000000000000000000000000000000000"
)

func Mint(libP2pPrivateKey libp2pCrypto.PrivKey, toAddress string) error {
	// Connect to the Ethereum client

	rpcEndpoint := os.Getenv("rpc.endpoint")
	if rpcEndpoint == "" {
		return errors.New("rpc.endpoint is not set in the environment")
	}
	intVal, err := strconv.Atoi(strings.TrimSpace(os.Getenv("chain.id")))
	if err != nil {
		return err
	}
	chainId := big.NewInt(int64(intVal))
	if rpcEndpoint == "" {
		return errors.New("chain.id is not set in the environment")
	}

	client, err := ethclient.Dial(rpcEndpoint)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}
	paymentMeth := os.Getenv("payment.method")
	if paymentMeth == "" {
		return errors.New("payment.method is not set in the environment")
	}
	limit, err := strconv.ParseInt(os.Getenv("gas.limit.wei"), 10, 64)
	if err != nil {
		return err
	}

	// Create a new instance of the contract
	contractAddress := common.HexToAddress("0xb6f59e114f2bF57B1891f08fC23B3C696b7D3b16")

	instance, err := contracts.NewEthereum(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to create a new instance of the contract: %v", err)
	}
	ecdsaPrivKey, err := LibP2pToEcdsa(libP2pPrivateKey)
	if err != nil {
		return err
	}

	// Create a new transactor
	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivKey, chainId)
	auth.GasLimit = uint64(limit)

	// Call the mint function
	paymentMethod := common.HexToAddress(os.Getenv("payment.method"))
	toAddressHex := common.HexToAddress(toAddress)
	tx, err := instance.Mint0(auth, paymentMethod, toAddressHex)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to call the mint function: %v", err))
	}

	log.Printf("Transaction sent: %s", tx.Hash().Hex())
	return nil
}
