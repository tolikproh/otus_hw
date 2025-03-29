package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrUnsupportedFile
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
	defer destFile.Close()

	progressBar := pb.Full.Start64(limit)
	sourceFileProgress := progressBar.NewProxyReader(sourceFile)

	_, err = io.CopyN(destFile, sourceFileProgress, limit)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	progressBar.Finish()

	return nil
}
