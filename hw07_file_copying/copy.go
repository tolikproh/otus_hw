package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrIncorrectFileName     = errors.New("incorrect file name")
	ErrEmptyFilePath         = errors.New("empty file path from or to")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrEmptyFilePath
	}

	var err error

	fromPath, err = filepath.Abs(fromPath)
	if err != nil {
		return err
	}

	toPath, err = filepath.Abs(toPath)
	if err != nil {
		return err
	}

	if fromPath == toPath {
		return ErrIncorrectFileName
	}

	sourceFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	sourceFileSize := sourceFileInfo.Size()
	if sourceFileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > sourceFileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > sourceFileSize-offset {
		limit = sourceFileSize - offset
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if offset > 0 {
		if _, err := sourceFile.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	progressBar := pb.Full.Start64(limit)

	defer func() {
		progressBar.Finish()
		destFile.Close()
		if err != nil {
			os.Remove(toPath)
		}
	}()

	sourceFileProgress := progressBar.NewProxyReader(sourceFile)
	if _, err = io.CopyN(destFile, sourceFileProgress, limit); err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
			return nil
		}
		return err
	}

	return nil
}
