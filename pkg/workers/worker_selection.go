package workers

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/types"
)

// GetEligibleWorkers Uses the new NodeTracker method to get the eligible workers for a given message type
// I'm leaving this returning an array so that we can easily increase the number of workers in the future
func GetEligibleWorkers(node *masa.OracleNode, category pubsub.WorkerCategory, config *WorkerConfig) ([]data_types.Worker, *data_types.Worker) {

	var workers []data_types.Worker
	nodes := node.NodeTracker.GetEligibleWorkerNodes(category)
	var localWorker *data_types.Worker

	// Shuffle the node list first to avoid always selecting the same node
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	logrus.Info("checking connections to eligible workers")
	start := time.Now()
	for _, eligible := range nodes {
		if eligible.PeerId.String() == node.Host.ID().String() {
			localWorker = &data_types.Worker{IsLocal: true, NodeData: eligible}
			continue
		}
		addr, err := multiaddr.NewMultiaddr(eligible.MultiaddrsString)
		if err != nil {
			logrus.Errorf("error creating multiaddress: %s", err.Error())
			continue
		}
		peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			logrus.Errorf("Failed to get peer info: %s", err)
			continue
		}
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
		defer cancel() // Cancel the context when done to release resources
		if err := node.Host.Connect(ctxWithTimeout, *peerInfo); err != nil {
			logrus.Debugf("Failed to connect to peer: %v", err)
			continue
		}
		workers = append(workers, data_types.Worker{IsLocal: false, NodeData: eligible, AddrInfo: peerInfo})
		// print duration of worker selection in seconds with floating point precision
		dur := time.Since(start).Milliseconds()
		logrus.Infof("Worker selection took %v milliseconds", dur)
		break
	}
	// make sure we get the local node in the list
	if localWorker == nil {
		nd := node.NodeTracker.GetNodeData(node.Host.ID().String())
		if nd.CanDoWork(category) {
			localWorker = &data_types.Worker{IsLocal: true, NodeData: *nd}
		}
	}
	return workers, localWorker
}
