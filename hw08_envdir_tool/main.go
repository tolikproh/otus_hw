package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Необходимо использовать два аргумента:")
		fmt.Println("   - первый агумент: директория с переменным окружением")
		fmt.Println("   - второй аргумет: программа для запуска ")
		return
	}

	dir := args[1]
	command := args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(RunCmd(command, env))
}
