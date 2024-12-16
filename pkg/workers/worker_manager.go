package workers

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/node"
	"github.com/masa-finance/masa-oracle/pkg/event"
	"github.com/masa-finance/masa-oracle/pkg/pubsub"
	"github.com/masa-finance/masa-oracle/pkg/workers/handlers"
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

func NewWorkHandlerManager(opts ...WorkerOptionFunc) *WorkHandlerManager {
	options := &WorkerOption{}
	options.Apply(opts...)

	whm := &WorkHandlerManager{
		handlers:     make(map[data_types.WorkerType]*WorkHandlerInfo),
		eventTracker: event.NewEventTracker(nil),
	}

	if options.isTwitterWorker {
		whm.addWorkHandler(data_types.Twitter, &handlers.TwitterQueryHandler{MasaDir: options.masaDir})
		whm.addWorkHandler(data_types.TwitterFollowers, &handlers.TwitterFollowersHandler{MasaDir: options.masaDir})
		whm.addWorkHandler(data_types.TwitterProfile, &handlers.TwitterProfileHandler{MasaDir: options.masaDir})
	}

	if options.isWebScraperWorker {
		whm.addWorkHandler(data_types.Web, &handlers.WebHandler{})
	}

	return whm
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

func (whm *WorkHandlerManager) DistributeWork(node *node.OracleNode, workRequest data_types.WorkRequest) (response data_types.WorkResponse) {
	category := data_types.WorkerTypeToCategory(workRequest.WorkType)
	var remoteWorkers []data_types.Worker
	var localWorker *data_types.Worker

	if category == pubsub.CategoryTwitter {
		// Use priority-based selection for Twitter work
		remoteWorkers, localWorker = GetEligibleWorkers(node, category, workerConfig.MaxRemoteWorkers)
		logrus.Info("Starting priority-based worker selection for Twitter work")
	} else {
		// Use existing selection for other work types
		remoteWorkers, localWorker = GetEligibleWorkers(node, category, 0)
		// Shuffle the workers to maintain round-robin behavior
		rand.Shuffle(len(remoteWorkers), func(i, j int) {
			remoteWorkers[i], remoteWorkers[j] = remoteWorkers[j], remoteWorkers[i]
		})
		logrus.Info("Starting round-robin worker selection for non-Twitter work")
	}

	remoteWorkersAttempted := 0
	var errorList []string

	// Try remote workers first, up to MaxRemoteWorkers
	for _, worker := range remoteWorkers {
		if remoteWorkersAttempted >= workerConfig.MaxRemoteWorkers {
			logrus.Infof("Reached maximum remote workers (%d), stopping remote worker attempts", workerConfig.MaxRemoteWorkers)
			break
		}
		remoteWorkersAttempted++

		// Attempt to connect to the worker
		ctx, cancel := context.WithTimeout(context.Background(), workerConfig.FindPeerTimeout)
		peerInfo, err := node.DHT.FindPeer(ctx, worker.NodeData.PeerId)
		cancel()
		if err != nil {
			if err == context.DeadlineExceeded {
				logrus.Warnf("Timeout while finding peer %s in DHT", worker.NodeData.PeerId.String())
			} else {
				logrus.Warnf("Failed to find peer %s in DHT: %v", worker.NodeData.PeerId.String(), err)
			}
			if category == pubsub.CategoryTwitter {
				err := node.NodeTracker.UpdateNodeDataTwitter(worker.NodeData.PeerId.String(), pubsub.NodeData{
					LastNotFoundTime: time.Now(),
					NotFoundCount:    1,
				})
				if err != nil {
					logrus.Warnf("Failed to update node data for peer %s: %v", worker.NodeData.PeerId.String(), err)
				}
			}
			continue
		}

		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), workerConfig.ConnectionTimeout)
		err = node.Host.Connect(ctxWithTimeout, peerInfo)
		cancel()
		if err != nil {
			logrus.Warnf("Failed to connect to peer %s: %v", worker.NodeData.PeerId.String(), err)
			continue
		}

		worker.AddrInfo = &peerInfo

		logrus.Infof("Attempting remote worker %s (attempt %d/%d)", worker.NodeData.PeerId, remoteWorkersAttempted, workerConfig.MaxRemoteWorkers)
		response = whm.sendWorkToWorker(node, worker, workRequest)
		if response.Error != "" {
			errorMsg := fmt.Sprintf("Worker %s: %s", worker.NodeData.PeerId, response.Error)
			errorList = append(errorList, errorMsg)

			whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
			logrus.Errorf("error sending work to worker: %s: %s", response.WorkerPeerId, response.Error)
			logrus.Infof("Remote worker %s failed, moving to next worker", worker.NodeData.PeerId)

			// Check if the error is related to Twitter authentication
			if strings.Contains(response.Error, "unable to get twitter profile: there was an error authenticating with your Twitter credentials") {
				logrus.Warnf("Worker %s failed due to Twitter authentication error. Skipping to the next worker.", worker.NodeData.PeerId)
				continue
			}
		} else {
			return response
		}
	}

	// Fallback to local execution if local worker is eligible and all remote workers failed
	if localWorker != nil {
		var reason string
		if len(remoteWorkers) > 0 {
			reason = "all remote workers failed"
		} else {
			reason = "no remote workers available"
		}
		whm.eventTracker.TrackLocalWorkerFallback(workRequest.WorkType, reason, localWorker.AddrInfo.ID.String())

		response = whm.ExecuteWork(workRequest)
		whm.eventTracker.TrackWorkCompletion(workRequest.WorkType, response.Error == "", localWorker.AddrInfo.ID.String())

		if response.Error != "" {
			errorList = append(errorList, fmt.Sprintf("Local worker: %s", response.Error))
		} else {
			return response
		}
	}

	// If we reach here, all attempts failed
	if len(errorList) == 0 {
		response.Error = "no eligible workers found"
	} else {
		response.Error = fmt.Sprintf("All workers failed. Errors: %s", strings.Join(errorList, "; "))
	}
	return response
}

