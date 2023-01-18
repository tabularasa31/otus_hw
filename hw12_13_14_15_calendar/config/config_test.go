package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	file, _ := os.CreateTemp("/tmp", "config.*.yml")
	_, err := file.Write([]byte(
		`logger:
  level: debug
storage:
  type: memory
  dsn: ''
http:
  addr: 'localhost:8080'
grpc:
  addr: 'localhost:8081'
`))
	require.NoError(t, err)
	conf, err := NewConfig(file.Name())

	require.NoError(t, err)
	require.Equal(t, "debug", conf.Logger.Level)
	require.Equal(t, "memory", conf.Storage.Type)
	require.Equal(t, "localhost:8080", conf.HTTP.Addr)
	require.Equal(t, "localhost:8081", conf.GRPC.Addr)

	_ = os.Remove(file.Name())
}
