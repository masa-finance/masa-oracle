package data_types

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
)

type Worker struct {
	IsLocal  bool
	IPAddr   string
	AddrInfo *peer.AddrInfo
	NodeData pubsub.NodeData
	Node     *node.OracleNode
}

func NewWorker(isLocal bool, nd *pubsub.NodeData) *Worker {
	var ma multiaddr.Multiaddr
	if len(nd.Multiaddrs) > 0 {
		ma = nd.Multiaddrs[0].Multiaddr
	} else {
		var err error
		ma, err = multiaddr.NewMultiaddr(nd.MultiaddrsString)
		if err != nil {
			logrus.Error(err)
			return nil
		}
	}
	ip, err := ma.ValueForProtocol(multiaddr.P_IP4)
	if err != nil {
		logrus.Error(err)
	}
	return &Worker{
		IsLocal:  isLocal,
		IPAddr:   ip,
		AddrInfo: nil,
		NodeData: pubsub.NodeData{},
		Node:     nil,
	}
}

type WorkRequest struct {
	WorkType     WorkerType `json:"workType,omitempty"`
	RequestId    string     `json:"requestId,omitempty"`
	Data         []byte     `json:"data,omitempty"`
	WorkerPeerId string     `json:"workerPeerId,omitempty"`
}

type WorkResponse struct {
	WorkRequest  *WorkRequest `json:"workRequest,omitempty"`
	Data         interface{}  `json:"data,omitempty"`
	Error        string       `json:"error,omitempty"`
	WorkerPeerId string       `json:"workerPeerId,omitempty"`
	RecordCount  int          `json:"recordCount,omitempty"`
	LoginEvent   *LoginEvent  `json:"loginEvent,omitempty"`
}
