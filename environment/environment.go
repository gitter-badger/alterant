package environment

import (
	"os"

	"github.com/autonomy/alterant/logger"
	"github.com/autonomy/alterant/machine"
)

// Environment represents the environment variables in alter.yaml
type Environment struct {
	logger  *logWrapper.LogWrapper
	machine string
}

// Set exports the variables to the environment
func (e *Environment) Set(mn string, mp *machine.Machine) {
	for variable, value := range mp.Environment {
		e.logger.Info("Exporting %s: %s", os.ExpandEnv(variable), os.ExpandEnv(value))
		os.Setenv(os.ExpandEnv(variable), os.ExpandEnv(value))
	}
}

// NewEnvironment returns an instance of `Environment`
func NewEnvironment(machine string, logger *logWrapper.LogWrapper) *Environment {
	return &Environment{
		logger:  logger,
		machine: machine,
	}
}
