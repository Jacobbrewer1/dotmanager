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
	fmt.Println("Copying", src, "to", dst)

	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0o400)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0o500)
	if err != nil {
		return fmt.Errorf("error opening destination file: %w", err)
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
}

func addFile(localPath string) error {
	_, file := filepath.Split(localPath)
	fileName := strings.Replace(file, ".", "dot_", 1)

	// Copy the file to the repository.
	if err := copyFile(localPath, fileName); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	return nil
}
