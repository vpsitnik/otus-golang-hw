package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(map[string]EnvValue)
	entriesDir, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, entry := range entriesDir {
		entryName := entry.Name()

		if !entry.Type().IsRegular() || strings.Contains(entryName, "=") {
			continue
		}

		envVal, err := ReadEnv(filepath.Join(dir, entryName))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		envs[entryName] = envVal
	}

	return envs, nil
}

func ReadEnv(filePath string) (EnvValue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return EnvValue{"", true}, err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		log.Println(err)
		return EnvValue{"", true}, err
	}
	if fileStat.Size() == 0 {
		return EnvValue{"", true}, nil
	}

	scanner := bufio.NewScanner(file)
	envVal := ""
	for scanner.Scan() {
		envVal = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
		return EnvValue{"", true}, err
	}

	envVal = strings.ReplaceAll(strings.TrimRight(envVal, " \t"), "0x00", "\n")
	return EnvValue{envVal, false}, nil
}

func ProcessEnv(envs Environment) {
	for envName, envValue := range envs {
		if _, ok := os.LookupEnv(envName); ok || envValue.NeedRemove {
			os.Unsetenv(envName)
		}

		if !envValue.NeedRemove {
			os.Setenv(envName, envValue.Value)
		}
	}
}
