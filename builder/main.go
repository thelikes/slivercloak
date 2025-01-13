package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

// Configuration constants
const (
	RepoURL    = "https://github.com/BishopFox/sliver.git"
	Version1_5 = "v1.5.42" // tag
	Version1_6 = "master"  // branch
)

func main() {
	// Get target version from environment, fallback to flag if not set
	//targetVersion := flag.String("target", "1.5", "Target version (1.5 or 1.6)")
	targetVersion := os.Getenv("TARGET_VERSION")
	if targetVersion == "" {
		targetVersionFlag := flag.String("target", "1.5", "Target version (1.5 or 1.6)")
		flag.Parse()
		targetVersion = *targetVersionFlag
	}

	modules := flag.String("modules", "", "Comma-separated list of modules to run")
	verbose := flag.Bool("verbose", false, "Show build output")
	flag.Parse()

	// Create the run environment
	config, err := NewConfig(targetVersion)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	log.Println("Target version:", config.Target.Tag)
	log.Println("Run directory:", config.RunDir)

	// create our builder
	builder := NewBuilder(config, *verbose)

	// Register the 'example' module
	exampleModule := NewExampleModule()
	builder.RegisterModule(exampleModule)

	// Register the 'branding' module
	brandingModule := NewBrandingModule()
	builder.RegisterModule(brandingModule)

	// Register the 'donOtamsi' module
	donOtamsi := NewDoNotAmsiModule()
	builder.RegisterModule(donOtamsi)

	// Register the 'elastic' module
	elasticModule := NewElasticModule()
	builder.RegisterModule(elasticModule)

	// Register new modules here
	// ...

	// process user input list
	var moduleList []string
	if *modules != "" {
		moduleList = strings.Split(*modules, ",")
	}

	// clone, process, build
	if err := builder.Run(moduleList); err != nil {
		log.Fatal(err)
	}
}
