package workers

import (
	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
	"github.com/multiformats/go-multiaddr"
)

func GetEligibleWorkers(node *masa.OracleNode, message *messages.Work) []Worker {
	var workers []Worker

	if node.IsStaked && node.IsWorker() {
		workers = append(workers, Worker{IsLocal: true, NodeData: pubsub.NodeData{PeerId: node.Host.ID()}})
	}

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
	}
}

func (r *roundRobinIterator) HasNext() bool {
	return len(r.workers) > 0
}

func (r *roundRobinIterator) Next() Worker {
	r.index = (r.index + 1) % len(r.workers)
	return r.workers[r.index]
}
