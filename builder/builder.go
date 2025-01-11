package main

import (
	"fmt"
	"log"
)

// Module defines the interface for build modules.
// Each module represents a discrete step in the build process and must provide:
// - A Name() method to identify the module
// - A Run() method to execute the module's functionality
type Module interface {
	Name() string
	Run(config *Config, verbose bool) error
}

// Builder orchestrates the build process by managing a collection of modules.
// It provides a centralized way to configure and execute multiple build steps
type Builder struct {
	modules map[string]Module
	config  *Config
	verbose bool // Controls command output display
}

// NewBuilder creates a new Builder instance with the provided configuration.
// It initializes an empty module registry and sets global build settings.
//
// Parameters:
//   - config: Configuration settings for the build process
//   - verbose: Whether to display detailed command output
//
// Returns:
//   - A pointer to the newly created Builder
func NewBuilder(config *Config, verbose bool) *Builder {
	return &Builder{
		modules: make(map[string]Module),
		config:  config,
		verbose: verbose,
	}
}

// RegisterModule adds a new module to the builder's registry.
// Each module is stored in the modules map using its name as the key.
// If a module with the same name already exists, it will be overwritten.
//
// Parameters:
//   - m: The module to register, must implement the Module interface
func (b *Builder) RegisterModule(m Module) {
	b.modules[m.Name()] = m
}

// Run executes the complete build process in three main steps:
// 1. Clones the repository
// 2. Runs specified modules (if any)
// 3. Executes make commands for compilation
//
// Parameters:
//   - moduleNames: Slice of module names to execute. If empty or ["all"],
//     only repo cloning and compilation will be performed
//
// Returns:
//   - error: If any step in the build process fails
//
// The function follows a sequential process where:
// - First, the repository is always cloned
// - Then, if specific modules are requested, they are executed in order
// - Finally, make commands are run to compile the project
func (b *Builder) Run(moduleNames []string) error {
	log.Println("Cloning Sliver...")
	if err := b.cloneRepo(); err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	// Handle module execution
	if len(moduleNames) > 0 {
		if moduleNames[0] == "all" {
			// Get all registered module names
			allModules := make([]string, 0, len(b.modules))
			for name := range b.modules {
				allModules = append(allModules, name)
			}
			if err := b.runModules(allModules); err != nil {
				return err
			}
		} else {
			if err := b.runModules(moduleNames); err != nil {
				return err
			}
		}
	}

	// Run make commands
	log.Println("Compiling...")
	if err := b.runMake(); err != nil {
		return err
	}

	return nil
}

// runModules executes a sequence of modules in the order specified.
// It validates each module's existence before execution and handles errors.
//
// Parameters:
//   - moduleNames: A slice of module names to execute in order
//
// Returns:
//   - error: If a module is not found or if any module execution fails
//
// The function will stop execution and return an error immediately if:
//   - A requested module is not found in the registry
//   - Any module's Run() method returns an error
func (b *Builder) runModules(moduleNames []string) error {
	for _, name := range moduleNames {
		module, exists := b.modules[name]
		if !exists {
			return fmt.Errorf("module %s not found", name)
		}

		log.Println("Running module:", module.Name())
		if err := module.Run(b.config, b.verbose); err != nil {
			return fmt.Errorf("module %s failed: %w", name, err)
		}
	}
	return nil
}
