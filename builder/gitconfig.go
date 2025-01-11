package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type BuildTarget struct {
	Tag         string
	GitRef      string
	UseGoLegacy bool // true for 1.5 which needs Go 1.18
}

type Config struct {
	RepoURL string
	RunDir  string // Path to current run directory
	Target  BuildTarget
}

// NewConfig sets up the run directory and repo targets
func NewConfig(targetVersion string) (*Config, error) {
	// Create unique run directory
	timestamp := time.Now().Format("20060102_150405")
	runDir := filepath.Join("/tmp/output", fmt.Sprintf("run_%s_%s", targetVersion, timestamp))
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create run directory: %w", err)
	}

	// Map target version to its configuration
	targets := map[string]BuildTarget{
		"1.5": {
			Tag:         Version1_5,
			GitRef:      Version1_5,
			UseGoLegacy: true,
		},
		"1.6": {
			Tag:         Version1_6,
			GitRef:      Version1_6,
			UseGoLegacy: false,
		},
	}

	target, exists := targets[targetVersion]
	if !exists {
		return nil, fmt.Errorf("invalid target version: %s", targetVersion)
	}

	return &Config{
		RepoURL: RepoURL,
		RunDir:  runDir,
		Target:  target,
	}, nil
}

func (b *Builder) cloneRepo() error {
	// Clone into the run directory
	cmd := exec.Command("git", "clone", b.config.RepoURL)
	cmd.Dir = b.config.RunDir
	if b.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Handle tag checkout for 1.5
	if b.config.Target.UseGoLegacy {
		repoDir := filepath.Join(b.config.RunDir, "sliver")

		cmd = exec.Command("git", "fetch", "--all", "--tags")
		cmd.Dir = repoDir
		if b.verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to fetch tags: %w", err)
		}

		cmd = exec.Command("git", "checkout", "tags/"+b.config.Target.Tag)
		cmd.Dir = repoDir
		if b.verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to checkout tag: %w", err)
		}
		if b.verbose {
			log.Printf("Successfully checked out tag %s", b.config.Target.Tag)
		}
	}

	return nil
}
