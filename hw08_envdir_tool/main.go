package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	envDirectory := args[0]
	commandArgs := args[1:]

	if len(args) < 2 {
		log.Fatalf("there are no any command arguments")
	}

	env, err := ReadDir(envDirectory)
	if err != nil {
		log.Fatalf("can't find directory \"%v\", error : %v", envDirectory, err)
	}

	os.Exit(RunCmd(commandArgs, env))
}
