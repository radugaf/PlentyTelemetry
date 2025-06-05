package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radugaf/PlentyTelemetry/mocks"
	p "github.com/radugaf/PlentyTelemetry/ports"
	"github.com/spf13/viper"
	"go.uber.org/mock/gomock"
)

func TestRegisterDriver(t *testing.T) {
	driverRegistry = make(map[string]DriverFactory)

	testFactory := func(settings map[string]string) p.LogWriter {
		return nil
	}

	RegisterDriver("test", testFactory)

	if len(driverRegistry) != 1 {
		t.Errorf("Expected 1 driver, got %d", len(driverRegistry))
	}

	if _, exists := driverRegistry["test"]; !exists {
		t.Error("Driver 'test' was not registered")
	}
}

func TestCreateDriver_Success(t *testing.T) {
	driverRegistry = make(map[string]DriverFactory)

	ctrl := gomock.NewController(t)

	mockWriter := mocks.NewMockLogWriter(ctrl)

	testFactory := func(settings map[string]string) p.LogWriter {
		return mockWriter
	}
	RegisterDriver("test", testFactory)

	writer := CreateDriver("test", map[string]string{})

	if writer != mockWriter {
		t.Error("Expected to get mock writer back")
	}
}

func TestCreateDriver_WithSettings(t *testing.T) {
	driverRegistry = make(map[string]DriverFactory)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWriter := mocks.NewMockLogWriter(ctrl)

	testFactory := func(settings map[string]string) p.LogWriter {
		if settings["filename"] == "test.log" {
			return mockWriter
		}
		return nil
	}
	RegisterDriver("file", testFactory)

	settings := map[string]string{"filename": "test.log"}
	writer := CreateDriver("file", settings)

	if writer != mockWriter {
		t.Error("Expected to get mock writer with correct settings")
	}
}

func TestLoadConfig_WithValidFile(t *testing.T) {
	viper.Reset()

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	configContent := `drivers:
  - type: cli
    enabled: true
    settings: {}
  - type: json
    enabled: false
    settings:
      filename: test.json`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tempDir)

	config, err := LoadConfig()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(config.Drivers) != 2 {
		t.Errorf("Expected 2 drivers, got %d", len(config.Drivers))
	}

	// Check first driver
	if config.Drivers[0].Type != "cli" {
		t.Errorf("Expected first driver type 'cli', got '%s'", config.Drivers[0].Type)
	}
	if !config.Drivers[0].Enabled {
		t.Error("Expected first driver to be enabled")
	}

	// Check second driver
	if config.Drivers[1].Type != "json" {
		t.Errorf("Expected second driver type 'json', got '%s'", config.Drivers[1].Type)
	}
	if config.Drivers[1].Enabled {
		t.Error("Expected second driver to be disabled")
	}
	if config.Drivers[1].Settings["filename"] != "test.json" {
		t.Errorf("Expected filename 'test.json', got '%s'", config.Drivers[1].Settings["filename"])
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	viper.Reset()

	// Create temp directory and invalid config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")

	invalidYAML := `drivers:
  - type: cli
    enabled: true
    settings: {
  - invalid yaml here`

	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tempDir)

	_, err = LoadConfig()

	// Should return error
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestDriverConfig_Structure(t *testing.T) {
	config := DriverConfig{
		Type:    "test",
		Enabled: true,
		Settings: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	if config.Type != "test" {
		t.Errorf("Expected type 'test', got '%s'", config.Type)
	}
	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
	if len(config.Settings) != 2 {
		t.Errorf("Expected 2 settings, got %d", len(config.Settings))
	}
	if config.Settings["key1"] != "value1" {
		t.Errorf("Expected setting key1='value1', got '%s'", config.Settings["key1"])
	}
}
