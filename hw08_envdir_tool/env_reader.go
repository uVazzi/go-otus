package main

import (
	"bufio"
	"errors"
	"io"
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

var (
	ErrReadEnvDir = errors.New("can't read env directory")
	ErrFileInfo   = errors.New("can't get file info")
	ErrFileName   = errors.New("filename must not contain an equal sign")
	ErrOpenFile   = errors.New("can't open file")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrReadEnvDir
	}

	env := make(Environment)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		fileInfo, err := entry.Info()
		if err != nil {
			return nil, ErrFileInfo
		}

		if strings.Contains(filename, "=") {
			return nil, ErrFileName
		}
		if fileInfo.Size() == 0 {
			env[filename] = EnvValue{"", true}
			continue
		}

		file, err := os.Open(filepath.Join(dir, filename))
		if err != nil {
			return nil, ErrOpenFile
		}
		reader := bufio.NewReader(file)
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, ErrOpenFile
		}
		file.Close()

		line = strings.TrimRight(line, " \t\n")
		line = strings.ReplaceAll(line, "\x00", "\n")

		env[filename] = EnvValue{line, false}
	}

	return env, nil
}
