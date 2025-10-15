package utils

import (
	"fmt"
	"os"
	"strings"
)

// IsGitRepo checks if the current directory is a git repository
func IsGitRepo() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

// CommonDotFiles returns a list of all the dot files that are in both the current and the home directory.
func CommonDotFiles() ([]string, error) {
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
