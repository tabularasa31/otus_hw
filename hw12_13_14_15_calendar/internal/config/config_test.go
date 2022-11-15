package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	file, _ := os.CreateTemp("/tmp", "config.*.yaml")
	_, err := file.Write([]byte(
		`logger:
  level: debug
storage:
  type: memory
  dsn: ''
http:
  host: 'localhost'
  port: 8080
grpc:
  host: ''
  port: 8081
`))
	require.NoError(t, err)
	conf, err := NewConfig(file.Name())

	require.NoError(t, err)
	require.Equal(t, "debug", conf.Logger.Level)
	require.Equal(t, "memory", conf.Storage.Type)
	require.Equal(t, "localhost", conf.HTTP.Host)
	require.Equal(t, "8080", conf.HTTP.Port)
	require.Equal(t, "", conf.GRPC.Host)
	require.Equal(t, "8081", conf.GRPC.Port)

	_ = os.Remove(file.Name())
}
