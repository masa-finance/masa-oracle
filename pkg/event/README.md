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

## Contributing

Contributions are welcome! Please submit a pull request or create an issue for any bugs or feature requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.