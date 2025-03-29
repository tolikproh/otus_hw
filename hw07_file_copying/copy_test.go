package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("error file is not exist", func(t *testing.T) {
		pathFrom := "./tmp1"
		pathTo := "./tmp2"
		t.Cleanup(func() {
			_ = os.Remove(pathTo)
		})
		err := Copy(pathFrom, pathTo, 0, 0)
		require.Error(t, err)
		require.True(t, errors.Is(err, os.ErrNotExist), "Ожидается ошибка os.ErrNotExist, получена: %v", err)
	})

	t.Run("file size null", func(t *testing.T) {
		pathFrom := "/dev/urandom"
		pathTo := "./tmp2"

		t.Cleanup(func() {
			_ = os.Remove(pathTo)
		})

		err := Copy(pathFrom, pathTo, 0, 0)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrUnsupportedFile), "Ожидается ошибка ErrUnsupportedFile, получина: %v", err)
	})

	t.Run("offset is not valid", func(t *testing.T) {
		pathFrom := "testdata/input.txt"
		pathTo := "tmp.txt"

		fileFromInfo, err := os.Stat(pathFrom)
		if err != nil {
			t.Fatal(err)
		}

		fileFromSize := fileFromInfo.Size()
		offset := fileFromSize + 1

		t.Cleanup(func() {
			_ = os.Remove(pathTo)
		})

		err = Copy(pathFrom, pathTo, offset, 0)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize), "Ожидается ошибка ErrOffsetExceedsFileSize, получена: %v", err)
	})
}
