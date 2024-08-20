package workers

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/event"
	"github.com/masa-finance/masa-oracle/pkg/workers/handlers"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

var (
	instance *WorkHandlerManager
	once     sync.Once
)

func GetWorkHandlerManager() *WorkHandlerManager {
	once.Do(func() {
		instance = &WorkHandlerManager{
			handlers:     make(map[data_types.WorkerType]*WorkHandlerInfo),
			eventTracker: event.NewEventTracker(nil),
		}
		instance.setupHandlers()
	})
	return instance
}

// ErrHandlerNotFound is an error returned when a work handler cannot be found.
var ErrHandlerNotFound = errors.New("work handler not found")

// WorkHandler defines the interface for handling different types of work.
type WorkHandler interface {
	HandleWork(data []byte) data_types.WorkResponse
}

// WorkHandlerInfo contains information about a work handler, including metrics.
type WorkHandlerInfo struct {
	Handler      WorkHandler
	CallCount    int64
	TotalRuntime time.Duration
}

// WorkHandlerManager manages work handlers and tracks their execution metrics.
type WorkHandlerManager struct {
	handlers     map[data_types.WorkerType]*WorkHandlerInfo
	mu           sync.RWMutex
	eventTracker *event.EventTracker
}

func (whm *WorkHandlerManager) setupHandlers() {
	cfg := config.GetInstance()
	if cfg.TwitterScraper {
		whm.addWorkHandler(data_types.Twitter, &handlers.TwitterQueryHandler{})
		whm.addWorkHandler(data_types.TwitterFollowers, &handlers.TwitterFollowersHandler{})
		whm.addWorkHandler(data_types.TwitterProfile, &handlers.TwitterProfileHandler{})
		whm.addWorkHandler(data_types.TwitterSentiment, &handlers.TwitterSentimentHandler{})
		whm.addWorkHandler(data_types.TwitterTrends, &handlers.TwitterTrendsHandler{})
	}
	if cfg.WebScraper {
		whm.addWorkHandler(data_types.Web, &handlers.WebHandler{})
		whm.addWorkHandler(data_types.WebSentiment, &handlers.WebSentimentHandler{})
	}
	if cfg.LlmServer {
		whm.addWorkHandler(data_types.LLMChat, &handlers.LLMChatHandler{})
	}
	if cfg.DiscordScraper {
		whm.addWorkHandler(data_types.Discord, &handlers.DiscordProfileHandler{})
	}
}

// addWorkHandler registers a new work handler under a specific name.
func (whm *WorkHandlerManager) addWorkHandler(wType data_types.WorkerType, handler WorkHandler) {
	whm.mu.Lock()
	defer whm.mu.Unlock()
	whm.handlers[wType] = &WorkHandlerInfo{Handler: handler}
}

// getWorkHandler retrieves a registered work handler by name.
func (whm *WorkHandlerManager) getWorkHandler(wType data_types.WorkerType) (WorkHandler, bool) {
	whm.mu.RLock()
	defer whm.mu.RUnlock()
	info, exists := whm.handlers[wType]
	if !exists {
		return nil, false
	}
	return info.Handler, true
}

