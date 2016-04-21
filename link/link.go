package link

import (
	"os"
	"path"

	"gopkg.in/yaml.v2"

	"github.com/autonomy/alterant/cache"
)

// SymlinkTarget is a custom type for symlink targets
type SymlinkTarget string

// SymlinkDestination is a custom type for symlink destinations
type SymlinkDestination string

// Link represents a link in the machine yaml
type Link struct {
	Target      SymlinkTarget
	Destination SymlinkDestination
	Encrypted   bool
	SHA1        string
}

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

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (l *Link) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Target      SymlinkTarget      `yaml:"target"`
		Destination SymlinkDestination `yaml:"destination"`
		Encrypted   bool               `yaml:"encrypted"`
	}

	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(&aux)
	if err != nil {
		return err
	}

	*l = Link{
		Target:      aux.Target,
		Destination: aux.Destination,
		Encrypted:   aux.Encrypted,
		SHA1:        cache.SHAFromBytes(b),
	}

	return nil
}
