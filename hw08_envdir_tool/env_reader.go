package main

import (
	"bytes"
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

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	fileList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envList := make(Environment, len(fileList))

	for _, file := range fileList {
		name := file.Name()

		if strings.Contains(name, "=") {
			continue
		}

		if env, ok := readFileEnv(dir, name); ok {
			envList[name] = env
		}
	}

	return envList, nil
}

func readFileEnv(dir, name string) (EnvValue, bool) {
	filePath := filepath.Join(dir, name)

	sourceFileInfo, err := os.Stat(filePath)
	if err != nil {
		return EnvValue{}, false
	}
	if !sourceFileInfo.Mode().IsRegular() {
		return EnvValue{}, false
	}

	sourceFileSize := sourceFileInfo.Size()

	if sourceFileSize == 0 {
		return EnvValue{"", true}, true
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return EnvValue{}, false
	}

	str := formatedStr(data)

	return EnvValue{str, false}, true
}

func formatedStr(data []byte) string {
	str := strings.Split(string(data), "\n")[0]
	str = string(bytes.ReplaceAll([]byte(str), []byte{0x00}, []byte("\n")))
	str = strings.TrimRight(str, " \t")
	return str
}
