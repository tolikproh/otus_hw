package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testCases := []struct {
		name     string
		cmd      []string
		expected int
	}{
		{
			name:     "exit code = 0 (success)",
			cmd:      []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"},
			expected: 0,
		},
		{
			name:     "exit code = 127 (command not found)",
			cmd:      []string{"/bin/bash", "./notfile", "arg1=1", "arg2=2"},
			expected: 127,
		},
		{
			name:     "exit code = 127 (general error)",
			cmd:      []string{"test"},
			expected: 1,
		},
		{
			name:     "exit code = 0 (success)",
			cmd:      []string{"test", "-test"},
			expected: 0,
		},
		{
			name:     "exit code = -1 (executable file not found)",
			cmd:      []string{"-test-"},
			expected: -1,
		},
	}

	env, err := ReadDir("./testdata/env")
	require.Nil(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			retCode := RunCmd(tc.cmd, env)
			require.Equal(t, tc.expected, retCode)
		})
	}
}
