package main

import "os"

func prepareEnvironment(m *machine) {
	// export the machine name to the environment
	os.Setenv("MACHINE", m.name)

	for variable, value := range m.Environment {
		os.Setenv(variable, value)
	}
}
