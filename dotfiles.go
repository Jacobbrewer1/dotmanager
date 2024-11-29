package main

import (
	"fmt"
	"os"
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
