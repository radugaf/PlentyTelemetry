package config

import (
	"fmt"

	p "github.com/radugaf/PlentyTelemetry/ports"
	"github.com/spf13/viper"
)

type DriverConfig struct {
	Type     string            `mapstructure:"type"`
	Enabled  bool              `mapstructure:"enabled"`
	Settings map[string]string `mapstructure:"settings"`
}

type Config struct {
	Drivers []DriverConfig `mapstructure:"drivers"`
}

// Registry implementation
type DriverFactory func(settings map[string]string) p.LogWriter

var driverRegistry = make(map[string]DriverFactory)

func RegisterDriver(name string, factory DriverFactory) {
	driverRegistry[name] = factory
}

func CreateDriver(driverType string, settings map[string]string) p.LogWriter {
	if factory, exists := driverRegistry[driverType]; exists {
		return factory(settings)
	}
	return nil
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("drivers", []DriverConfig{
		{Type: "cli", Enabled: true, Settings: map[string]string{}},
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

	return &config, nil
}
