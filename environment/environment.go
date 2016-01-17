package environment

import (
	"os"

	"github.com/autonomy/alterant/logger"
)

// Environment represents the environment variables in alter.yaml
type Environment struct {
	logger  *logWrapper.LogWrapper
	machine string
}

// Set exports the variables to the environment
func (e *Environment) Set(environment map[string]string) {
	for variable, value := range environment {
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
