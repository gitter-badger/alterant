package task

import "github.com/autonomy/alterant/linker"

// Task represents a task
type Task struct {
	Links    map[linker.SymlinkTarget]linker.SymlinkDestination `yaml:"links"`
	Commands []string                                           `yaml:"commands"`
}
