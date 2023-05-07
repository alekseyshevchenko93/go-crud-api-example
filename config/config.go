package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort string `mapstructure:"HTTP_PORT"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
