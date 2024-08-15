package data_types

import (
	"github.com/libp2p/go-libp2p/core/peer"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type Worker struct {
	IsLocal  bool
	IPAddr   string
	AddrInfo *peer.AddrInfo
	NodeData pubsub.NodeData
	Node     *masa.OracleNode
}

type WorkRequest struct {
	WorkType  WorkerType
	RequestId string
	Data      []byte
}

type WorkResponse struct {
	WorkRequest  WorkRequest
	Data         interface{}
	Error        error
	WorkerPeerId string
}
