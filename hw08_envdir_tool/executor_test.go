package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {

	t.Run("invalid command", func(t *testing.T) {
		code := RunCmd([]string{"invalid"}, Environment{})
		require.NotEqual(t, 0, code)
	})

	t.Run("success simple", func(t *testing.T) {
		code := RunCmd([]string{"ls", "-t"}, Environment{})
		require.Equal(t, 0, code)
	})

	t.Run("success simple with env", func(t *testing.T) {
		code := RunCmd([]string{"ls", "-t"}, Environment{"TEST": EnvValue{Value: "My test env"}})
		require.Equal(t, 0, code)
		require.Contains(t, os.Environ(), "TEST=My test env")
	})

}
