package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	// TODO
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type string
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
