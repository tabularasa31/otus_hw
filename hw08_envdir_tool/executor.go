package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		err := os.Unsetenv(k)
		if err != nil {
			log.Println(err)
			return 1
		}
		if !v.NeedRemove {
			err := os.Setenv(k, v.Value)
			if err != nil {
				log.Println(err)
				return 1
			}
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		return 1
	}
	return 0
}
