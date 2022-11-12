package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type string
	Dsn  string
}

func NewConfig(configFile string) (Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %v", err)
	}
	cfg := Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config file: %v", err)
	}
	return cfg, nil
}
