package adapters

import (
	"fmt"

	p "github.com/radugaf/PlentyTelemetry/ports"
)

type CLIDriver struct{}

func NewCLIDriver() *CLIDriver {
	return &CLIDriver{}
}

func (d *CLIDriver) Write(logEntry p.LogEntry) {
	fmt.Printf("[%s] %s: %s | Tags: %v\n",
		logEntry.Timestamp.Format("2006-01-02 15:04:05"),
		logEntry.Level,
		logEntry.Message,
		logEntry.Tags,
	)
}
