package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		log.Println("No command")
		return 1
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	// process environments
	ProcessEnv(env)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		log.Println(err)

		var exErr *exec.ExitError

		if errors.As(err, &exErr) {
			return exErr.ExitCode()
		}
		return 1
	}
	return
}
