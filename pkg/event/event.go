package event

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Event struct {
	Name      string
	Timestamp time.Time `json:"timestamp"`
	Data      map[string]interface{}
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

func (a *EventTracker) TrackEvent(name string, data map[string]interface{}) {
	if a == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	event := Event{
		Name:      name,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	a.events = append(a.events, event)
	a.logger.WithFields(logrus.Fields{
		"event_name": name,
		"data":       data,
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

func (a *EventTracker) TrackAndSendEvent(name string, data map[string]interface{}, client *EventClient) error {
	if a == nil {
		return fmt.Errorf("analytics is nil")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	event := Event{
		Name:      name,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	a.events = append(a.events, event)
	a.logger.WithFields(logrus.Fields{
		"event_name": name,
		"data":       data,
	}).Info("Event tracked")

	if client != nil {
		return client.SendEvent(event)
	}

	return nil
}
