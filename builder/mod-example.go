package main

import (
	"log"
)

type ExampleModule struct {
	// Add any module-specific configuration here
}

func NewExampleModule() *ExampleModule {
	return &ExampleModule{}
}

func (m *ExampleModule) Name() string {
	return "example"
}

func (m *ExampleModule) Run(config *Config, verbose bool) error {
	if verbose {
		log.Println("hey from example module \\o/")
	}

	// do stuff

	return nil
}
