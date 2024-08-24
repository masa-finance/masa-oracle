package event

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	WorkRequest                 = "work_request"
	WorkCompletion              = "work_completion"
	WorkFailure                 = "worker_failure"
	WorkDistribution            = "work_distribution"
	WorkExecutionStart          = "work_execution_start"
	WorkExecutionTimeout        = "work_execution_timeout"
	RemoteWorkerConnection      = "remote_work_connection"
	StreamCreation              = "stream_creation"
	WorkRequestSerialization    = "work_request_serialized"
	WorkResponseDeserialization = "work_response_serialized"
	LocalWorkerFallback         = "local_work_executed"
)

type Event struct {
	Name         string    `json:"name"`
	Timestamp    time.Time `json:"timestamp"`
	PeerID       string    `json:"peer_id"`
	WorkType     string    `json:"work_type"`
	RemoteWorker bool      `json:"remote_worker"`
	Success      bool      `json:"success"`
	Error        string    `json:"error"`
}

type EventTracker struct {
	events    []Event
	mu        sync.Mutex
	logger    *logrus.Logger
	config    *Config
	apiClient *EventClient
}

func NewEventTracker(config *Config) *EventTracker {
	if config == nil {
		config = DefaultConfig()
	}
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	if level, err := logrus.ParseLevel(config.LogLevel); err == nil {
		logger.SetLevel(level)
	}
	return &EventTracker{
		events:    make([]Event, 0),
		logger:    logger,
		config:    config,
		apiClient: NewEventClient(config.BaseURL, logger, config.HTTPTimeout),
	}
}

func (a *EventTracker) TrackEvent(event Event) {
	if a == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	event.Timestamp = time.Now().UTC()
	a.events = append(a.events, event)
	a.logger.WithFields(logrus.Fields{
		"event_name": event.Name,
		"data":       event,
	}).Info("Event tracked")
}

func (a *EventTracker) GetEvents() []Event {
	if a == nil {
		return nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	return append([]Event{}, a.events...)
}

func (a *EventTracker) ClearEvents() {
	if a == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.events = make([]Event, 0)
	a.logger.Info("Events cleared")
}

func (a *EventTracker) TrackAndSendEvent(event Event, client *EventClient) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	event.Timestamp = time.Now().UTC()
	a.events = append(a.events, event)
	a.logger.WithFields(logrus.Fields{
		"event_name": event.Name,
		"data":       event,
	}).Info("Event tracked")

	if client != nil {
		return client.SendEvent(event)
	} else {
		if a.apiClient != nil {
			err := validateEvent(event)
			if err != nil {
				return err
			}
			return a.apiClient.SendEvent(event)
		}
	}
	return fmt.Errorf("no client available")
}

func validateEvent(event Event) error {
	if event.Name == "" {
		return errors.New("Event name is required")
	}
	if event.Timestamp.IsZero() {
		return errors.New("Invalid timestamp")
	}
	if event.Timestamp.After(time.Now().UTC()) {
		return errors.New("Timestamp cannot be in the future")
	}
	if event.PeerID == "" {
		return errors.New("Peer ID is required")
	}
	if event.WorkType == "" {
		return errors.New("Work type is required")
	}
	return nil
}
