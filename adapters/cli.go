package adapters

import (
	"fmt"
	"strings"

	p "github.com/radugaf/PlentyTelemetry/ports"
)

type CLIDriver struct{}

func NewCLIDriver() *CLIDriver {
	return &CLIDriver{}
}

func (d *CLIDriver) Write(logEntry p.LogEntry) {
	var output strings.Builder

	fmt.Fprintf(&output, "[%s] %s: %s",
		logEntry.Timestamp.Format("2006-01-02 15:04:05.000"),
		logEntry.Level,
		logEntry.Message,
	)

	if len(logEntry.Tags) > 0 {
		fmt.Fprintf(&output, " | Tags: %v", logEntry.Tags)
	}

	if logEntry.TransactionID != nil {
		fmt.Fprintf(&output, " | TxID: %s", *logEntry.TransactionID)
	}

	fmt.Println(output.String())
}
