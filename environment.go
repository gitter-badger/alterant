package main

import (
	"log"
	"os"
)

func prepareEnvironment(mn string, mp *machine) {
	for variable, value := range mp.Environment {
		log.Printf("Exporting %s: %s", variable, value)
		os.Setenv(variable, value)
	}
}