func (whm *WorkHandlerManager) DistributeWork(node *masa.OracleNode, workRequest data_types.WorkRequest) (response data_types.WorkResponse) {
	category := data_types.WorkerTypeToCategory(workRequest.WorkType)
	remoteWorkers, localWorker := GetEligibleWorkers(node, category, workerConfig)

	remoteWorkersAttempted := 0
	logrus.Info("Starting round-robin worker selection")

	// Try remote workers first, up to MaxRemoteWorkers
	for _, worker := range remoteWorkers {
		if remoteWorkersAttempted >= workerConfig.MaxRemoteWorkers {
			logrus.Infof("Reached maximum remote workers (%d), stopping remote worker attempts", workerConfig.MaxRemoteWorkers)
			break
		}
		remoteWorkersAttempted++
		logrus.Infof("Attempting remote worker %s (attempt %d/%d)", worker.NodeData.PeerId, remoteWorkersAttempted, workerConfig.MaxRemoteWorkers)
		response = whm.sendWorkToWorker(node, worker, workRequest)
		if response.Error != "" {
			whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
			logrus.Errorf("error sending work to worker: %s: %s", response.WorkerPeerId, response.Error)
			logrus.Infof("Remote worker %s failed, moving to next worker", worker.NodeData.PeerId)
			continue
		}
		whm.eventTracker.TrackWorkCompletion(workRequest.WorkType, response.Error == "", worker.AddrInfo.ID.String())
		return response
	}
	// Fallback to local execution if local worker is eligible
	if localWorker != nil {
		var reason string
		if len(remoteWorkers) > 0 {
			reason = "all remote workers failed"
		} else {
			reason = "no remote workers available"
		}
		whm.eventTracker.TrackLocalWorkerFallback(reason)
		whm.eventTracker.TrackWorkExecutionStart(workRequest.WorkType, false, localWorker.AddrInfo.ID.String())
		return whm.ExecuteWork(workRequest)
	}
	if response.Error == "" {
		response.Error = "no eligible workers found"
	} else {
		response.Error = fmt.Sprintf("no workers could process: remote attempt failed due to: %s", response.Error)
	}
	return response
}

func (whm *WorkHandlerManager) sendWorkToWorker(node *masa.OracleNode, worker data_types.Worker, workRequest data_types.WorkRequest) (response data_types.WorkResponse) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), workerConfig.WorkerResponseTimeout)
	defer cancel() // Cancel the context when done to release resources

	if err := node.Host.Connect(ctxWithTimeout, *worker.AddrInfo); err != nil {
		response.Error = fmt.Sprintf("failed to connect to remote peer %s: %v", worker.AddrInfo.ID.String(), err)
		whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
		return
	} else {
		whm.eventTracker.TrackRemoteWorkerConnection(worker.AddrInfo.ID.String())
		logrus.Debugf("[+] Connection established with node: %s", worker.AddrInfo.ID.String())
		protocol := config.ProtocolWithVersion(config.WorkerProtocol)
		stream, err := node.Host.NewStream(ctxWithTimeout, worker.AddrInfo.ID, protocol)
		if err != nil {
			response.Error = fmt.Sprintf("error opening stream: %v", err)
			whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
			return
		}
		// the stream should be closed by the receiver, but keeping this here just in case
		whm.eventTracker.TrackStreamCreation(worker.AddrInfo.ID.String(), string(protocol))
		defer func(stream network.Stream) {
			err := stream.Close()
			if err != nil {
				logrus.Debugf("[-] Error closing stream: %s", err)
			}
		}(stream) // Close the stream when done.S

		// Write the request to the stream with length prefix
		bytes, err := json.Marshal(workRequest)
		if err != nil {
			response.Error = fmt.Sprintf("error marshaling work request: %v", err)
			return
		}
		lengthBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBuf, uint32(len(bytes)))
		_, err = stream.Write(lengthBuf)
		if err != nil {
			response.Error = fmt.Sprintf("error writing length to stream: %v", err)
			return
		}
		whm.eventTracker.TrackWorkRequestSerialization(workRequest.WorkType, len(bytes))
		_, err = stream.Write(bytes)
		if err != nil {
			response.Error = fmt.Sprintf("error writing to stream: %v", err)
			whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
			return
		}
		whm.eventTracker.TrackWorkDistribution(workRequest.WorkType, true, worker.AddrInfo.ID.String())
		// Read the response length
		lengthBuf = make([]byte, 4)
		_, err = io.ReadFull(stream, lengthBuf)
		if err != nil {
			response.Error = fmt.Sprintf("error reading response length: %v", err)
			whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
			return
		}
		responseLength := binary.BigEndian.Uint32(lengthBuf)

		// Read the actual response
		responseBuf := make([]byte, responseLength)
		_, err = io.ReadFull(stream, responseBuf)
		if err != nil {
			response.Error = fmt.Sprintf("error reading response: %v", err)
			return
		}
		err = json.Unmarshal(responseBuf, &response)
		if err != nil {
			response.Error = fmt.Sprintf("error unmarshaling response: %v", err)
			return
		}
		whm.eventTracker.TrackWorkResponseDeserialization(workRequest.WorkType, true)
	}
	return response
}

