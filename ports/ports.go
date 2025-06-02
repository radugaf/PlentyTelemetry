package ports

import "time"

type LogLevel string

const (
	Debug   LogLevel = "DEBUG"
	Info    LogLevel = "INFO"
	Warning LogLevel = "WARNING"
	Error   LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Tags      map[string]string
}

type LoggingService interface {
	Log(level LogLevel, msg string, tags map[string]string)
	Info(msg string, tags map[string]string)
	Debug(msg string, tags map[string]string)
	Warning(msg string, tags map[string]string)
	Error(msg string, tags map[string]string)
}

type LogWriter interface {
	Write(entry LogEntry)
}
