package machine

// Machine represents a machine in alter.yaml
type Machine struct {
	Environment map[string]string `yaml:"environment"`
	Requests    []string          `yaml:"requests"`
}
