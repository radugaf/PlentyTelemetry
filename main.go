package main

import (
	a "github.com/radugaf/PlentyTelemetry/adapters"
	logger "github.com/radugaf/PlentyTelemetry/domain"
)

func main() {
	cliAdapter := a.NewCLIDriver()
	logger := logger.NewLogger(cliAdapter)
	logger.Info("Test Started", map[string]string{
		"origin":     "http",
		"customerId": "123",
	})
	logger.Debug("Debug", nil)
	logger.Error("Error", nil)
}
