package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Logger  `yaml:"logger"`
		Storage `yaml:"storage" validate:"required"`
		HTTP    `yaml:"http" validate:"required"`
		GRPC    `yaml:"grpc" validate:"required"`
	}

	Logger struct {
		Level string `yaml:"level"`
	}

	Storage struct {
		Type string `yaml:"type" validate:"required"`
		Dsn  string `yaml:"dsn"`
	}

	HTTP struct {
		Addr string `yaml:"addr" validate:"required,hostname_port" env:"HTTP_ADDR"`
	}

	GRPC struct {
		Addr string `yaml:"addr" validate:"required,hostname_port" env:"GRPC_ADDR"`
	}
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

	v := validator.New()
	if err = v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("failed config: %w", err)
	}

	return cfg, nil
}
