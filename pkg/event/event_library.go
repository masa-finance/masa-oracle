package event

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Add this function to the existing file

// TrackWorkRequest records when a work request is initiated.
//
// Parameters:
// - workType: String indicating the type of work being requested (e.g., "SearchTweetsRecent")
// - peerId: String containing the peer ID (or client IP in this case)
func (a *EventTracker) TrackWorkRequest(workType string, peerId string, payload string, dataSource string) {
	event := Event{
		Name:       WorkRequest,
		PeerID:     peerId,
		Payload:    payload,
		DataSource: dataSource,
		WorkType:   workType,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking work request event: %s", err)
	}

	logrus.Infof("[+] %s input: %s", workType, payload)
}

// TrackWorkDistribution records the distribution of work to a worker.
//
// Parameters:
// - remoteWorker: Boolean indicating if the work is sent to a remote worker (true) or executed locally (false)
// - peerId: String containing the peer ID
func (a *EventTracker) TrackWorkDistribution(remoteWorker bool, peerId string) {
	event := Event{
		Name:         WorkDistribution,
		PeerID:       peerId,
		WorkType:     WorkDistribution,
		RemoteWorker: remoteWorker,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking work distribution event: %s", err)
	}
}

// TrackWorkCompletion records the completion of a work item.
//
// Parameters:
// - success: Boolean indicating if the work was completed successfully
// - peerId: String containing the peer ID
func (a *EventTracker) TrackWorkCompletion(success bool, peerId string) {
	event := Event{
		Name:     WorkCompletion,
		PeerID:   peerId,
		WorkType: WorkCompletion,
		Success:  success,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking work completion event: %s", err)
	}
}

// TrackWorkerFailure records a failure that occurred during work execution.
//
// Parameters:
// - errorMessage: A string describing the error that occurred
// - peerId: String containing the peer ID
func (a *EventTracker) TrackWorkerFailure(errorMessage string, peerId string) {
	event := Event{
		Name:     WorkFailure,
		PeerID:   peerId,
		WorkType: WorkFailure,
		Error:    errorMessage,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking worker failure event: %s", err)
	}
}

// TrackWorkExecutionStart records the start of work execution.
//
// Parameters:
// - remoteWorker: Boolean indicating if the work is executed by a remote worker (true) or locally (false)
// - peerId: String containing the peer ID
func (a *EventTracker) TrackWorkExecutionStart(remoteWorker bool, peerId string) {
	event := Event{
		Name:         WorkExecutionStart,
		PeerID:       peerId,
		WorkType:     WorkExecutionStart,
		RemoteWorker: remoteWorker,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking work execution start event: %s", err)
	}
}

// TrackWorkExecutionTimeout records when work execution times out.
//
// Parameters:
// - timeoutDuration: The duration of the timeout
// - peerId: String containing the peer ID
func (a *EventTracker) TrackWorkExecutionTimeout(timeoutDuration time.Duration, peerId string) {
	event := Event{
		Name:     WorkExecutionTimeout,
		PeerID:   peerId,
		WorkType: WorkExecutionTimeout,
		Error:    fmt.Sprintf("timeout after %s", timeoutDuration),
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking work execution timeout event: %s", err)
	}
}

// TrackRemoteWorkerConnection records when a connection is established with a remote worker.
//
// Parameters:
// - peerId: String containing the peer ID
func (a *EventTracker) TrackRemoteWorkerConnection(peerId string) {
	event := Event{
		Name:     RemoteWorkerConnection,
		PeerID:   peerId,
		WorkType: RemoteWorkerConnection,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking remote worker connection event: %s", err)
	}
}

// TrackStreamCreation records when a new stream is created for communication with a remote worker.
//
// Parameters:
// - peerId: String containing the peer ID
// - protocol: The protocol used for the stream
func (a *EventTracker) TrackStreamCreation(peerId string, protocol string) {
	event := Event{
		Name:     StreamCreation,
		PeerID:   peerId,
		WorkType: StreamCreation,
		Error:    protocol, // Assuming protocol is stored in Error field for now

	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking stream creation event: %s", err)
	}
}

// TrackWorkRequestSerialization records when a work request is serialized for transmission.
//
// Parameters:
// - dataSize: The size of the serialized data
// - peerId: String containing the peer ID
func (a *EventTracker) TrackWorkRequestSerialization(dataSize int, peerId string) {
	event := Event{
		Name:     WorkRequestSerialization,
		PeerID:   peerId,
		WorkType: WorkRequestSerialization,
		Error:    fmt.Sprintf("data size: %d", dataSize), // Assuming data size is stored in Error field for now

	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking work request serialization event: %s", err)
	}
}

// TrackWorkResponseDeserialization records when a work response is deserialized after reception.
//
// Parameters:
// - success: Boolean indicating if the deserialization was successful
// - peerId: String containing the peer ID
func (a *EventTracker) TrackWorkResponseDeserialization(success bool, peerId string) {
	event := Event{
		Name:     WorkResponseDeserialization,
		PeerID:   peerId,
		WorkType: WorkResponseDeserialization,
		Success:  success,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking work response deserialization event: %s", err)
	}
}

// TrackLocalWorkerFallback records when the system falls back to using a local worker.
//
// Parameters:
// - reason: The reason for the fallback
// - peerId: String containing the peer ID
func (a *EventTracker) TrackLocalWorkerFallback(reason string, peerId string) {
	event := Event{
		Name:     LocalWorkerFallback,
		PeerID:   peerId,
		WorkType: LocalWorkerFallback,
		Error:    reason,
	}
	err := a.TrackAndSendEvent(event, nil)
	if err != nil {
		logrus.Errorf("error tracking local worker fallback event: %s", err)
	}
}
