package task

import "github.com/autonomy/alterant/linker"

type link struct {
	Target      linker.SymlinkTarget      `yaml:"target"`
	Destination linker.SymlinkDestination `yaml:"destination"`
	Encrypted   bool                      `yaml:"encrypted"`
}

// Task represents a task
type Task struct {
	Links    []link   `yaml:"links"`
	Commands []string `yaml:"commands"`
}
