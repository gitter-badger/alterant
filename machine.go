package main

import (
	"log"

	"github.com/andrewrynhard/go-ordered-map"
)

type taskMap struct {
	Map *ordered.OrderedMap
}

type machine struct {
	Environment map[string]string `yaml:"environment"`
	Requests    []string          `yaml:"requests"`
}

func provisionMachine(tasks []*task, a *alterant, flags *provisionFlags) error {
	log.Printf("Provisioning: %s", a.machineName)

	for _, task := range tasks {
		// log.Printf("Performing task: %s", task.name)

		// export environment variables specific to the specified machine
		prepareEnvironment(a.machineName, a.machinePtr)

		if flags.links {
			// create the links specified in the task
			err := task.createLinks(flags)
			if err != nil {
				return err
			}
		}

		if flags.commands {
			// execute the commands specified in the task
			err := task.executeCommands(a.cfg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO: clean this up
func cleanMachine(machine string, tasks []*task, cfg *config) error {
	// for _, task := range tasks {
	// 	for _, link := range task.Links {
	// 		if removeTasks != nil {
	// 			if _, ok := removeTasks[task]; ok {
	// 				err := link.removeLink()
	// 				if err != nil {
	// 					return err
	// 				}
	// 			} else {
	// 				continue
	// 			}
	// 		} else {
	// 			err := link.removeLink()
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }

	return nil
}
