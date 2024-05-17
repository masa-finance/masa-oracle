package config

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
)

// SetupLogging configures the logging output for the application.
// It creates the log directory if needed, opens the log file, sets up
// a MultiWriter to log to both stdout and the log file, and configures
// the log level based on the Config.LogLevel field.
func (c *AppConfig) SetupLogging() {
	if _, err := os.Stat(c.MasaDir); os.IsNotExist(err) {
		err = os.MkdirAll(c.MasaDir, 0755)
		if err != nil {
			logrus.Fatal("could not create directory:", err)
		}
	}

	// Open output file for logging
	f, err := os.OpenFile(c.LogFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)

	var logger *logrus.Logger
	if c.LogLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if c.LogLevel == "error" {
		logger = logrus.New()
		logger.Out = os.Stdout
		logger.SetLevel(logrus.ErrorLevel)

		// Setup Google Cloud Logging
		ctx := context.Background()
		client, err := logging.NewClient(ctx, "masa-chain")
		if err != nil {
			logrus.Error(err)
		}
		defer client.Close()

		cloudLogger := client.Logger("cf_chat_log")

		// Configure Logrus to use Google Cloud Logging
		logger.Hooks.Add(&GoogleCloudLoggingHook{
			Client: cloudLogger,
		})

		// Example function that generates an error and logs it
		// gcpLoggerFunction(logger)

	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

}

// func gcpLoggerFunction(logger *logrus.Logger) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			logger.WithField("stack", string(debug.Stack())).Errorf("Panic: %v", r)
// 		}
// 	}()

// 	// Simulate an error
// 	var a *int
// 	*a = 1
// }

// GoogleCloudLoggingHook is a Logrus hook for Google Cloud Logging
type GoogleCloudLoggingHook struct {
	Client *logging.Logger
}

func (hook *GoogleCloudLoggingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *GoogleCloudLoggingHook) Fire(entry *logrus.Entry) error {
	severity := logging.Error
	switch entry.Level {
	case logrus.PanicLevel, logrus.FatalLevel:
		severity = logging.Emergency
	case logrus.ErrorLevel:
		severity = logging.Error
	case logrus.WarnLevel:
		severity = logging.Warning
	case logrus.InfoLevel:
		severity = logging.Info
	case logrus.DebugLevel, logrus.TraceLevel:
		severity = logging.Debug
	}

	payload := map[string]interface{}{
		"message": entry.Message,
		"level":   entry.Level.String(),
		"time":    entry.Time,
	}

	for k, v := range entry.Data {
		payload[k] = v
	}

	hook.Client.Log(logging.Entry{
		Payload:   payload,
		Severity:  severity,
		Timestamp: entry.Time,
	})

	return nil
}
