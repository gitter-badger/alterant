package linker

import "github.com/autonomy/alterant/link"

// Linker is the interface to a symlink handler
type Linker interface {
	RemoveLinks([]*link.Link) error
	CreateLinks([]*link.Link) error
}
