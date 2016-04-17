package commander

import (
	"os"
	"os/exec"

	"github.com/autonomy/alterant/logger"
	"github.com/autonomy/alterant/task"
)

//DefaultCommander is a basic command executer and is the default
type DefaultCommander struct {
	logger  *logWrapper.LogWrapper
	enabled bool
}

// Execute executes the command on the system
func (dc *DefaultCommander) Execute(t *task.Task) error {
	// do not execute commands if not enabled
	if !dc.enabled {
		return nil
	}

	for _, taskCmd := range t.Commands {
		cmdName := "bash"
		cmdArgs := []string{"-c", taskCmd}
		cmd := exec.Command(cmdName, cmdArgs...)

		if dc.logger.Verbose {
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
		}

		dc.logger.Info(2, "Executing command: \n%s", taskCmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// NewDefaultCommander returns an instance of `DefaultLinker`
func NewDefaultCommander(enabled bool, logger *logWrapper.LogWrapper) *DefaultCommander {
	return &DefaultCommander{
		logger:  logger,
		enabled: enabled,
	}
}
