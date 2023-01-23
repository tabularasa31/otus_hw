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
		PoolMax int    `yaml:"pool_max" env:"PG_POOL_MAX"`
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

	v := val.New()
	if err = v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("failed config: %w", err)
	}

	return cfg, nil
}
