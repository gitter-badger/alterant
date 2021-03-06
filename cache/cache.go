package cache

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/autonomy/alterant/config"
	"github.com/autonomy/alterant/task"
)

type Task struct {
	Links    []string `yaml:"links"`
	Commands []string `yaml:"commands"`
	SHA1     string   `yaml:"sha1"`
}

type Machine struct {
	Tasks map[string]Task `yaml:"tasks"`
	SHA1  string          `yaml:"sha1"`
}

type Cache struct {
	Machines map[string]Machine `yaml:"machines"`
}

func (c *Cache) AddTask(machine Machine, task *task.Task) {
	t := Task{}

	for _, link := range task.Links {
		t.Links = append(t.Links, link.SHA1)
	}

	for _, command := range task.Commands {
		t.Commands = append(t.Commands, command.SHA1)
	}

	t.SHA1 = task.SHA1

	machine.Tasks[task.Name] = t
}

func WriteToFile(cfg *config.Config) error {
	cache := Cache{}
	m := Machine{}

	cache.Machines = make(map[string]Machine)
	m.Tasks = make(map[string]Task)

	for _, task := range cfg.Tasks {
		cache.AddTask(m, task)
	}

	m.SHA1 = cfg.SHA1

	cache.Machines[cfg.Machine] = m

	d, err := yaml.Marshal(&cache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.Getenv("ALTERANT_HOME")+"/db.yaml", d, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadCache() (*Cache, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("ALTERANT_HOME") + "/db.yaml")
	if err != nil {
		return nil, err
	}

	c := &Cache{}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
