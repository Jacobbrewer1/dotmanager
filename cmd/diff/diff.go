package diff

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jacobbrewer1/dotmanager/pkg/utils"
)

func PrintDiff(ctx context.Context) error {
	if !utils.IsGitRepo() {
		wd, _ := os.Getwd()
		return fmt.Errorf("current directory (%s) is not a git repository: %w", wd, os.ErrInvalid)
	}

	commonDotFiles, err := utils.CommonDotFiles()
	if err != nil {
		return fmt.Errorf("error getting common dotfiles: %w", err)
	}

	var merr error
	for _, repoDot := range commonDotFiles {
		// Is the file in the home directory?
		// If not, show the diff as deleted.
		homeDotName := strings.Replace(repoDot, "dot_", ".", 1)

		// Find the users home directory path.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			merr = fmt.Errorf("error getting user home directory: %w", err)
			continue
		}

		homeDotPath := filepath.Join(homeDir, homeDotName)
		if _, err := os.Stat(homeDotPath); os.IsNotExist(err) {
			fmt.Println("File does not exist in home directory:", homeDotPath)
			continue
		} else if err != nil {
			merr = fmt.Errorf("error stating file %s: %w", homeDotPath, err)
			continue
		}

		// Show the diff between the repository dotfile and the local dotfile.
		// If there is no diff, then skip to the next file.
		// If there is a diff, show the diff.

		// Get the diff between the two files.
		diff, err := utils.GetFileDiff(repoDot, homeDotPath)
		if err != nil {
			merr = fmt.Errorf("error getting diff for %s and %s: %w", repoDot, homeDotPath, err)
			continue
		}

		if diff == "" {
			continue
		}

		fmt.Printf("\n")
		fmt.Printf("Diff for %s:\n%s\n", homeDotPath, diff)
	}

	if merr != nil {
		return merr
	}

	return nil
}
