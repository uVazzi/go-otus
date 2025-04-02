package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("not found command", func(t *testing.T) {
		cmd := []string{"/notExistCommand", "./testdata/env", "arg1=1", "arg2=2"}
		exitCode := RunCmd(cmd, Environment{})
		require.Equal(t, 125, exitCode)
	})

	t.Run("success by testdata", func(t *testing.T) {
		os.Setenv("HELLO", "SHOULD_REPLACE")
		os.Setenv("FOO", "SHOULD_REPLACE")
		os.Setenv("UNSET", "SHOULD_REMOVE")
		os.Setenv("ADDED", "from original env")
		os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

		cmd := []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}
		env := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}

		originalStdout := os.Stdout
		var buf bytes.Buffer
		readStdout, writeStdout, _ := os.Pipe()
		os.Stdout = writeStdout

		exitCode := RunCmd(cmd, env)

		writeStdout.Close()
		os.Stdout = originalStdout
		io.Copy(&buf, readStdout)

		output := buf.String()
		output = strings.TrimRight(output, " \t\n")
		outputSplit := strings.Split(output, "\n")

		expected := []string{
			"HELLO is (\"hello\")",
			"BAR is (bar)",
			"FOO is (   foo",
			"with new line)",
			"UNSET is ()",
			"ADDED is (from original env)",
			"EMPTY is ()",
			"arguments are arg1=1 arg2=2",
		}

		require.Equal(t, 0, exitCode)
		require.Equal(t, expected, outputSplit)
	})
}
