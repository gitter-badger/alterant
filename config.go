package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type config struct {
	Machines  map[string]*machine `yaml:"machines"`
	Tasks     map[string]*task    `yaml:"tasks"`
	Encrypted []string            `yaml:"encrypted"`
}

func newConfig() *config {
	return &config{}
}

func loadConfig(file string) (*config, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := newConfig()

	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func acquireConfig() (*config, error) {
	// require that the config is named 'alter.yaml'
	file := "alter.yaml"

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// require that alter.yaml exists in the cwd
	if _, err := os.Stat(path.Join(cwd, file)); os.IsNotExist(err) {
		return nil, err
	}

	cfg, err := loadConfig(file)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// match the machine's requested tasks to the tasks defined in the config
func (c *config) filterTasks(argMachine string, argTasks []string) error {
	if _, ok := c.Machines[argMachine]; !ok {
		return fmt.Errorf("Machine %s is not defined in alter.yaml", argMachine)
	}

	for mn := range c.Machines {
		// remove irrelevant machines
		if mn != argMachine {
			delete(c.Machines, mn)
			continue
		}
	}

	// remove irrelevant tasks
	auxTasks := map[string]*task{}
	auxRequests := []string{}
	for _, mr := range c.Machines[argMachine].Requests {
		if _, ok := c.Tasks[mr]; ok {
			auxTasks[mr] = c.Tasks[mr]
			auxRequests = append(auxRequests, mr)
		} else {
			return fmt.Errorf("The requested task %s is not defined in alter.yaml", mr)
		}
	}

	c.Tasks = auxTasks
	c.Machines[argMachine].Requests = auxRequests

	// remove tasks not indicated as arguments and check if they tasks are valid
	// for the machine
	if len(argTasks) > 0 {
		auxTasks = map[string]*task{}

		for _, argTask := range argTasks {
			if _, ok := c.Tasks[argTask]; !ok {
				return fmt.Errorf("The requested task %s is not specified for %s in alter.yaml", argTask, argMachine)
			}

			auxTasks[argTask] = c.Tasks[argTask]
		}
	}

	c.Tasks = auxTasks

	return nil
}
