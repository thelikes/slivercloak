package subs

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SearchAndReplace recursively searchers file content for the searchStr
// and replaces with replaceStr, while preserving file permissions
func SearchAndReplace(rootDir, searchStr, replaceStr string, ignoreDirs []string, verbose bool) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Handle directories
		if info.IsDir() {
			// Check if this directory should be ignored
			baseName := filepath.Base(path)
			for _, ignoreDir := range ignoreDirs {
				if ignoreDir != "" && baseName == ignoreDir {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", path, err)
		}

		// Check if file contains the search string
		if !strings.Contains(string(content), searchStr) {
			return nil
		}

		// Get the original file permissions
		fileMode := info.Mode()

		// Create a temporary file
		tempFile, err := os.CreateTemp(filepath.Dir(path), "temp_*")
		if err != nil {
			return fmt.Errorf("error creating temp file for %s: %v", path, err)
		}
		tempFilePath := tempFile.Name()
		defer os.Remove(tempFilePath) // Clean up in case of failure

		// Process the file line by line to maintain original line endings
		reader := bufio.NewReader(strings.NewReader(string(content)))
		writer := bufio.NewWriter(tempFile)
		modified := false

		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				tempFile.Close()
				return fmt.Errorf("error reading line from %s: %v", path, err)
			}

			if strings.Contains(line, searchStr) {
				line = strings.ReplaceAll(line, searchStr, replaceStr)
				modified = true
			}

			if _, err := writer.WriteString(line); err != nil {
				tempFile.Close()
				return fmt.Errorf("error writing to temp file for %s: %v", path, err)
			}

			if err == io.EOF {
				break
			}
		}

		if err := writer.Flush(); err != nil {
			tempFile.Close()
			return fmt.Errorf("error flushing temp file for %s: %v", path, err)
		}
		tempFile.Close()

		// Only replace the original file if modifications were made
		if modified {
			// Set the same permissions on the temp file before renaming
			if err := os.Chmod(tempFilePath, fileMode); err != nil {
				return fmt.Errorf("error setting permissions on temp file for %s: %v", path, err)
			}

			if err := os.Rename(tempFilePath, path); err != nil {
				return fmt.Errorf("error replacing original file %s: %v", path, err)
			}
			if verbose {
				log.Printf("Modified file: %s\n", path)
			}
		} else {
			// Clean up temp file if no modifications were needed
			os.Remove(tempFilePath)
		}

		return nil
	})
}

// SearchAndRenameFiles recursively searches for and renames files with paths that match
// searchStr, while preserving the original file's file permissions
func SearchAndRenameFiles(rootDir, searchStr, replaceStr string, ignoreDirs []string, verbose bool) error {
	// Walk through all files and directories under rootDir
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's a directory
		if info.IsDir() {
			// Check if this directory should be ignored
			for _, ignoreDir := range ignoreDirs {
				if strings.HasSuffix(path, ignoreDir) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Get the directory and filename separately
		dir, filename := filepath.Split(path)

		// Check if filename contains the search string
		if strings.Contains(filename, searchStr) {
			// Create the new filename with the replaced string
			newFilename := strings.Replace(filename, searchStr, replaceStr, -1)
			newPath := filepath.Join(dir, newFilename)

			// Log the renaming operation
			if verbose {
				log.Printf("Renaming file: %s -> %s", path, newPath)
			}

			// Get the current file permissions
			fileMode := info.Mode()

			// Perform the rename operation
			if err := os.Rename(path, newPath); err != nil {
				return fmt.Errorf("failed to rename file %s to %s: %w", path, newPath, err)
			}

			// Restore the file permissions
			if err := os.Chmod(newPath, fileMode); err != nil {
				return fmt.Errorf("failed to restore permissions for %s: %w", newPath, err)
			}
		}

		return nil
	})
}

// SearchAndRenameDirectories recursively searches for and renames directories that match searchStr
func SearchAndRenameDirectories(rootDir, searchStr, replaceStr string, ignoreDirs []string, verbose bool) error {
	// Clean and get absolute path for proper comparison
	absRootDir, err := filepath.Abs(filepath.Clean(rootDir))
	if err != nil {
		return fmt.Errorf("failed to get absolute path for root directory: %w", err)
	}

	type dirInfo struct {
		path    string
		depth   int
		mode    os.FileMode
		newPath string
	}
	var dirs []dirInfo

	// First pass: collect all directories that need to be renamed
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if this directory should be ignored
		for _, ignoreDir := range ignoreDirs {
			if strings.HasSuffix(path, ignoreDir) {
				return filepath.SkipDir
			}
		}

		// Get absolute path for comparison
		absPath, err := filepath.Abs(filepath.Clean(path))
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		// Skip if this is the root directory itself
		if absPath == absRootDir {
			return nil
		}

		base := filepath.Base(path)
		if strings.Contains(base, searchStr) {
			depth := len(strings.Split(path, string(os.PathSeparator)))
			newName := strings.Replace(base, searchStr, replaceStr, -1)
			newPath := filepath.Join(filepath.Dir(path), newName)

			// Check if destination already exists
			if _, err := os.Stat(newPath); err == nil {
				return fmt.Errorf("destination path already exists: %s", newPath)
			}

			dirs = append(dirs, dirInfo{
				path:    path,
				depth:   depth,
				mode:    info.Mode(),
				newPath: newPath,
			})
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory tree: %w", err)
	}

	// Sort directories by depth in descending order (deepest first)
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].depth > dirs[j].depth
	})

	// Second pass: rename directories from deepest to shallowest
	for _, dir := range dirs {
		if verbose {
			log.Printf("Renaming directory with all contents: %s -> %s", dir.path, dir.newPath)
		}

		// Perform the atomic rename/move operation
		if err := os.Rename(dir.path, dir.newPath); err != nil {
			return fmt.Errorf("failed to rename directory %s to %s: %w", dir.path, dir.newPath, err)
		}

		// Restore the directory permissions
		if err := os.Chmod(dir.newPath, dir.mode); err != nil {
			return fmt.Errorf("failed to restore permissions for %s: %w", dir.newPath, err)
		}
	}

	return nil
}