func (whm *WorkHandlerManager) sendWorkToWorker(node *node.OracleNode, worker data_types.Worker, workRequest data_types.WorkRequest) (response data_types.WorkResponse) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), workerConfig.WorkerResponseTimeout)
	defer cancel() // Cancel the context when done to release resources

	if err := node.Host.Connect(ctxWithTimeout, *worker.AddrInfo); err != nil {
		response.Error = fmt.Sprintf("failed to connect to remote peer %s: %v", worker.AddrInfo.ID.String(), err)
		whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
		return
	} else {
		//whm.eventTracker.TrackRemoteWorkerConnection(worker.AddrInfo.ID.String())
		logrus.Debugf("[+] Connection established with node: %s", worker.AddrInfo.ID.String())
		stream, err := node.ProtocolStream(ctxWithTimeout, worker.AddrInfo.ID, node.Options.WorkerProtocol)
		if err != nil {
			response.Error = fmt.Sprintf("error opening stream: %v", err)
			whm.eventTracker.TrackWorkerFailure(workRequest.WorkType, response.Error, worker.AddrInfo.ID.String())
			return
		}
		// the stream should be closed by the receiver, but keeping this here just in case
		defer func(stream network.Stream) {
			err := stream.Close()
			if err != nil {
				logrus.Debugf("[-] Error closing stream: %s", err)
			}
		}(stream) // Close the stream when done

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
		// Update metrics only if the work category is Twitter
		if data_types.WorkerTypeToCategory(workRequest.WorkType) == pubsub.CategoryTwitter {
			if response.Error == "" {
				err = node.NodeTracker.UpdateNodeDataTwitter(worker.NodeData.PeerId.String(), pubsub.NodeData{
					LastReturnedTweet: time.Now(),
				})
			} else {
				err = node.NodeTracker.UpdateNodeDataTwitter(worker.NodeData.PeerId.String(), pubsub.NodeData{
					TweetTimeout:     true,
					TweetTimeouts:    1,
					LastTweetTimeout: time.Now(),
				})
			}
			if err != nil {
				logrus.Warnf("Failed to update node data for peer %s: %v", worker.NodeData.PeerId.String(), err)
			}
		}
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
		duration := time.Since(startTime)
		whm.mu.Lock()
		handlerInfo := whm.handlers[workRequest.WorkType]
		handlerInfo.CallCount++
		handlerInfo.TotalRuntime += duration
		whm.mu.Unlock()

		if workResponse.Error != "" {
			logrus.Errorf("[-] Work error for %s: %s", workRequest.WorkType, workResponse.Error)
		} else if workResponse.Data == "" {
			logrus.Warnf("[-] Work response for %s: No data returned", workRequest.WorkType)
		}

		responseChan <- workResponse
	}()

	select {
	case <-ctx.Done():
		// Context timed out
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
			logrus.Infof("[-] Error closing stream in handler: %s", err)
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
	workResponse := whm.ExecuteWork(workRequest)
	if workResponse.Error != "" {
		logrus.Errorf("error from remote worker %s: executing work: %s", peerId, workResponse.Error)
	}
	workResponse.WorkerPeerId = peerId
	whm.eventTracker.TrackWorkCompletion(workRequest.WorkType, workResponse.Error == "", peerId)

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
