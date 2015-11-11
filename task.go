package main

type task struct {
	Links    map[symlinkTarget]symlinkDestination `yaml:"links"`
	Commands []string                             `yaml:"commands"`
}
