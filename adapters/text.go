package adapters

import (
	"fmt"
	"os"
	"strings"

	c "github.com/radugaf/PlentyTelemetry/config"
	p "github.com/radugaf/PlentyTelemetry/ports"
)

type TextDriver struct {
	filename string
}

func init() {
	c.RegisterDriver("text", func(settings map[string]string) p.LogWriter {
		return NewTextDriver(settings["filename"])
	})
}

func NewTextDriver(filename string) *TextDriver {
	return &TextDriver{filename: filename}
}

func (d *TextDriver) Write(logEntry p.LogEntry) {
	file, err := os.OpenFile(d.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

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

	fmt.Fprintln(file, output.String())
}
