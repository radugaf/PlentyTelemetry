package config

import (
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

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Defaults
	viper.SetDefault("drivers", []DriverConfig{
		{Type: "cli", Enabled: true, Settings: map[string]string{}},
	})

	viper.ReadInConfig()

	var config Config
	return &config, viper.Unmarshal(&config)
}
