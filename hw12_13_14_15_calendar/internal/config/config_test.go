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
server:
  addr: ":8080"
`))
	require.NoError(t, err)
	conf, err := NewConfig(file.Name())

	require.NoError(t, err)
	require.Equal(t, "debug", conf.Logger.Level)
	require.Equal(t, "memory", conf.Storage.Type)
	//require.Equal(t, ":8080", conf.Server.Addr)

	_ = os.Remove(file.Name())
}
