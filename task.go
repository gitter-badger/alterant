package main

import (
	"errors"
	"log"
)

type task struct {
	Destination string            `yaml:"destination"`
	Links       map[string]string `yaml:"links"`
	Commands    []string          `yaml:"commands"`
}

func provisionMachine(machine string, cfg *config, flags *provisionFlags) error {
	m := cfg.Machines[machine]
	for _, task := range m.Tasks {
		t, ok := cfg.Tasks[task]
		if !ok {
			return errors.New("Task " + task + " is not defined in alter.yaml")
		}

		log.Printf("Performing task: %s", task)

		// export environment variables specific to the specified machine
		m.prepareEnvironment(machine)

		if flags.links {
			// create the links specified in the task
			err := t.createLinks(cfg.cwd, flags)
			if err != nil {
				return err
			}
		}

		if flags.commands {
			// execute the commands specified in the task
			err := t.executeCommands(cfg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
