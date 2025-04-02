package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("can't read env directory", func(t *testing.T) {
		_, err := ReadDir("./notExistDir/env")
		require.Truef(t, errors.Is(err, ErrReadEnvDir), "actual err - %v", err)
	})

	t.Run("filename must not contain an equal sign", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile, _ := os.CreateTemp(tempDir, "TEST=TEST")
		tempFile.Close()

		_, err := ReadDir(tempDir)
		require.Truef(t, errors.Is(err, ErrFileName), "actual err - %v", err)
	})

	t.Run("success by testdata", func(t *testing.T) {
		outputEnv, err := ReadDir("./testdata/env")
		checkedEnv := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}

		require.NoError(t, err)
		require.Equal(t, checkedEnv, outputEnv)
	})
}
