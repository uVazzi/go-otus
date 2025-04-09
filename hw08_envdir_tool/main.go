package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var ErrArgs = errors.New("not enough args")

func main() {
	args, err := getArgs()
	if err != nil {
		os.Exit(GetErrorCodeAndPrintError(err))
	}

	env, err := ReadDir(args[1])
	if err != nil {
		os.Exit(GetErrorCodeAndPrintError(err))
	}

	code := RunCmd(args[2:], env)
	os.Exit(code)
}

func GetErrorCodeAndPrintError(err error) int {
	fmt.Println("Error:", err)
	var exitErr *exec.ExitError

	switch {
	case errors.As(err, &exitErr):
		return exitErr.ExitCode()
	case errors.Is(err, ErrArgs):
		return 1
	case errors.Is(err, ErrReadEnvDir):
		return 2
	case errors.Is(err, ErrFileInfo):
		return 3
	case errors.Is(err, ErrFileName):
		return 4
	case errors.Is(err, ErrOpenFile):
		return 5
	default:
		return 125 // unexpected error
	}
}

func getArgs() ([]string, error) {
	args := os.Args
	if len(args) < 3 {
		return nil, ErrArgs
	}

	return args, nil
}
