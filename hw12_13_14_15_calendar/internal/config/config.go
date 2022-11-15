package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
	HTTP    HTTPConfig  `yaml:"http"`
	GRPC    GRPCConfig  `yaml:"grpc"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageConf struct {
	Type string `yaml:"type"`
	Dsn  string `yaml:"dsn"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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
