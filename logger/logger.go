package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps the zerolog.Logger to provide a structured logging interface.
type Logger struct {
	Log zerolog.Logger
}

// NewLogger initializes and returns a new Logger instance with console output.
func NewLogger() *Logger {
	// Create a new zerolog logger with console writer for output.
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,    // Output logs to standard output.
		TimeFormat: time.RFC3339, // Use RFC3339 format for timestamps.
	}).
		With().
		Timestamp(). // Add timestamp to log entries.
		Logger()

	return &Logger{Log: logger}
}
