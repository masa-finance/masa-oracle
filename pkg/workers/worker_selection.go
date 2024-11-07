package workers

import (
	"math/rand"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

// GetEligibleWorkers returns eligible workers for a given message type.
// For Twitter workers, it uses a balanced approach between high-performing workers and fair distribution.
// For other worker types, it returns all eligible workers without modification.
func GetEligibleWorkers(node *node.OracleNode, category pubsub.WorkerCategory, limit int) ([]data_types.Worker, *data_types.Worker) {
	nodes := node.NodeTracker.GetEligibleWorkerNodes(category)

	logrus.Infof("Getting eligible workers for category: %s", category)

	if category == pubsub.CategoryTwitter {
		return getTwitterWorkers(node, nodes, limit)
	}

	// For non-Twitter categories, return all eligible workers without modification
	return getAllWorkers(node, nodes, limit)
}

// getTwitterWorkers selects and shuffles a pool of top-performing Twitter workers
func getTwitterWorkers(node *node.OracleNode, nodes []pubsub.NodeData, limit int) ([]data_types.Worker, *data_types.Worker) {
	poolSize := calculatePoolSize(len(nodes), limit)
	topPerformers := nodes[:poolSize]

	// Shuffle the top performers
	rand.Shuffle(len(topPerformers), func(i, j int) {
		topPerformers[i], topPerformers[j] = topPerformers[j], topPerformers[i]
	})

	return createWorkerList(node, topPerformers, limit)
}

// getAllWorkers returns all eligible workers for non-Twitter categories
func getAllWorkers(node *node.OracleNode, nodes []pubsub.NodeData, limit int) ([]data_types.Worker, *data_types.Worker) {
	return createWorkerList(node, nodes, limit)
}

// createWorkerList creates a list of workers from the given nodes, respecting the limit
func createWorkerList(node *node.OracleNode, nodes []pubsub.NodeData, limit int) ([]data_types.Worker, *data_types.Worker) {
	workers := make([]data_types.Worker, 0, limit)
	var localWorker *data_types.Worker

	for _, eligible := range nodes {
		if eligible.PeerId.String() == node.Host.ID().String() {
			localAddrInfo := peer.AddrInfo{
				ID:    node.Host.ID(),
				Addrs: node.Host.Addrs(),
			}
			localWorker = &data_types.Worker{IsLocal: true, NodeData: eligible, AddrInfo: &localAddrInfo}
			continue
		}
		workers = append(workers, data_types.Worker{IsLocal: false, NodeData: eligible})

		// Apply limit if specified
		if limit > 0 && len(workers) >= limit {
			break
		}
	}

	logrus.Infof("Found %d eligible remote workers", len(workers))
	return workers, localWorker
}

// calculatePoolSize determines the size of the top performers pool for Twitter workers
func calculatePoolSize(totalNodes, limit int) int {
	if limit <= 0 {
		return totalNodes // If no limit, consider all nodes
	}
	// Use the larger of 5, double the limit, or 20% of total nodes
	poolSize := max(5, limit*2)
	poolSize = max(poolSize, totalNodes/5)
	return min(poolSize, totalNodes)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
