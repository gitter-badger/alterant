package command

import (
	"github.com/autonomy/alterant/hasher"
	"gopkg.in/yaml.v2"
)

type Command struct {
	Contents string
	Queued   bool
	SHA1     string
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux string

	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(&aux)
	if err != nil {
		return err
	}

	*c = Command{
		Contents: aux,
		Queued:   true,
		SHA1:     hasher.SHA1FromBytes(b),
	}

	return nil
}
