package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const dataPath = "./testdata/env"

var expectedMap = Environment{
	"BAR":   {Value: "bar", NeedRemove: false},
	"EMPTY": {Value: "", NeedRemove: false},
	"FOO": {Value: `   foo
with new line`, NeedRemove: false},
	"HELLO": {Value: "\"hello\"", NeedRemove: false},
	"UNSET": {Value: "", NeedRemove: true},
}

func TestReadDir(t *testing.T) {

	t.Run("Case equal compare testdata for ReadDir", func(t *testing.T) {
		envResult, err := ReadDir(dataPath)
		require.NoError(t, err)
		require.Equal(t, expectedMap, envResult)
	})

	t.Run("Negative case with uncorrect dir path", func(t *testing.T) {
		_, err := ReadDir("/tmp123_nonexisting")
		require.Error(t, err)
	})
}

func TestProcessEnv(t *testing.T) {
	t.Run("Case set and unset environments", func(t *testing.T) {
		ProcessEnv(expectedMap)

		for envName, envValue := range expectedMap {
			if !envValue.NeedRemove {
				require.Equal(t, envValue.Value, os.Getenv(envName))
			} else {
				_, ok := os.LookupEnv(envName)
				require.False(t, ok)
			}
		}
	})

	t.Run("Case overwrite environments", func(t *testing.T) {
		envKey := "MYBAR"
		envMyVal := "mybarvalue123"
		expected := Environment{envKey: {Value: "12345", NeedRemove: false}}

		os.Setenv(envKey, envMyVal)
		require.Equal(t, envMyVal, os.Getenv(envKey))

		ProcessEnv(expected)

		require.NotEqual(t, envMyVal, os.Getenv(envKey))
		require.Equal(t, expected[envKey].Value, os.Getenv(envKey))
	})

	t.Run("Case with uncorrect environment", func(t *testing.T) {
		envKey := "MYBAR="
		expected := Environment{envKey: {Value: "12345", NeedRemove: false}}

		ProcessEnv(expected)

		_, ok := os.LookupEnv(envKey)
		require.False(t, ok)
	})
}

func TestReadEnv(t *testing.T) {
	t.Run("Case with environment files", func(t *testing.T) {
		for envName, envValue := range expectedMap {
			resultEnv, err := ReadEnv(filepath.Join(dataPath, envName))
			require.NoError(t, err)
			require.Equal(t, envValue, resultEnv)
		}
	})

	t.Run("Negative case with environment files", func(t *testing.T) {
		resultEnv, err := ReadEnv("testdata/env/BAR=")
		require.Error(t, err)
		require.Equal(t, EnvValue{"", true}, resultEnv)
	})
}
