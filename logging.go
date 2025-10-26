package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func startLogging() {

	logLevel := parseLogLevel()
	logrus.SetLevel(logLevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Error creating log directory: %s\n", err)
		os.Exit(1)
	}

	logFileName := filepath.Join(logDir, fmt.Sprintf("application_%s.log", time.Now().Format("2006-01-02")))

	// Create or open log file with append mode
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err)
		os.Exit(1)
	}

	// Set logrus to write to the file
	logrus.SetOutput(logFile)

}

// parseLogLevel tries to parse the environment variable as a direct log level name
func parseLogLevel() logrus.Level {

	ll := os.Getenv("LogLevel")

	// Convert to lowercase for case-insensitive comparison
	ll = strings.ToLower(ll)

	switch strings.ToLower(ll) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		// Default to Info level if the string doesn't match any known level
		return logrus.InfoLevel
	}
}
