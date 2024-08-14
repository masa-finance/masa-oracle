package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
)

// GetEligibleWorkers Uses the new NodeTracker method to get the eligible workers for a given message type
// I'm leaving this returning an array so that we can easily increase the number of workers in the future
func GetEligibleWorkers(node *masa.OracleNode, category pubsub.WorkerCategory, config *WorkerConfig) ([]Worker, *Worker) {

	var workers []Worker
	nodes := node.NodeTracker.GetEligibleWorkerNodes(category)
	var localWorker *Worker

	// Shuffle the node list first to avoid always selecting the same node
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	logrus.Info("checking connections to eligible workers")
	start := time.Now()
	for _, eligible := range nodes {
		if eligible.PeerId.String() == node.Host.ID().String() {
			localWorker = &Worker{IsLocal: true, NodeData: eligible}
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
		workers = append(workers, Worker{IsLocal: false, NodeData: eligible, AddrInfo: peerInfo})
		// print duration of worker selection in seconds with floating point precision
		dur := time.Since(start).Milliseconds()
		logrus.Infof("Worker selection took %v milliseconds", dur)
		break
	}
	// make sure we get the local node in the list
	if localWorker == nil {
		nd := node.NodeTracker.GetNodeData(node.Host.ID().String())
		if nd.CanDoWork(category) {
			localWorker = &Worker{IsLocal: true, NodeData: *nd}
		}
	}
	return workers, localWorker
}

// GetEligibleWorkersOld Uses the new NodeTracker method to get the eligible workers for a given message type
// I'm leaving this returning an array so that we can easily increase the number of workers in the future
func GetEligibleWorkersOld(node *masa.OracleNode, message *messages.Work, config *WorkerConfig) []Worker {

	var workers []Worker
	category := getCategorytForMessage(message)
	// Get the eligible workers for the given message type. This will include the local node only if it is eligible
	// for this category of work.

	nodes := node.NodeTracker.GetEligibleWorkerNodes(category)
	var localWorker *Worker

	// Shuffle the node list first to avoid always selecting the same node
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	logrus.Info("checking connections to eligible workers")
	start := time.Now()
	workerFound := false
	for _, eligible := range nodes {
		if eligible.PeerId.String() == node.Host.ID().String() {
			localWorker = &Worker{IsLocal: true, NodeData: eligible}
			continue
		}
		for _, addr := range eligible.Multiaddrs {
			ipAddr, _ := addr.ValueForProtocol(multiaddr.P_IP4)
			realAddr := fmt.Sprintf("/ip4/%s/udp/4001/quic-v1/p2p/%s", ipAddr, eligible.PeerId.String())
			addr, err := multiaddr.NewMultiaddr(realAddr)
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
			workers = append(workers, Worker{IsLocal: false, NodeData: eligible, IPAddr: ipAddr})
			// print duration of worker selection in seconds with floating point precision
			dur := time.Since(start).Milliseconds()
			logrus.Infof("Worker selection took %v milliseconds", dur)
			workerFound = true
			break
		}
		if workerFound {
			break
		}
	}
	// make sure we get the local node in the list
	if localWorker == nil {
		nd := node.NodeTracker.GetNodeData(node.Host.ID().String())
		if nd.CanDoWork(category) {
			workers = append(workers, Worker{IsLocal: true, NodeData: *nd})
		}
	}
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
	logrus.Error("No workers available")
	return Worker{}
}
