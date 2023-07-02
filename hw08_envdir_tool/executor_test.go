package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Negative case with no command", func(t *testing.T) {
		const expectedRc = 1

		cmd := []string{}
		rc := RunCmd(cmd, expectedMap)
		require.Equal(t, expectedRc, rc)
	})

	t.Run("Case shell exit code", func(t *testing.T) {
		expected := []struct {
			rc      int
			command string
		}{
			{0, "exit 0"},
			{1, "exit 1"},
			{126, "exit 126"},
			{127, "exit 127"},
			{128, "exit 128"},
			{130, "exit 130"},
		}

		for _, entry := range expected {
			rc := RunCmd([]string{"sh", "-c", entry.command}, expectedMap)
			require.Equal(t, entry.rc, rc)
		}
	})

	t.Run("Case with comparing local env", func(t *testing.T) {
		command := []string{"sh", "-c", "[ $BAR = \"bar\" ] || exit 1"}
		rc := RunCmd(command, expectedMap)
		require.Equal(t, 0, rc)

		command = []string{"sh", "-c", "[ $BAR = \"barfoo\" ] || exit 1"}
		rc = RunCmd(command, expectedMap)
		require.Equal(t, 1, rc)
	})
}
