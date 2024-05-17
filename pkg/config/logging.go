package config

import (
	"io"
	"log"
	"os"

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

	if c.LogLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if c.LogLevel == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

}
