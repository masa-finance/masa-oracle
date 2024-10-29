package cosmos

import (
	"context"
	"log"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	protoCodec := codec.NewProtoCodec(interfaceRegistry)

	clientCtx := client.Context{}.
		WithChainID(chainID).
		WithClient(rpcClient).
		WithCodec(protoCodec)

	return &BlockchainClient{clientCtx: clientCtx}, nil
}

// SendTransaction sends a transaction to the blockchain.
func (bc *BlockchainClient) SendTransaction(msg sdk.Msg) error {
	//txFactory, err := tx.NewFactoryCLI(bc.clientCtx, nil)
	//if err != nil {
	//	return err
	//}
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
func (bc *BlockchainClient) QueryNodeStatus(peerId string) (*myTypes.NodeData, error) {
	queryClient := myTypes.NewNodeStatusQueryClient(bc.clientCtx)
	req := &myTypes.QueryNodeDataRequest{NodeId: peerId}

	res, err := queryClient.NodeData(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res.Status, nil
}
