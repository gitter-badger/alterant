package task

import "github.com/autonomy/alterant/linker"

// Task represents a task
type Task struct {
	Dependencies []string       `yaml:"dependencies"`
	Links        []*linker.Link `yaml:"links"`
	Commands     []string       `yaml:"commands"`
	Name         string
}
