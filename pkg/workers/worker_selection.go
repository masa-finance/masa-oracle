package workers

import (
	"github.com/multiformats/go-multiaddr"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
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
