package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrRequiredPath          = errors.New("fromPath and toPath must be specified")
	ErrNotValidOffsetOrLimit = errors.New("offset and limit must be greater than zero")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrRequiredPath
	}
	if offset < 0 || limit < 0 {
		return ErrNotValidOffsetOrLimit
	}

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if limit == 0 {
		limit = fileInfo.Size()
	}
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	inputFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	_, err = inputFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	copySize := fileInfo.Size() - offset
	if limit < copySize {
		copySize = limit
	}

	bar := pb.Full.Start64(copySize)
	defer bar.Finish()

	reader := io.LimitReader(inputFile, copySize)
	barReader := bar.NewProxyReader(reader)
	_, err = io.CopyN(outputFile, barReader, copySize)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
