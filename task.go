package main

type task struct {
	Links    map[linkTarget]linkDestination `yaml:"links"`
	Commands []string                       `yaml:"commands"`
	name     string
	machine  *machine
}
