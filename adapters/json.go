package adapters

import (
	"encoding/json"
	"os"

	p "github.com/radugaf/PlentyTelemetry/ports"
)

type JSONDriver struct {
	filename string
}

func NewJSONDriver(filename string) *JSONDriver {
	return &JSONDriver{filename: filename}
}

func (d *JSONDriver) Write(logEntry p.LogEntry) {
	file, err := os.OpenFile(d.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := json.Marshal(logEntry)
	if err != nil {
		return
	}

	file.Write(data)
	file.WriteString("\n")
}
