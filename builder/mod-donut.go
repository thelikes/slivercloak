package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DoNotAmsiModule struct {
	generateFtnPath string
}

func NewDoNotAmsiModule() *DoNotAmsiModule {
	return &DoNotAmsiModule{
		generateFtnPath: "server/generate/donut.go",
	}
}

func (m *DoNotAmsiModule) Name() string {
	return "donotamsi"
}

func (m *DoNotAmsiModule) Run(config *Config, verbose bool) error {
	// https://github.com/Binject/go-donut/blob/master/main.go#L31
	// s/Bypass:     3,/Bypass:     1,/g
	// s/config.Bypass = 3/config.Bypass = 1/g

	filePath := filepath.Join(config.RunDir, "sliver", m.generateFtnPath)
	// Read the entire file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Perform the replacements
	newContent := strings.ReplaceAll(string(content), "Bypass:     3,", "Bypass:     1,")
	newContent = strings.ReplaceAll(newContent, "config.Bypass = 3", "config.Bypass = 1")

	// Write the modified content back to the file
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
