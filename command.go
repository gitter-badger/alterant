package main

import (
	"log"
	"os"
	"os/exec"
)

func (t *task) executeCommands(cfg *config) error {
	for _, taskCmd := range t.Commands {
		cmdName := "bash"
		cmdArgs := []string{"-c", taskCmd}
		cmd := exec.Command(cmdName, cmdArgs...)

		cmd.Dir = cfg.cwd
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		log.Printf("Executing command: %s", taskCmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
