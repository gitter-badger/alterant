package main

import "log"

type machine struct {
	Environment map[string]string `yaml:"environment"`
	Requests    []string          `yaml:"requests"`
}

func provisionMachine(machine string, cfg *config, flags *provisionFlags) error {
	log.Printf("Provisioning: %s", machine)

	// iterate over the request array to preserve the order of tasks
	for _, request := range cfg.Machines[machine].Requests {
		task := cfg.Tasks[request]
		if task == nil {
			continue
		}

		log.Printf("Performing task: %s", request)

		// export environment variables specific to the specified machine
		prepareEnvironment(machine, cfg.Machines[machine])

		if flags.links {
			// create the links specified in the task
			err := task.createLinks(flags)
			if err != nil {
				return err
			}
		}

		if flags.commands {
			// execute the commands specified in the task
			err := task.executeCommands(cfg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func cleanMachine(machine string, cfg *config) error {
	log.Printf("Cleaning: %s", machine)
	for _, task := range cfg.Tasks {
		for _, link := range task.Links {
			removeLink(link)
		}
	}

	return nil
}
