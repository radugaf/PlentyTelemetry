package ports

//go:generate mockgen -source=ports.go -destination=../mocks/mock_ports.go -package=mocks

import "time"

type LogLevel string

const (
	Debug   LogLevel = "DEBUG"
	Info    LogLevel = "INFO"
	Warning LogLevel = "WARNING"
	Error   LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp     time.Time
	Level         LogLevel
	Message       string
	Tags          map[string]string
	TransactionID *string
}

type LoggingService interface {
	Log(level LogLevel, msg string, tags map[string]string, txID ...string)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warning(msg string, args ...any)
	Error(msg string, args ...any)
	StartTransaction() string
}

type LogWriter interface {
	Write(entry LogEntry)
}
