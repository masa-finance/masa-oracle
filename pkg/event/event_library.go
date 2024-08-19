package event

import (
	data_types "github.com/masa-finance/masa-oracle/pkg/workers/types"
	"github.com/sirupsen/logrus"
)

// TrackWorkDistribution records the distribution of work to a worker.
//
// Parameters:
// - "peer_id": String containing the peer ID
// - workType: The type of work being distributed (e.g., Twitter, Web, Discord)
// - remoteWorker: Boolean indicating if the work is sent to a remote worker (true) or executed locally (false)
// - client: The EventClient used to send the event
//
// Returns:
// - error: nil if the event was successfully tracked and sent, otherwise an error describing what went wrong
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
// - "work_type": The WorkerType as a string
// - "remote_worker": Boolean indicating if it's a remote worker
func (a *EventTracker) TrackWorkDistribution(workType data_types.WorkerType, remoteWorker bool, peerId string, client *EventClient) error {
	if a == nil {
		logrus.Error("EventTracker is nil")
		return nil
	}
	return a.TrackAndSendEvent("work_distribution", map[string]interface{}{
		"peer_id":       peerId,
		"work_type":     workType,
		"remote_worker": remoteWorker,
	}, client)
}

// TrackWorkCompletion records the completion of a work item.
//
// Parameters:
// - "peer_id": String containing the peer ID
// - workType: The type of work that was completed
// - success: Boolean indicating if the work was completed successfully
// - client: The EventClient used to send the event
//
// Returns:
// - error: nil if the event was successfully tracked and sent, otherwise an error describing what went wrong
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
// - "work_type": The WorkerType as a string
// - "success": Boolean indicating if the work was successful
func (a *EventTracker) TrackWorkCompletion(workType data_types.WorkerType, success bool, peerId string, client *EventClient) error {
	if a == nil {
		logrus.Error("EventTracker is nil")
		return nil
	}
	return a.TrackAndSendEvent("work_completion", map[string]interface{}{
		"peer_id":   peerId,
		"work_type": workType,
		"success":   success,
	}, client)
}

// TrackWorkerFailure records a failure that occurred during work execution.
//
// Parameters:
// - "peer_id": String containing the peer ID
// - workType: The type of work that failed
// - errorMessage: A string describing the error that occurred
// - client: The EventClient used to send the event
//
// Returns:
// - error: nil if the event was successfully tracked and sent, otherwise an error describing what went wrong
//
// The event will contain the following data:
// - "peer_id": String containing the peer ID
// - "work_type": The WorkerType as a string
// - "error": String containing the error message
func (a *EventTracker) TrackWorkerFailure(workType data_types.WorkerType, errorMessage string, peerId string, client *EventClient) error {
	if a == nil {
		logrus.Error("EventTracker is nil")
		return nil
	}
	return a.TrackAndSendEvent("worker_failure", map[string]interface{}{
		"peer_id":   peerId,
		"work_type": workType,
		"error":     errorMessage,
	}, client)
}
