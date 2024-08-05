package workers

import (
	"fmt"
	"sync"
	"time"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"
	"github.com/multiformats/go-multiaddr"

	"github.com/asynkron/protoactor-go/actor"

	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

const (
	workerTimeout = 30 * time.Second
	queueTimeout  = 8 * time.Second
	maxWorkers    = 5
)

type Worker struct {
	IsLocal  bool
	NodeData pubsub.NodeData
	IPAddr   string
	Node     *masa.OracleNode
}

func SendWork(node *masa.OracleNode, m *pubsub2.Message) {
	logrus.Infof("Sending work to node %s", node.Host.ID())
	props := actor.PropsFromProducer(NewWorker(node))
	pid := node.ActorEngine.Spawn(props)
	message := createWorkMessage(m, pid)

	responseCollector := make(chan *pubsub2.Message, 100)
	timeout := time.After(queueTimeout)

	eligibleWorkers := getEligibleWorkers(node, message)
	workerIterator := newRoundRobinIterator(eligibleWorkers)

	var wg sync.WaitGroup
	for i := 0; i < maxWorkers && workerIterator.HasNext(); i++ {
		worker := workerIterator.Next()
		wg.Add(1)
		go func(w Worker) {
			defer wg.Done()
			if w.IsLocal {
				handleLocalWorker(node, pid, message, responseCollector)
			} else {
				handleRemoteWorker(node, w.NodeData, w.IPAddr, props, message, responseCollector)
			}
		}(worker)
	}

	go queueResponses(responseCollector, timeout)

	wg.Wait()
}

func createWorkMessage(m *pubsub2.Message, pid *actor.PID) *messages.Work {
	return &messages.Work{
		Data:   string(m.Data),
		Sender: pid,
		Id:     m.ReceivedFrom.String(),
		Type:   int64(pubsub.CategoryTwitter),
	}
}

func getEligibleWorkers(node *masa.OracleNode, message *messages.Work) []Worker {
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

type roundRobinIterator struct {
	workers []Worker
	index   int
}

func newRoundRobinIterator(workers []Worker) *roundRobinIterator {
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

func handleLocalWorker(node *masa.OracleNode, pid *actor.PID, message *messages.Work, responseCollector chan<- *pubsub2.Message) {
	logrus.Info("Sending work to local worker")
	future := node.ActorEngine.RequestFuture(pid, message, workerTimeout)
	result, err := future.Result()
	if err != nil {
		handleWorkerError(err, responseCollector)
		return
	}
	processWorkerResponse(result, node.Host.ID(), responseCollector)
}

func isEligibleRemoteWorker(p pubsub.NodeData, node *masa.OracleNode, message *messages.Work) bool {
	return (p.PeerId.String() != node.Host.ID().String()) &&
		p.IsStaked &&
		node.NodeTracker.GetNodeData(p.PeerId.String()).CanDoWork(pubsub.WorkerCategory(message.Type))
}

func handleRemoteWorker(node *masa.OracleNode, p pubsub.NodeData, ipAddr string, props *actor.Props, message *messages.Work, responseCollector chan<- *pubsub2.Message) {
	logrus.WithFields(logrus.Fields{
		"ip":   ipAddr,
		"peer": p.PeerId,
	}).Info("Handling remote worker")

	spawned, err := node.ActorRemote.SpawnNamed(fmt.Sprintf("%s:4001", ipAddr), "worker", "peer", -1)
	if err != nil {
		logrus.WithError(err).WithField("ip", ipAddr).Error("Failed to spawn remote worker")
		handleWorkerError(err, responseCollector)
		return
	}

	spawnedPID := spawned.Pid
	if spawnedPID == nil {
		err := fmt.Errorf("failed to spawn remote worker: PID is nil for IP %s", ipAddr)
		logrus.WithFields(logrus.Fields{
			"ip":   ipAddr,
			"peer": p.PeerId,
		}).Error(err)
		handleWorkerError(err, responseCollector)
		return
	}

	client := node.ActorEngine.Spawn(props)
	node.ActorEngine.Send(spawnedPID, &messages.Connect{Sender: client})

	future := node.ActorEngine.RequestFuture(spawnedPID, message, workerTimeout)
	result, err := future.Result()
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"ip":   ipAddr,
			"peer": p.PeerId,
		}).Error("Error getting result from remote worker")
		handleWorkerError(err, responseCollector)
		return
	}

	logrus.WithFields(logrus.Fields{
		"ip":   ipAddr,
		"peer": p.PeerId,
	}).Info("Successfully processed remote worker response")
	processWorkerResponse(result, p.PeerId, responseCollector)
}

func handleWorkerError(err error, responseCollector chan<- *pubsub2.Message) {
	logrus.Errorf("[-] Error with worker: %v", err)
	responseCollector <- &pubsub2.Message{
		ValidatorData: map[string]interface{}{"error": err.Error()},
	}
}

func processWorkerResponse(result interface{}, workerID interface{}, responseCollector chan<- *pubsub2.Message) {
	response, ok := result.(*messages.Response)
	if !ok {
		logrus.Errorf("[-] Invalid response type from worker")
		return
	}
	msg, err := getResponseMessage(response)
	if err != nil {
		logrus.Errorf("[-] Error converting worker response: %v", err)
		return
	}
	logrus.Infof("Received response from worker %v, sending to responseCollector", workerID)
	responseCollector <- msg
}

func queueResponses(responseCollector <-chan *pubsub2.Message, timeout <-chan time.Time) {
	var responses []*pubsub2.Message
	for {
		select {
		case response := <-responseCollector:
			responses = append(responses, response)
			logrus.Infof("Adding response from %s to responses list. Total responses: %d", response.ReceivedFrom, len(responses))
		case <-timeout:
			logrus.Infof("Timeout reached, sending all responses to workerDoneCh. Total responses: %d", len(responses))
			for _, resp := range responses {
				workerDoneCh <- resp
			}
			return
		}
	}
}
