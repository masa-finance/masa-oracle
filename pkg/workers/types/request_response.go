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
	WorkType  WorkerType `json:"workType,omitempty"`
	RequestId string     `json:"requestId,omitempty"`
	Data      []byte     `json:"data,omitempty"`
}

type WorkResponse struct {
	WorkRequest  *WorkRequest `json:"workRequest,omitempty"`
	Data         interface{}  `json:"data,omitempty"`
	Error        string       `json:"error,omitempty"`
	WorkerPeerId string       `json:"workerPeerId,omitempty"`
}
