package main

type machine struct {
	Environment map[string]string `yaml:"environment"`
	Tasks       []string          `yaml:"tasks"`
}
