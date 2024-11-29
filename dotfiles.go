package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getRepositoryDotFiles() ([]string, error) {
	// Get all the files in the repository that have the `dot_` prefix.
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	dotFiles := make([]string, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "dot_") {
			dotFiles = append(dotFiles, file.Name())
		}
	}

	return dotFiles, nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			fmt.Printf("Error closing source file: %s\n", err)
		}
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			fmt.Printf("Error closing destination file: %s\n", err)
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
}

func addFile(localPath string) error {
	// Get the file name from the path.
	fileName := filepath.Base(localPath)
	fileName = strings.Replace(fileName, "dot_", ".", 1)

	// Copy the file to the repository.
	if err := copyFile(localPath, fileName); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	return nil
}