// ExecuteWork finds and executes the work handler associated with the given name.
// It tracks the call count and execution duration for the handler.
func (whm *WorkHandlerManager) ExecuteWork(workRequest data_types.WorkRequest) (response data_types.WorkResponse) {
	handler, exists := whm.getWorkHandler(workRequest.WorkType)
	if !exists {
		return data_types.WorkResponse{Error: ErrHandlerNotFound.Error()}
	}

	// Create a context with a 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), workerConfig.WorkerResponseTimeout)
	defer cancel()

	// Channel to receive the work response
	responseChan := make(chan data_types.WorkResponse, 1)

	// Execute the work in a separate goroutine
	go func() {
		startTime := time.Now()
		workResponse := handler.HandleWork(workRequest.Data)
		if workResponse.Error == "" {
			duration := time.Since(startTime)
			whm.mu.Lock()
			handlerInfo := whm.handlers[workRequest.WorkType]
			handlerInfo.CallCount++
			handlerInfo.TotalRuntime += duration
			whm.mu.Unlock()
		}
		responseChan <- workResponse
	}()

	select {
	case <-ctx.Done():
		// Context timed out
		whm.eventTracker.TrackWorkExecutionTimeout(workRequest.WorkType, workerConfig.WorkerResponseTimeout)

		return data_types.WorkResponse{Error: "work execution timed out"}
	case response = <-responseChan:
		// Work completed within the timeout
		return response
	}
}

func (whm *WorkHandlerManager) HandleWorkerStream(stream network.Stream) {
	defer func(stream network.Stream) {
		err := stream.Close()
		if err != nil {
			logrus.Errorf("[-] Error closing stream in handler: %s", err)
		}
	}(stream)

	// Read the length of the message
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(stream, lengthBuf)
	if err != nil {
		logrus.Errorf("error reading message length: %v", err)
		return
	}
	messageLength := binary.BigEndian.Uint32(lengthBuf)

	// Read the actual message
	messageBuf := make([]byte, messageLength)
	_, err = io.ReadFull(stream, messageBuf)
	if err != nil {
		logrus.Errorf("error reading message: %v", err)
		return
	}

	var workRequest data_types.WorkRequest
	err = json.Unmarshal(messageBuf, &workRequest)
	if err != nil {
		logrus.Errorf("error unmarshaling work request: %v", err)
		return
	}
	peerId := stream.Conn().LocalPeer().String()
	whm.eventTracker.TrackWorkExecutionStart(workRequest.WorkType, true, peerId)
	workResponse := whm.ExecuteWork(workRequest)
	if workResponse.Error != "" {
		logrus.Errorf("error from remote worker %s: executing work: %s", peerId, workResponse.Error)
	}
	workResponse.WorkerPeerId = peerId

	// Write the response to the stream
	responseBytes, err := json.Marshal(workResponse)
	if err != nil {
		logrus.Errorf("error marshaling work response: %v", err)
		return
	}

	// Prefix the response with its length
	responseLength := uint32(len(responseBytes))
	lengthBuf = make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, responseLength)

	_, err = stream.Write(lengthBuf)
	if err != nil {
		logrus.Errorf("error writing response length to stream: %v", err)
		return
	}

	_, err = stream.Write(responseBytes)
	if err != nil {
		logrus.Errorf("error writing response to stream: %v", err)
		return
	}
}
