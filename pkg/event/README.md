# Masa Protocol Event Tracking Package

A Go package for tracking and sending analytics events.

## Features

- In-memory event storage
- Configurable event sending to external API
- Thread-safe operations
- Comprehensive logging with logrus
- Convenience methods for common event types

## Usage

```go
import "github.com/masa-finance/masa-oracle/pkg/event"

// Create a new event tracker with default config
tracker := event.NewEventTracker(nil)

// Track a custom event
tracker.TrackEvent("custom_event", map[string]interface{}{"key": "value"})

// Use convenience method to track and send a login event
client := event.NewEventClient("https://api.example.com", logger, 10*time.Second)
err := tracker.TrackUserLogin("user123", client)
if err != nil {
    log.Fatal(err)
}

// Retrieve all tracked events
events := tracker.GetEvents()

// Clear all tracked events
tracker.ClearEvents()
```

## Event Library

The package provides a set of predefined events for common scenarios:

### Work Distribution

```go
func (a *EventTracker) TrackWorkDistribution(workType data_types.WorkerType, remoteWorker bool, peerId string, client *EventClient) error
```

Tracks the distribution of work to a worker. Event data includes:
- `peer_id`: String containing the peer ID
- `work_type`: The WorkerType as a string
- `remote_worker`: Boolean indicating if it's a remote worker

### Work Completion

```go
func (a *EventTracker) TrackWorkCompletion(workType data_types.WorkerType, success bool, peerId string, client *EventClient) error
```

Records the completion of a work item. Event data includes:
- `peer_id`: String containing the peer ID
- `work_type`: The WorkerType as a string
- `success`: Boolean indicating if the work was successful

### Worker Failure

```go
func (a *EventTracker) TrackWorkerFailure(workType data_types.WorkerType, errorMessage string, peerId string, client *EventClient) error
```

Records a failure that occurred during work execution. Event data includes:
- `peer_id`: String containing the peer ID
- `work_type`: The WorkerType as a string
- `error`: String containing the error message

## Contributing

Contributions are welcome! Please submit a pull request or create an issue for any bugs or feature requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.