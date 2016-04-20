package command

import (
	"github.com/autonomy/alterant/cache"
	"gopkg.in/yaml.v2"
)

type Command struct {
	Contents string
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
		SHA1:     cache.SHAFromBytes(b),
	}

	return nil
}
