package adapters

import (
	"encoding/json"
	"os"

	c "github.com/radugaf/PlentyTelemetry/config"
	p "github.com/radugaf/PlentyTelemetry/ports"
)

type JSONDriver struct {
	filename string
}

func init() {
	c.RegisterDriver("json", func(settings map[string]string) p.LogWriter {
		filename := settings["filename"]
		if filename == "" {
			filename = "logs.json"
		}
		return NewJSONDriver(filename)
	})
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
