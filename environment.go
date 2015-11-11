package main

import (
	"log"
	"os"
)

func prepareEnvironment(mn string, mp *machine) {
	for variable, value := range mp.Environment {
		log.Printf("Exporting %s: %s", os.ExpandEnv(variable), os.ExpandEnv(value))
		os.Setenv(os.ExpandEnv(variable), os.ExpandEnv(value))
	}
}
