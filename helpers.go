package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pmezard/go-difflib/difflib"
)

func isGitRepo() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

func getDiff(repoDotPath, homeDotPath string) (string, error) {
	repoContent, err := os.Open(repoDotPath)
	if err != nil {
		return "", fmt.Errorf("error reading repository dotfile: %w", err)
	}
	defer func() {
		if err := repoContent.Close(); err != nil {
			slog.Warn("Error closing repository dotfile", slog.String(loggingKeyError, err.Error()))
		}
	}()

	homeContent, err := os.Open(homeDotPath)
	if err != nil {
		return "", fmt.Errorf("error reading home dotfile: %w", err)
	}
	defer func() {
		if err := homeContent.Close(); err != nil {
			slog.Warn("Error closing home dotfile", slog.String(loggingKeyError, err.Error()))
		}
	}()

	diff := difflib.UnifiedDiff{
		A: difflib.SplitLines(readFile(repoContent)),
		B: difflib.SplitLines(readFile(homeContent)),
	}

	return difflib.GetUnifiedDiffString(diff)
}

func readFile(f *os.File) string {
	content, err := os.ReadFile(f.Name())
	if err != nil {
		slog.Warn("Error reading file", slog.String(loggingKeyError, err.Error()))
		return ""
	}

	return string(content)
}
