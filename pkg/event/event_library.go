package event

import (
	"time"

	"github.com/sirupsen/logrus"

	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
)

// TrackWorkDistribution records the distribution of work to a worker.
//
// Parameters:
// - workType: The type of work being distributed (e.g., Twitter, Web, Discord)
// - remoteWorker: Boolean indicating if the work is sent to a remote worker (true) or executed locally (false)
// - peerId: String containing the peer ID
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
// - "work_type": The WorkerType as a string
// - "remote_worker": Boolean indicating if it's a remote worker
func (a *EventTracker) TrackWorkDistribution(workType data_types.WorkerType, remoteWorker bool, peerId string) {
	err := a.TrackAndSendEvent(WorkDistribution, map[string]interface{}{
		"peer_id":       peerId,
		"work_type":     workType,
		"remote_worker": remoteWorker,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking work distribution event: %s", err)
	}
}

// TrackWorkCompletion records the completion of a work item.
//
// Parameters:
// - workType: The type of work that was completed
// - success: Boolean indicating if the work was completed successfully
// - peerId: String containing the peer ID
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
// - "work_type": The WorkerType as a string
// - "success": Boolean indicating if the work was successful
func (a *EventTracker) TrackWorkCompletion(workType data_types.WorkerType, success bool, peerId string) {
	err := a.TrackAndSendEvent(WorkCompletion, map[string]interface{}{
		"peer_id":   peerId,
		"work_type": workType,
		"success":   success,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking work completion event: %s", err)
	}
}

// TrackWorkerFailure records a failure that occurred during work execution.
//
// Parameters:
// - workType: The type of work that failed
// - errorMessage: A string describing the error that occurred
// - peerId: String containing the peer ID
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
// - "work_type": The WorkerType as a string
// - "error": String containing the error message
func (a *EventTracker) TrackWorkerFailure(workType data_types.WorkerType, errorMessage string, peerId string) {
	err := a.TrackAndSendEvent(WorkFailure, map[string]interface{}{
		"peer_id":   peerId,
		"work_type": workType,
		"error":     errorMessage,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking worker failure event: %s", err)
	}
}

// TrackWorkExecutionStart records the start of work execution.
//
// Parameters:
// - workType: The type of work being executed
// - remoteWorker: Boolean indicating if the work is executed by a remote worker (true) or locally (false)
// - peerId: String containing the peer ID
//
// The event will contain the following data:
// - "work_type": The WorkerType as a string
// - "remote_worker": Boolean indicating if it's a remote worker
// - "peer_id": String containing the peer ID
func (a *EventTracker) TrackWorkExecutionStart(workType data_types.WorkerType, remoteWorker bool, peerId string) {
	err := a.TrackAndSendEvent(WorkExecutionStart, map[string]interface{}{
		"work_type":     workType,
		"remote_worker": remoteWorker,
		"peer_id":       peerId,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking work execution start event: %s", err)
	}
}

// TrackWorkExecutionTimeout records when work execution times out.
//
// Parameters:
// - workType: The type of work that timed out
// - timeoutDuration: The duration of the timeout
//
// The event will contain the following data:
// - "work_type": The WorkerType as a string
// - "timeout_duration": The duration of the timeout
func (a *EventTracker) TrackWorkExecutionTimeout(workType data_types.WorkerType, timeoutDuration time.Duration) {
	err := a.TrackAndSendEvent(WorkExecutionTimeout, map[string]interface{}{
		"work_type":        workType,
		"timeout_duration": timeoutDuration,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking work execution timeout event: %s", err)
	}
}

// TrackRemoteWorkerConnection records when a connection is established with a remote worker.
//
// Parameters:
// - peerId: String containing the peer ID
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
func (a *EventTracker) TrackRemoteWorkerConnection(peerId string) {
	err := a.TrackAndSendEvent(RemoteWorkerConnection, map[string]interface{}{
		"peer_id": peerId,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking remote worker connection event: %s", err)
	}
}

// TrackStreamCreation records when a new stream is created for communication with a remote worker.
//
// Parameters:
// - peerId: String containing the peer ID
// - protocol: The protocol used for the stream
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
// - "protocol": The protocol used for the stream
func (a *EventTracker) TrackStreamCreation(peerId string, protocol string) {
	err := a.TrackAndSendEvent(StreamCreation, map[string]interface{}{
		"peer_id":  peerId,
		"protocol": protocol,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking stream creation event: %s", err)
	}
}

// TrackWorkRequestSerialization records when a work request is serialized for transmission.
//
// Parameters:
// - workType: The type of work being serialized
// - dataSize: The size of the serialized data
//
// The event will contain the following data:
// - "work_type": The WorkerType as a string
// - "data_size": The size of the serialized data
func (a *EventTracker) TrackWorkRequestSerialization(workType data_types.WorkerType, dataSize int) {
	err := a.TrackAndSendEvent(WorkRequestSerialization, map[string]interface{}{
		"work_type": workType,
		"data_size": dataSize,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking work request serialization event: %s", err)
	}
}

// TrackWorkResponseDeserialization records when a work response is deserialized after reception.
//
// Parameters:
// - workType: The type of work being deserialized
// - success: Boolean indicating if the deserialization was successful
//
// The event will contain the following data:
// - "work_type": The WorkerType as a string
// - "success": Boolean indicating if the deserialization was successful
func (a *EventTracker) TrackWorkResponseDeserialization(workType data_types.WorkerType, success bool) {
	err := a.TrackAndSendEvent(WorkResponseDeserialization, map[string]interface{}{
		"work_type": workType,
		"success":   success,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking work response deserialization event: %s", err)
	}
}

// TrackLocalWorkerFallback records when the system falls back to using a local worker.
//
// Parameters:
// - reason: The reason for the fallback
//
// The event will contain the following data:
// - "reason": The reason for the fallback
func (a *EventTracker) TrackLocalWorkerFallback(reason string) {
	err := a.TrackAndSendEvent(LocalWorkerFallback, map[string]interface{}{
		"reason": reason,
	}, nil)
	if err != nil {
		logrus.Errorf("error tracking local worker fallback event: %s", err)
	}
}
