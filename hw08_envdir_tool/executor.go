package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
//
//nolint:gosec
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return -1
	}

	cmds := exec.Command(cmd[0], cmd[1:]...)
	cmds.Stdin = os.Stdin
	cmds.Stdout = os.Stdout
	cmds.Stderr = os.Stderr
	cmds.Env = updateEnv(env)
	cmds.Run()

	return cmds.ProcessState.ExitCode()
}

func updateEnv(env Environment) []string {
	for key, env := range env {
		if env.NeedRemove {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, env.Value)
		}
	}

	return os.Environ()
}
