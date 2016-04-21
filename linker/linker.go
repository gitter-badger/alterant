package linker

import "github.com/autonomy/alterant/link"

// Linker is the interface to a symlink handler
type Linker interface {
	RemoveLinks(map[string]*link.Link) error
	CreateLinks(map[string]*link.Link) error
}
