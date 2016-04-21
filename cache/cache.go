package cache

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/autonomy/alterant/config"
)

type Task struct {
	Links    []string `yaml:"links"`
	Commands []string `yaml:"commands"`
	SHA1     string   `yaml:"sha1"`
}

type Machine struct {
	Machine string          `yaml:"machine"`
	Tasks   map[string]Task `yaml:"tasks"`
	SHA1    string          `yaml:"sha1"`
}

type Cache struct {
	Machines []Machine
}

func WriteToFile(cfg *config.Config) error {
	cache := Cache{}
	m := Machine{}
	m.Machine = cfg.Machine

	m.Tasks = make(map[string]Task)

	for _, task := range cfg.Tasks {
		t := Task{}

		for _, link := range task.Links {
			t.Links = append(t.Links, link.SHA1)
		}

		for _, command := range task.Commands {
			t.Commands = append(t.Commands, command.SHA1)
		}

		t.SHA1 = task.SHA1
		m.Tasks[task.Name] = t
	}

	m.SHA1 = cfg.SHA1

	cache.Machines = append(cache.Machines, m)

	d, err := yaml.Marshal(&cache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/tmp/db.yaml", d, 0644)
	if err != nil {
		return err
	}

	return nil
}
