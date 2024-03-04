package config

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

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
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
