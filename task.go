package main

type task struct {
	Links    links    `yaml:"links"`
	Commands []string `yaml:"commands"`
}
