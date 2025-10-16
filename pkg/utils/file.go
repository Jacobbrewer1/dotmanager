package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
)

func AddFileToWd(localPath string) error {
	_, file := filepath.Split(localPath)
	fileName := strings.Replace(file, ".", "dot_", 1)

	// Copy the file to the repository.
	if err := CopyFile(localPath, fileName); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	return nil
}

func CopyFile(src, dst string) error {
	fmt.Println("Copying", src, "to", dst)

	srcFile, err := os.OpenFile(src, os.O_RDONLY, 0o400)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}

	defer func() {
		_ = srcFile.Close()
	}()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0o500)
	if err != nil {
		return fmt.Errorf("error opening destination file: %w", err)
	}

	defer func() {
		_ = dstFile.Close()
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
}

func GetFileDiff(repoDotPath, homeDotPath string) (string, error) {
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
		fmt.Println("Error reading file:", err)
		return ""
	}

	return string(content)
}
