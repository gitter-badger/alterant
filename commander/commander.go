package commander

import "github.com/autonomy/alterant/task"

// Commander is the interface for command execution
type Commander interface {
	Execute(*task.Task) error
}
