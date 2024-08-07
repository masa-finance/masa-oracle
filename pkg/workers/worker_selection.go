package workers

import (
	"math/rand/v2"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
	"github.com/multiformats/go-multiaddr"
)

func GetEligibleWorkers(node *masa.OracleNode, message *messages.Work) []Worker {
	var workers []Worker

	// Always include a local worker
	workers = append(workers, Worker{IsLocal: true, NodeData: pubsub.NodeData{PeerId: node.Host.ID()}, Node: node})

	peers := node.NodeTracker.GetAllNodeData()
	for _, p := range peers {
		if isEligibleRemoteWorker(p, node, message) {
			for _, addr := range p.Multiaddrs {
				ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
				workers = append(workers, Worker{IsLocal: false, NodeData: p, IPAddr: ipAddr})
				break
			}
		}
	}

	// Shuffle the workers list
	rand.Shuffle(len(workers), func(i, j int) {
		workers[i], workers[j] = workers[j], workers[i]
	})

	return workers
}

func isEligibleRemoteWorker(p pubsub.NodeData, node *masa.OracleNode, message *messages.Work) bool {
	return (p.PeerId.String() != node.Host.ID().String()) &&
		p.IsStaked &&
		node.NodeTracker.GetNodeData(p.PeerId.String()).CanDoWork(pubsub.WorkerCategory(message.Type))
}

func NewRoundRobinIterator(workers []Worker) *roundRobinIterator {
	return &roundRobinIterator{
		workers: workers,
		index:   -1,
		tried:   make(map[int]bool),
	}
}

func (r *roundRobinIterator) HasNext() bool {
	return len(r.workers) > 0 && len(r.tried) < len(r.workers)
}

func (r *roundRobinIterator) Next() Worker {
	localWorker := Worker{}
	hasLocalWorker := false

	// First, try to find a remote worker
	for i := 0; i < len(r.workers); i++ {
		r.index = (r.index + 1) % len(r.workers)
		if !r.tried[r.index] {
			r.tried[r.index] = true
			if !r.workers[r.index].IsLocal {
				return r.workers[r.index]
			} else {
				localWorker = r.workers[r.index]
				hasLocalWorker = true
			}
		}
	}

	// If no untried remote worker is found, return the local worker if available
	if hasLocalWorker {
		return localWorker
	}

	// This should never happen if HasNext() is checked before calling Next()
	panic("No workers available")
}
