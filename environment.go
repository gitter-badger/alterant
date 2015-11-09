package main

import "os"

func (m *machine) prepareEnvironment() {
	for variable, value := range m.Environment {
		os.Setenv(variable, value)
	}
}
