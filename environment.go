package main

import "os"

func (m *machine) prepareEnvironment(machine string) {
	// export the machine name to the environment
	os.Setenv("MACHINE", machine)

	for variable, value := range m.Environment {
		os.Setenv(variable, value)
	}
}
