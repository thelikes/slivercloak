package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (b *Builder) runMake() error {
	makeDir := filepath.Join(b.config.RunDir, "sliver")
	if _, err := os.Stat(makeDir); err != nil {
		return fmt.Errorf("make directory not found: %w", err)
	}

	// First run 'make pb'
	cmd := exec.Command("make", "pb")
	cmd.Dir = makeDir
	if b.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make pb failed: %w", err)
	}

	// Then run 'make'
	cmd = exec.Command("make")
	cmd.Dir = makeDir
	if b.verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make failed: %w", err)
	}

	return nil
}
