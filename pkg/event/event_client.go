package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type EventClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     *logrus.Logger
}

func NewEventClient(baseURL string, logger *logrus.Logger, timeout time.Duration) *EventClient {
	return &EventClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: timeout},
		Logger:     logger,
	}
}

func (c *EventClient) SendEvent(event Event) error {
	if c == nil {
		return fmt.Errorf("EventClient is nil")
	}

	url := fmt.Sprintf("%s/%s/events", c.BaseURL, APIVersion)
	payload, err := json.Marshal(event)
	if err != nil {
		c.Logger.WithError(err).Error("Failed to marshal event")
		return err
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		c.Logger.WithError(err).Error("Failed to send event")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("event service returned non-OK status: %d", resp.StatusCode)
		c.Logger.WithError(err).Error("Failed to send event")
		return err
	}

	c.Logger.WithFields(logrus.Fields{
		"event_name": event.Name,
		"timestamp":  event.Timestamp,
	}).Info("Event sent")

	return nil
}
