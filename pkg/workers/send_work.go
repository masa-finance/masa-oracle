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

	success := tryWorkersRoundRobin(node, eligibleWorkers, message, responseCollector)
	if !success {
		logrus.Error("Failed to process the work")
	}
}

func tryWorkersRoundRobin(node *masa.OracleNode, workers []Worker, message *messages.Work, responseCollector chan *pubsub2.Message) bool {
	var localWorker *Worker
	remoteWorkersAttempted := 0

	logrus.Info("Starting round-robin worker selection")

	// Try remote workers first, up to MaxRemoteWorkers
	for _, worker := range workers {
		if !worker.IsLocal {
			if remoteWorkersAttempted >= workerConfig.MaxRemoteWorkers {
				logrus.Infof("Reached maximum remote workers (%d), stopping remote worker attempts", workerConfig.MaxRemoteWorkers)
				break
			}
			remoteWorkersAttempted++
			logrus.Infof("Attempting remote worker %s (attempt %d/%d)", worker.NodeData.PeerId, remoteWorkersAttempted, workerConfig.MaxRemoteWorkers)
			if tryWorker(node, worker, message, responseCollector) {
				logrus.Infof("Remote worker %s successfully completed the work", worker.NodeData.PeerId)
				return true
			}
			logrus.Infof("Remote worker %s failed, moving to next worker", worker.NodeData.PeerId)
		} else {
			localWorker = &worker
			logrus.Info("Found local worker, saving for later if needed")
		}
	}

	// If remote workers fail or don't exist, try local worker
	if localWorker != nil {
		logrus.Info("Attempting local worker")
		return tryWorker(node, *localWorker, message, responseCollector)
	}

	// If no workers are available, create a local worker as last resort
	logrus.Warn("No workers available, creating last resort local worker")
	lastResortLocalWorker := Worker{
		IsLocal:  true,
		NodeData: pubsub.NodeData{PeerId: node.Host.ID()},
		Node:     node,
	}
	return tryWorker(node, lastResortLocalWorker, message, responseCollector)
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
		select {
		case response := <-responseCollector:
			if isSuccessfulResponse(response) {
				if worker.IsLocal {
					logrus.Infof("Local worker with PeerID %s successfully completed the work", node.Host.ID())
				} else {
					logrus.Infof("Remote worker with PeerID %s and IP %s successfully completed the work", worker.NodeData.PeerId, worker.IPAddr)
				}
				processAndSendResponse(response)
				return true
			}
		case <-time.After(workerConfig.WorkerResponseTimeout):
			if worker.IsLocal {
				logrus.Warnf("Local worker with PeerID %s failed to respond in time", node.Host.ID())
			} else {
				logrus.Warnf("Remote worker with PeerID %s and IP %s failed to respond in time", worker.NodeData.PeerId, worker.IPAddr)
			}
		}
	case <-time.After(workerConfig.WorkerTimeout):
		if worker.IsLocal {
			logrus.Warnf("Local worker with PeerID %s timed out", node.Host.ID())
		} else {
			logrus.Warnf("Remote worker with PeerID %s and IP %s timed out", worker.NodeData.PeerId, worker.IPAddr)
		}
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
