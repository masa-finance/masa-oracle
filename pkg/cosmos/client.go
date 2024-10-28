package cosmos

import (
	"context"
	"log"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	myTypes "github.com/masa-finance/bobtestchain/x/bobtestchain/types"
)

type BlockchainClient struct {
	clientCtx client.Context
}

// NewBlockchainClient initializes a new BlockchainClient.
func NewBlockchainClient(nodeURI string, chainID string) (*BlockchainClient, error) {
	rpcClient, err := rpchttp.New(nodeURI, "/websocket")
	if err != nil {
		return nil, err
	}

	clientCtx := client.Context{}.
		WithChainID(chainID).
		WithClient(rpcClient).
		WithCodec(codec.NewProtoCodec(myTypes.ModuleCdc))

	return &BlockchainClient{clientCtx: clientCtx}, nil
}

// SendTransaction sends a transaction to the blockchain.
func (bc *BlockchainClient) SendTransaction(msg types.Msg) error {
	txFactory := tx.NewFactoryCLI(bc.clientCtx, nil)
	txBuilder := bc.clientCtx.TxConfig.NewTxBuilder()

	if err := txBuilder.SetMsgs(msg); err != nil {
		return err
	}

	txBytes, err := bc.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	res, err := bc.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	log.Printf("Transaction sent: %s", res.TxHash)
	return nil
}

// QueryNodeStatus queries the blockchain for a node's status.
func (bc *BlockchainClient) QueryNodeStatus(peerId string) (myTypes.NodeData, error) {
	queryClient := myTypes.NewQueryClient(bc.clientCtx)
	req := &myTypes.QueryNodeDataRequest{PeerId: peerId}

	res, err := queryClient.NodeStatus(context.Background(), req)
	if err != nil {
		return NodeData, err
	}

	return res.Status, nil
}
