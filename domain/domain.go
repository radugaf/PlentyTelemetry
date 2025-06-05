package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	p "github.com/radugaf/PlentyTelemetry/ports"
)

type Logger struct {
	writers []p.LogWriter
}

func NewLogger(writers ...p.LogWriter) p.LoggingService {
	// fmt.Printf("Creating logger with %d writers\n", len(writers))
	return &Logger{writers: writers}
}

func (l *Logger) Log(level p.LogLevel, msg string, tags map[string]string, txID ...string) {
	entry := p.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Tags:      tags,
	}

	// Handle transaction ID if provided
	if len(txID) > 0 && txID[0] != "" {
		entry.TransactionID = &txID[0]
	}

	// Write to all configured writers
	for _, writer := range l.writers {
		// fmt.Printf("Logging to writer %d\n", i)
		writer.Write(entry)
	}
}

func (l *Logger) Info(msg string, args ...any) {
	tags, txID := parseArgs(args...)
	l.Log(p.Info, msg, tags, txID...)
}

func (l *Logger) Debug(msg string, args ...any) {
	tags, txID := parseArgs(args...)
	l.Log(p.Debug, msg, tags, txID...)
}

func (l *Logger) Warning(msg string, args ...any) {
	tags, txID := parseArgs(args...)
	l.Log(p.Warning, msg, tags, txID...)
}

func (l *Logger) Error(msg string, args ...any) {
	tags, txID := parseArgs(args...)
	l.Log(p.Error, msg, tags, txID...)
}

func parseArgs(args ...any) (map[string]string, []string) {
	var tags map[string]string
	var txID []string

	for i, arg := range args {
		switch v := arg.(type) {
		case map[string]string:
			tags = v
			// fmt.Printf("Found tags at position %d\n", i)
		case string:
			if v != "" {
				txID = []string{v}
				// fmt.Printf("Found txID at position %d: %s\n", i, v)
			}
		default:
			fmt.Printf("Unknown argument type at position %d: %T\n", i, v)
		}
	}

	// Ensure there is always a tags map
	if tags == nil {
		tags = make(map[string]string)
	}

	return tags, txID
}

func (l *Logger) StartTransaction() string {
	txID := uuid.New().String()
	// fmt.Printf("Started transaction: %s\n", txID)
	return txID
}
