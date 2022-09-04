package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("invalid dir path", func(t *testing.T) {
		_, err := ReadDir("invalid/dir/path")
		require.Error(t, err)
	})
	t.Run("read testdata dir", func(t *testing.T) {
		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}
		dir := "./testdata/env"
		env, err := ReadDir(dir)
		assert.NoError(t, err)
		assert.Equal(t, expected, env)
	})
}
