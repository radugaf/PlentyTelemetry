package domain

import (
	"time"

	"github.com/google/uuid"
	p "github.com/radugaf/PlentyTelemetry/ports"
)

type Logger struct {
	writers []p.LogWriter
}

func NewLogger(writers ...p.LogWriter) p.LoggingService {
	return &Logger{
		writers: writers,
	}
}

func (l *Logger) Log(level p.LogLevel, msg string, tags map[string]string, txID ...string) {
	entry := p.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Tags:      tags,
	}

	if len(txID) > 0 && txID[0] != "" {
		entry.TransactionID = &txID[0]
	}

	for _, writer := range l.writers {
		writer.Write(entry)
	}
}

func (l *Logger) Info(msg string, tags map[string]string, txID ...string) {
	l.Log(p.Info, msg, tags, txID...)
}

func (l *Logger) Debug(msg string, tags map[string]string, txID ...string) {
	l.Log(p.Debug, msg, tags, txID...)
}

func (l *Logger) Warning(msg string, tags map[string]string, txID ...string) {
	l.Log(p.Warning, msg, tags, txID...)
}

func (l *Logger) Error(msg string, tags map[string]string, txID ...string) {
	l.Log(p.Error, msg, tags, txID...)
}

func (l *Logger) StartTransaction() string {
	return uuid.New().String()
}
