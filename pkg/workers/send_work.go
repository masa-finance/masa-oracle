package workers

import (
	"fmt"
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

	responseCollector := make(chan *pubsub2.Message, 1)

	eligibleWorkers := getEligibleWorkers(node, message)
	workerIterator := newRoundRobinIterator(eligibleWorkers)

	for workerIterator.HasNext() {
		worker := workerIterator.Next()

		go func(w Worker) {
			if w.IsLocal {
				handleLocalWorker(node, pid, message, responseCollector)
			} else {
				handleRemoteWorker(node, w.NodeData, w.IPAddr, props, message, responseCollector)
			}
		}(worker)

		select {
		case response := <-responseCollector:
			if isSuccessfulResponse(response) {
				processAndSendResponse(response)
				return
			}
			// If response is not successful, continue to next worker
		case <-time.After(workerTimeout):
			logrus.Warnf("Worker %v timed out, moving to next worker", worker.NodeData.PeerId)
		}
	}

	logrus.Error("All workers failed to process the work")
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

func isSuccessfulResponse(response *pubsub2.Message) bool {
	if response.ValidatorData == nil {
		return true
	}
	validatorData, ok := response.ValidatorData.(map[string]interface{})
	if !ok {
		return false
	}
	errorVal, exists := validatorData["error"]
	return !exists || errorVal == nil
}

func processAndSendResponse(response *pubsub2.Message) {
	logrus.Infof("Processing and sending successful response")
	workerDoneCh <- response
}

func init() {
	workerDoneCh = make(chan *pubsub2.Message, 100)
}
