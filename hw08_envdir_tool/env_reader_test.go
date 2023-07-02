package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	const dataPath = "./testdata/env"

	expectedMap := Environment{
		"BAR":   {Value: "bar", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: false},
		"FOO": {Value: `   foo
with new line`, NeedRemove: false},
		"HELLO": {Value: "\"hello\"", NeedRemove: false},
		"UNSET": {Value: "", NeedRemove: true},
	}

	t.Run("Case equal compare testdata for ReadDir", func(t *testing.T) {
		envResult, err := ReadDir(dataPath)
		require.NoError(t, err)
		require.Equal(t, expectedMap, envResult)
	})

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

	t.Run("Negative case with uncorrect dir path", func(t *testing.T) {
		_, err := ReadDir("/tmp123_nonexisting")
		require.Error(t, err)
	})
}
