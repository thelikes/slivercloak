package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//
// Cruft code from trying to use a single dockerfile and
// manage the Go environment from the entrypoint
//

func cleanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func cp(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return os.Chmod(dst, 0755)
}

func setupGoVersion(targetVersion string) error {
	// Save current Go if it's not already saved
	if _, err := os.Stat("/usr/local/go.modern"); err != nil {
		if err := os.Rename("/usr/local/go", "/usr/local/go.modern"); err != nil {
			return fmt.Errorf("failed to save modern go: %w", err)
		}
	}

	if targetVersion == "1.5" {
		// Switch to Go 1.18
		if err := os.RemoveAll("/usr/local/go"); err != nil {
			return fmt.Errorf("failed to remove current go: %w", err)
		}
		if err := os.Rename("/usr/local/go1.18.temp", "/usr/local/go"); err != nil {
			return fmt.Errorf("failed to move go 1.18 into place: %w", err)
		}

		// Switch GOPATH/bin to legacy tools
		if err := cleanDir("/go/bin"); err != nil {
			return fmt.Errorf("failed to clean bin directory: %w", err)
		}
		if err := os.Remove("/go/bin"); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove bin symlink: %w", err)
		}
		if err := os.Symlink("/go/bin.legacy", "/go/bin"); err != nil {
			return fmt.Errorf("failed to symlink legacy bin: %w", err)
		}

		// Copy legacy protoc
		if err := cp("/usr/local/protoc-legacy/bin/protoc", "/usr/local/bin/protoc"); err != nil {
			return fmt.Errorf("failed to copy legacy protoc: %w", err)
		}
	} else {
		// Switch back to modern Go
		if err := os.RemoveAll("/usr/local/go"); err != nil {
			return fmt.Errorf("failed to remove current go: %w", err)
		}
		if err := os.Rename("/usr/local/go.modern", "/usr/local/go"); err != nil {
			return fmt.Errorf("failed to restore modern go: %w", err)
		}

		// Switch GOPATH/bin to modern tools
		if err := cleanDir("/go/bin"); err != nil {
			return fmt.Errorf("failed to clean bin directory: %w", err)
		}
		if err := os.Remove("/go/bin"); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove bin symlink: %w", err)
		}
		if err := os.Symlink("/go/bin.modern", "/go/bin"); err != nil {
			return fmt.Errorf("failed to symlink modern bin: %w", err)
		}

		// Copy modern protoc
		if err := cp("/usr/local/protoc-modern/bin/protoc", "/usr/local/bin/protoc"); err != nil {
			return fmt.Errorf("failed to copy modern protoc: %w", err)
		}
	}

	// Verify the switch worked
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to verify go version: %w", err)
	}
	log.Printf("Using Go version: %s", string(output))

	return nil
}
