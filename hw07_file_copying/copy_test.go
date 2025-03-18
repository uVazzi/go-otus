package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func copyToTemp(t *testing.T, fromPath string, offset, limit int64) error {
	t.Helper()
	output, err := os.CreateTemp(os.TempDir(), "output")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.Remove(output.Name())
		if err != nil {
			t.Error(err)
		}
	}()

	return Copy(fromPath, output.Name(), offset, limit)
}

func checkCopyFile(t *testing.T, checkPath string, offset, limit int64) {
	t.Helper()
	fromPath := "testdata/input.txt"
	output, err := os.CreateTemp(os.TempDir(), "output")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.Remove(output.Name())
		if err != nil {
			t.Error(err)
		}
	}()

	err = Copy(fromPath, output.Name(), offset, limit)
	require.NoError(t, err)

	outputFile, err := os.ReadFile(output.Name())
	require.NoError(t, err)

	checkedFile, err := os.ReadFile(checkPath)
	require.NoError(t, err)

	require.Equal(t, checkedFile, outputFile)
}

func TestCopy(t *testing.T) {
	t.Run("fromPath and toPath must be specified", func(t *testing.T) {
		err := copyToTemp(t, "", 0, 0)
		require.Truef(t, errors.Is(err, ErrRequiredPath), "actual err - %v", err)

		err = Copy("testdata/input.txt", "", 0, 0)
		require.Truef(t, errors.Is(err, ErrRequiredPath), "actual err - %v", err)
	})

	t.Run("offset and limit must be greater than zero", func(t *testing.T) {
		err := copyToTemp(t, "testdata/input.txt", -100, 0)
		require.Truef(t, errors.Is(err, ErrNotValidOffsetOrLimit), "actual err - %v", err)

		err = copyToTemp(t, "testdata/input.txt", 0, -100)
		require.Truef(t, errors.Is(err, ErrNotValidOffsetOrLimit), "actual err - %v", err)
	})

	t.Run("not exist file", func(t *testing.T) {
		err := copyToTemp(t, "testdataNotFound/input.txt", 0, 0)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := copyToTemp(t, "testdata/input.txt", 7000, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("unsupported file", func(t *testing.T) {
		err := copyToTemp(t, "/dev/urandom", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})

	t.Run("success by testdata", func(t *testing.T) {
		checkCopyFile(t, "testdata/out_offset0_limit0.txt", 0, 0)
		checkCopyFile(t, "testdata/out_offset0_limit10.txt", 0, 10)
		checkCopyFile(t, "testdata/out_offset0_limit1000.txt", 0, 1000)
		checkCopyFile(t, "testdata/out_offset0_limit10000.txt", 0, 10000)
		checkCopyFile(t, "testdata/out_offset100_limit1000.txt", 100, 1000)
		checkCopyFile(t, "testdata/out_offset6000_limit1000.txt", 6000, 1000)
	})
}
