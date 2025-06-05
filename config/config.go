package config

import (
	"fmt"

	p "github.com/radugaf/PlentyTelemetry/ports"
	"github.com/spf13/viper"
)

type Config struct {
	Drivers []DriverConfig `mapstructure:"drivers"`
}

type DriverConfig struct {
	Type     string            `mapstructure:"type"`
	Enabled  bool              `mapstructure:"enabled"`
	Settings map[string]string `mapstructure:"settings"`
}

// Registry implementation
type DriverFactory func(settings map[string]string) p.LogWriter

var driverRegistry = make(map[string]DriverFactory)

func RegisterDriver(name string, factory DriverFactory) {
	fmt.Printf("Registering driver: %s\n", name)
	driverRegistry[name] = factory
}

func CreateDriver(driverType string, settings map[string]string) p.LogWriter {
	fmt.Printf("Creating driver of type: %s\n", driverType)

	if factory, exists := driverRegistry[driverType]; exists {
		writer := factory(settings)
		if writer != nil {
			fmt.Printf("Successfully created %s driver\n", driverType)
		} else {
			fmt.Printf("Failed to create %s driver\n", driverType)
		}
		return writer
	}

	fmt.Printf("Unknown driver type: %s\n", driverType)
	return nil
}

func LoadConfig() (*Config, error) {
	fmt.Println("Loading configuration...")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Defaults in case config file is missing
	viper.SetDefault("drivers", []DriverConfig{
		{
			Type:    "cli",
			Enabled: true,
		},
	})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	fmt.Printf("Loaded %d driver configurations\n", len(config.Drivers))
	return &config, nil
}
