package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Error, not enough args. Type correct path and command with args")
	}

	envs, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	rc := RunCmd(os.Args[2:], envs)
	os.Exit(rc)
}
