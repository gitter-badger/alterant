package linker

import (
	"os"
	"path"
)

// Linker is the interface to a symlink handler
type Linker interface {
	RemoveLinks(map[SymlinkTarget]SymlinkDestination) error
	CreateLinks([]Link) error
}

// SymlinkTarget is a custom type for symlink targets
type SymlinkTarget string

// SymlinkDestination is a custom type for symlink destinations
type SymlinkDestination string

// UnmarshalYAML implementation for SymlinkTarget
func (t *SymlinkTarget) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux string

	if err := unmarshal(&aux); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	*t = SymlinkTarget(path.Join(cwd, os.ExpandEnv(aux)))

	return nil
}

// UnmarshalYAML implementation for SymlinkDestination
func (t *SymlinkDestination) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux string

	if err := unmarshal(&aux); err != nil {
		return err
	}

	home := os.Getenv("HOME")

	*t = SymlinkDestination(path.Join(home, os.ExpandEnv(aux)))

	return nil
}
