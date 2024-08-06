package workers

import (
	"fmt"
	"time"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/messages"

	"github.com/asynkron/protoactor-go/actor"

	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

var workerConfig *WorkerConfig

func init() {
	var err error
	workerConfig, err = LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load worker config: %v", err)
	}
	workerDoneCh = make(chan *pubsub2.Message, workerConfig.WorkerBufferSize)
}

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

	eligibleWorkers := GetEligibleWorkers(node, message)

	for retries := 0; retries < workerConfig.MaxRetries; retries++ {
		success := tryWorkersRoundRobin(node, eligibleWorkers, message, responseCollector)
		if success {
			return
		}
		logrus.Warnf("All workers failed, retry attempt %d of %d", retries+1, workerConfig.MaxRetries)
	}

	logrus.Error("All workers failed to process the work after maximum retries")
}

func tryWorkersRoundRobin(node *masa.OracleNode, workers []Worker, message *messages.Work, responseCollector chan *pubsub2.Message) bool {
	workerIterator := NewRoundRobinIterator(workers)

	for workerIterator.HasNext() {
		worker := workerIterator.Next()

		success := tryWorker(node, worker, message, responseCollector)
		if success {
			return true
		}
	}

	return false
}

func tryWorker(node *masa.OracleNode, worker Worker, message *messages.Work, responseCollector chan *pubsub2.Message) bool {
	workerDone := make(chan bool, 1)

	go func() {
		if worker.IsLocal {
			handleLocalWorker(node, node.ActorEngine.Spawn(actor.PropsFromProducer(NewWorker(node))), message, responseCollector)
		} else {
			handleRemoteWorker(node, worker.NodeData, worker.IPAddr, actor.PropsFromProducer(NewWorker(node)), message, responseCollector)
		}
		workerDone <- true
	}()

	select {
	case <-workerDone:
		// Worker finished, check response
		select {
		case response := <-responseCollector:
			if isSuccessfulResponse(response) {
				processAndSendResponse(response)
				return true
			}
		default:
			// No response in channel, continue to next worker
		}
	case <-time.After(workerConfig.WorkerTimeout):
		logrus.Warnf("Worker %v timed out", worker.NodeData.PeerId)
	}

	return false
}

func createWorkMessage(m *pubsub2.Message, pid *actor.PID) *messages.Work {
	return &messages.Work{
		Data:   string(m.Data),
		Sender: pid,
		Id:     m.ReceivedFrom.String(),
		Type:   int64(pubsub.CategoryTwitter),
	}
}

func handleLocalWorker(node *masa.OracleNode, pid *actor.PID, message *messages.Work, responseCollector chan<- *pubsub2.Message) {
	logrus.Info("Sending work to local worker")
	future := node.ActorEngine.RequestFuture(pid, message, workerConfig.WorkerTimeout)
	result, err := future.Result()
	if err != nil {
		handleWorkerError(err, responseCollector)
		return
	}

	// Log the full response from the local worker
	logrus.WithField("full_response", result).Info("Full response from local worker")

	processWorkerResponse(result, node.Host.ID(), responseCollector)
}

func handleRemoteWorker(node *masa.OracleNode, p pubsub.NodeData, ipAddr string, props *actor.Props, message *messages.Work, responseCollector chan<- *pubsub2.Message) {
	logrus.WithFields(logrus.Fields{
		"ip":   ipAddr,
		"peer": p.PeerId,
	}).Info("Handling remote worker")

	var spawned *actor.PID
	var err error

	// Attempt to spawn the remote worker multiple times
	for attempt := 1; attempt <= workerConfig.MaxSpawnAttempts; attempt++ {
		spawned, err = spawnRemoteWorker(node, ipAddr)
		if err == nil {
			break
		}
		logrus.WithError(err).WithFields(logrus.Fields{
			"ip":      ipAddr,
			"attempt": attempt,
		}).Warn("Failed to spawn remote worker, retrying")
		time.Sleep(time.Second * time.Duration(attempt)) // Exponential backoff
	}

	if err != nil {
		logrus.WithError(err).WithField("ip", ipAddr).Error("Failed to spawn remote worker after multiple attempts")
		handleWorkerError(err, responseCollector)
		return
	}

	client := node.ActorEngine.Spawn(props)
	node.ActorEngine.Send(spawned, &messages.Connect{Sender: client})

	future := node.ActorEngine.RequestFuture(spawned, message, workerConfig.WorkerTimeout)
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

func spawnRemoteWorker(node *masa.OracleNode, ipAddr string) (*actor.PID, error) {
	spawned, err := node.ActorRemote.SpawnNamed(fmt.Sprintf("%s:4001", ipAddr), "worker", "peer", -1)
	if err != nil {
		return nil, err
	}

	if spawned == nil || spawned.Pid == nil {
		return nil, fmt.Errorf("failed to spawn remote worker: PID is nil for IP %s", ipAddr)
	}

	return spawned.Pid, nil
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
