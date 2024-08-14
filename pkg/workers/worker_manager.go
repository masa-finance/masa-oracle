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
)

var (
	instance *WorkHandlerManager
	once     sync.Once
)

type WorkRequest struct {
	WorkType  WorkerType
	RequestId string
	Data      []byte
}

type WorkResponse struct {
	WorkRequest  WorkRequest
	Data         interface{}
	Error        error
	WorkerPeerId string
}

func GetWorkHandlerManager() *WorkHandlerManager {
	once.Do(func() {
		instance = &WorkHandlerManager{
			handlers: make(map[WorkerType]*WorkHandlerInfo),
		}
		instance.setupHandlers()
	})
	return instance
}

// ErrHandlerNotFound is an error returned when a work handler cannot be found.
var ErrHandlerNotFound = errors.New("work handler not found")

// WorkHandler defines the interface for handling different types of work.
type WorkHandler interface {
	HandleWork(data []byte) WorkResponse
}

// WorkHandlerInfo contains information about a work handler, including metrics.
type WorkHandlerInfo struct {
	Handler      WorkHandler
	CallCount    int64
	TotalRuntime time.Duration
}

// WorkHandlerManager manages work handlers and tracks their execution metrics.
type WorkHandlerManager struct {
	handlers map[WorkerType]*WorkHandlerInfo
	mu       sync.RWMutex
}

func (whm *WorkHandlerManager) setupHandlers() {
	cfg := config.GetInstance()
	if cfg.TwitterScraper {
		whm.addWorkHandler(Twitter, &TwitterQueryHandler{})
		whm.addWorkHandler(TwitterFollowers, &TwitterFollowersHandler{})
		whm.addWorkHandler(TwitterProfile, &TwitterProfileHandler{})
		whm.addWorkHandler(TwitterSentiment, &TwitterSentimentHandler{})
		whm.addWorkHandler(TwitterTrends, &TwitterTrendsHandler{})
	}
	if cfg.WebScraper {
		whm.addWorkHandler(Web, &WebHandler{})
		whm.addWorkHandler(WebSentiment, &WebSentimentHandler{})
	}
	if cfg.LlmServer {
		whm.addWorkHandler(LLMChat, &LLMChatHandler{})
	}
	if cfg.DiscordScraper {
		whm.addWorkHandler(Discord, &DiscordHandler{})
	}
}

// addWorkHandler registers a new work handler under a specific name.
func (whm *WorkHandlerManager) addWorkHandler(wType WorkerType, handler WorkHandler) {
	whm.mu.Lock()
	defer whm.mu.Unlock()
	whm.handlers[wType] = &WorkHandlerInfo{Handler: handler}
}

// getWorkHandler retrieves a registered work handler by name.
func (whm *WorkHandlerManager) getWorkHandler(wType WorkerType) (WorkHandler, bool) {
	whm.mu.RLock()
	defer whm.mu.RUnlock()
	info, exists := whm.handlers[wType]
	if !exists {
		return nil, false
	}
	return info.Handler, true
}

func (whm *WorkHandlerManager) DistributeWork(node *masa.OracleNode, workRequest WorkRequest) (response WorkResponse) {
	category := WorkerTypeToCategory(workRequest.WorkType)
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
		response := whm.sendWorkToWorker(node, worker, workRequest)
		if response.Error != nil {
			logrus.Errorf("error sending work to worker: %s", response.Error.Error())
			logrus.Infof("Remote worker %s failed, moving to next worker", worker.NodeData.PeerId)
			continue
		}
		return response
	}
	// Fallback to local execution if local worker is eligible
	if localWorker != nil {
		return whm.ExecuteWork(workRequest)
	}
	if response.Error == nil {
		response.Error = errors.New("no eligible workers found")
	}
	return response
}

func (whm *WorkHandlerManager) sendWorkToWorker(node *masa.OracleNode, worker Worker, workRequest WorkRequest) (response WorkResponse) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), workerConfig.WorkerResponseTimeout)
	defer cancel() // Cancel the context when done to release resources

	if err := node.Host.Connect(ctxWithTimeout, *worker.AddrInfo); err != nil {
		response.Error = fmt.Errorf("failed to connect to remote peer %s: %v", worker.AddrInfo.ID.String(), err)
		return
	} else {
		logrus.Debugf("[+] Connection established with node: %s", worker.AddrInfo.ID.String())
		stream, err := node.Host.NewStream(ctxWithTimeout, worker.AddrInfo.ID, config.ProtocolWithVersion(config.WorkerProtocol))
		if err != nil {
			response.Error = fmt.Errorf("error opening stream: %v", err)
			return
		}
		defer func(stream network.Stream) {
			err := stream.Close()
			if err != nil {
				logrus.Errorf("[-] Error closing stream: %s", err)
			}
		}(stream) // Close the stream when done

		// Write the request to the stream with length prefix
		bytes, err := json.Marshal(workRequest)
		if err != nil {
			response.Error = fmt.Errorf("error marshaling work request: %v", err)
			return
		}
		lengthBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBuf, uint32(len(bytes)))
		_, err = stream.Write(lengthBuf)
		if err != nil {
			response.Error = fmt.Errorf("error writing length to stream: %v", err)
			return
		}
		_, err = stream.Write(bytes)
		if err != nil {
			response.Error = fmt.Errorf("error writing to stream: %v", err)
			return
		}

		// Read the response length
		lengthBuf = make([]byte, 4)
		_, err = io.ReadFull(stream, lengthBuf)
		if err != nil {
			response.Error = fmt.Errorf("error reading response length: %v", err)
			return
		}
		responseLength := binary.BigEndian.Uint32(lengthBuf)

		// Read the actual response
		responseBuf := make([]byte, responseLength)
		_, err = io.ReadFull(stream, responseBuf)
		if err != nil {
			response.Error = fmt.Errorf("error reading response: %v", err)
			return
		}
		err = json.Unmarshal(responseBuf, &response)
		if err != nil {
			response.Error = fmt.Errorf("error unmarshaling response: %v", err)
			return
		}
	}
	return response
}

// ExecuteWork finds and executes the work handler associated with the given name.
// It tracks the call count and execution duration for the handler.
func (whm *WorkHandlerManager) ExecuteWork(workRequest WorkRequest) (response WorkResponse) {
	handler, exists := whm.getWorkHandler(workRequest.WorkType)
	if !exists {
		return WorkResponse{Error: ErrHandlerNotFound}
	}

	// Create a context with a 30-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), workerConfig.WorkerResponseTimeout)
	defer cancel()

	// Channel to receive the work response
	responseChan := make(chan WorkResponse, 1)

	// Execute the work in a separate goroutine
	go func() {
		startTime := time.Now()
		workResponse := handler.HandleWork(workRequest.Data)
		if workResponse.Error == nil {
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
		return WorkResponse{Error: errors.New("work execution timed out")}
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

	var workRequest WorkRequest
	err = json.Unmarshal(messageBuf, &workRequest)
	if err != nil {
		logrus.Errorf("error unmarshaling work request: %v", err)
		return
	}
	workResponse := whm.ExecuteWork(workRequest)
	if workResponse.Error != nil {
		logrus.Errorf("error from remote worker %s: executing work: %v", err)
	}
	workResponse.WorkerPeerId = stream.Conn().LocalPeer().String()

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
