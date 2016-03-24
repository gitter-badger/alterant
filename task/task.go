package task

import "github.com/autonomy/alterant/linker"

// Task represents a task
type Task struct {
	Links    []*linker.Link `yaml:"links"`
	Commands []string       `yaml:"commands"`
}
