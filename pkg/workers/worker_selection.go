package workers

import (
	"encoding/json"
	"math/rand/v2"

	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
)

// GetEligibleWorkers Uses the new NodeTracker method to get the eligible workers for a given message type
func GetEligibleWorkers(node *masa.OracleNode, message *messages.Work) []Worker {
	var workers []Worker
	category := getCategorytForMessage(message)
	// Get the eligible workers for the given message type. This will include the local node only if it is eligible
	// for this category of work.
	for _, eligible := range node.NodeTracker.GetEligibleWorkerNodes(category) {
		if eligible.PeerId.String() == node.Host.ID().String() {
			workers = append(workers, Worker{IsLocal: true, NodeData: pubsub.NodeData{PeerId: node.Host.ID()}})
			continue
		}
		for _, addr := range eligible.Multiaddrs {
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			workers = append(workers, Worker{IsLocal: false, NodeData: eligible, IPAddr: ipAddr})
			break
		}
	}
	// Shuffle the workers list
	rand.Shuffle(len(workers), func(i, j int) {
		workers[i], workers[j] = workers[j], workers[i]
	})
	return workers
}

// right now the message has the Twitter work type hard coded so we have to get it from the message data
func getCategorytForMessage(message *messages.Work) pubsub.WorkerCategory {
	// TODO: can we get this fixed in the protobuf code?
	var workData map[string]string
	err := json.Unmarshal([]byte(message.Data), &workData)
	if err != nil {
		logrus.Errorf("[-] Error parsing work data: %v", err)
		return -1
	}

	workType, err := StringToWorkerType(workData["request"])
	if err != nil {
		logrus.Errorf("[-] Error parsing work type: %v", err)
		return -1
	}
	return WorkerTypeToCategory(workType)

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
