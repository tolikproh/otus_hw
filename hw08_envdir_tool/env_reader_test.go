package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testData = "./testdata/env"

func TestReadDir(t *testing.T) {
	expCase := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": EnvValue{"", true},
	}

	t.Run("проверка на тестовых данных", func(t *testing.T) {
		envList, err := ReadDir(testData)
		require.Nil(t, err)
		require.Equal(t, expCase, envList)
	})

	t.Run("директории не существует", func(t *testing.T) {
		_, err := ReadDir("tmp_tmp_tmp")
		require.NotEqual(t, nil, err)
	})

	t.Run("недопустимый тип директории", func(t *testing.T) {
		env, err := ReadDir("/dev")
		require.Equal(t, Environment{}, env)
		require.Equal(t, nil, err)
	})

	t.Run("игнорирование файла с '='", func(t *testing.T) {
		envTest, err := os.CreateTemp(testData, "*=")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(envTest.Name())
		envList, err := ReadDir(testData)
		require.Nil(t, err)
		require.Equal(t, expCase, envList)
	})
}
