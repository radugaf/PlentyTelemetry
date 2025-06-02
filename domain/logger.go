package domain

import (
	"time"

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

func (l *Logger) Log(level p.LogLevel, msg string, tags map[string]string) {
	entry := p.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Tags:      tags,
	}

	for _, writer := range l.writers {
		writer.Write(entry)
	}
}

func (l *Logger) Info(msg string, tags map[string]string) {
	l.Log(p.Info, msg, tags)
}

func (l *Logger) Debug(msg string, tags map[string]string) {
	l.Log(p.Debug, msg, tags)
}

func (l *Logger) Warning(msg string, tags map[string]string) {
	l.Log(p.Warning, msg, tags)
}

func (l *Logger) Error(msg string, tags map[string]string) {
	l.Log(p.Error, msg, tags)
}
