package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for _, environItem := range os.Environ() {
		envItemKey, envItemValue, _ := strings.Cut(environItem, "=")
		if _, ok := env[envItemKey]; !ok {
			env[envItemKey] = EnvValue{envItemValue, false}
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	for key, envValue := range env {
		if !envValue.NeedRemove {
			command.Env = append(command.Env, key+"="+envValue.Value)
		}
	}

	err := command.Run()
	if err != nil {
		return GetErrorCodeAndPrintError(err)
	}
	return command.ProcessState.ExitCode()
}
