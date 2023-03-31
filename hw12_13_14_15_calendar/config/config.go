package config

import (
	"fmt"

	val "github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Logger   `yaml:"logger"`
		Storage  `yaml:"storage" validate:"required"`
		HTTP     `yaml:"http" validate:"required"`
		GRPC     `yaml:"grpc" validate:"required"`
		Postgres `yaml:"postgres"`
	}

	Logger struct {
		Level string `yaml:"level"`
	}

	Storage struct {
		Type string `yaml:"type" validate:"required"`
	}

	HTTP struct {
		Addr string `yaml:"addr" validate:"required,hostname_port" env:"HTTP_ADDR"`
	}

	GRPC struct {
		Addr string `yaml:"addr" validate:"required,hostname_port" env:"GRPC_ADDR"`
	}

	Postgres struct {
		Dsn     string `yaml:"dsn"`
		PoolMax int    `yaml:"poolMax" env:"PG_POOL_MAX"`
	}

	SchedulerConf struct {
		Scheduler  `yaml:"scheduler"`
		AMQPConfig `yaml:"rabbitmq"`
	}

	Scheduler struct {
		DeletePeriod int `yaml:"deletePeriod"`
	}

	AMQPConfig struct {
		Addr         string `yaml:"addr"`
		Exchange     string `yaml:"exchange"`
		ExchangeType string `yaml:"exchangeType"`
		Queue        string `yaml:"queue"`
		ConsumerTag  string `yaml:"consumerTag"`
		BindingKey   string `yaml:"bindingKey"`
		Reliable     string `yaml:"reliable"`
		Persistent   string `yaml:"persistent"`
	}

	SenderConf struct {
		Sender     `yaml:"sender"`
		AMQPConfig `yaml:"rabbitmq"`
	}

	Sender struct{}
)

// NewConfig returns app config.
func NewConfig(file string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(file, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	v := val.New()
	if err = v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("failed config: %w", err)
	}

	return cfg, nil
}

// NewSchedulerConfig returns scheduler config.
func NewSchedulerConfig(file string) (*SchedulerConf, error) {
	cfg := &SchedulerConf{}

	err := cleanenv.ReadConfig(file, cfg)
	if err != nil {
		return nil, fmt.Errorf("ampq config error: %w", err)
	}

	return cfg, nil
}

// NewSenderConfig returns sender config.
func NewSenderConfig(file string) (*SenderConf, error) {
	cfg := &SenderConf{}

	err := cleanenv.ReadConfig(file, cfg)
	if err != nil {
		return nil, fmt.Errorf("ampq config error: %w", err)
	}

	return cfg, nil
}
