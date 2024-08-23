package workers

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

// GetEligibleWorkers Uses the new NodeTracker method to get the eligible workers for a given message type
// I'm leaving this returning an array so that we can easily increase the number of workers in the future
func GetEligibleWorkers(node *masa.OracleNode, category pubsub.WorkerCategory, config *WorkerConfig) ([]data_types.Worker, *data_types.Worker) {
	var workers []data_types.Worker
	nodes := node.NodeTracker.GetEligibleWorkerNodes(category)
	var localWorker *data_types.Worker

	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	logrus.Info("Checking connections to eligible workers")
	start := time.Now()
	for _, eligible := range nodes {
		if eligible.PeerId.String() == node.Host.ID().String() {
			localWorker = &data_types.Worker{IsLocal: true, NodeData: eligible}
			continue
		}

		// Use the DHT to find the peer's address information
		peerInfo, err := node.DHT.FindPeer(context.Background(), eligible.PeerId)
		if err != nil {
			logrus.Warnf("Failed to find peer %s in DHT: %v", eligible.PeerId.String(), err)
			continue
		}

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
		err = node.Host.Connect(ctxWithTimeout, peerInfo)
		cancel()
		if err != nil {
			logrus.Warnf("Failed to connect to peer %s: %v", eligible.PeerId.String(), err)
			continue
		}

		workers = append(workers, data_types.Worker{IsLocal: false, NodeData: eligible, AddrInfo: &peerInfo})
		dur := time.Since(start).Milliseconds()
		logrus.Infof("Worker selection took %v milliseconds", dur)
		break
	}

	if localWorker == nil {
		nd := node.NodeTracker.GetNodeData(node.Host.ID().String())
		if nd != nil && nd.CanDoWork(category) {
			localWorker = &data_types.Worker{IsLocal: true, NodeData: *nd}
		}
	}

	if len(workers) == 0 && localWorker == nil {
		logrus.Warn("No eligible workers found, including local worker")
	}

	return workers, localWorker
}
