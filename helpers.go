package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
)

func isGitRepo() bool {
	_, err := os.Stat(".git")
	return !os.IsNotExist(err)
}

func getDiff(repoDotPath, homeDotPath string) (string, error) {
	repoContent, err := os.OpenFile(repoDotPath, os.O_RDONLY, 0o400)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}

	homeContent, err := os.OpenFile(homeDotPath, os.O_RDONLY, 0o400)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}

	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(readFile(repoContent)),
		B:       difflib.SplitLines(readFile(homeContent)),
		Context: 3,
	}

	diffStr, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", fmt.Errorf("error creating diff: %w", err)
	}

	return colorizeDiff(diffStr), nil
}

func colorizeDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "+") {
			lines[i] = color.GreenString(line)
		} else if strings.HasPrefix(line, "-") {
			lines[i] = color.RedString(line)
		}
	}
	return strings.Join(lines, "\n")
}

func readFile(f *os.File) string {
	content, err := os.ReadFile(f.Name())
	if err != nil {
		slog.Warn("Error reading file", slog.String(loggingKeyError, err.Error()))
		return ""
	}

	return string(content)
}
