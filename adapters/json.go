package adapters

import (
	"encoding/json"
	"fmt"
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
			filename = "logs.jsonl" // default filename
			fmt.Println("Using default filename: logs.json")
		}
		return NewJSONDriver(filename)
	})
}

func NewJSONDriver(filename string) *JSONDriver {
	fmt.Printf("Creating JSON driver with file: %s\n", filename)
	return &JSONDriver{filename: filename}
}

func (d *JSONDriver) Write(logEntry p.LogEntry) {
	file, err := os.OpenFile(d.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to open file %s: %v\n", d.filename, err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Printf("Failed to marshal log entry: %v\n", err)
		return
	}

	// Write the JSON + newline
	file.Write(data)
	file.WriteString("\n")
}
