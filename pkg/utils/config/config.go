package config

import (
	"fmt"

	testutils "github.com/ripple-mq/benchmark/test"
	"github.com/spf13/viper"
)

type Config struct {
	Broker struct {
		Consumer_buffer_size int
		Producer_buffer_size int
	}
}

var Conf, _ = LoadConfig(".")

func LoadConfig(path string) (*Config, error) {
	testutils.SetRoot()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}
