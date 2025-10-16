package pull

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/multierr"

	"github.com/jacobbrewer1/dotmanager/pkg/utils"
)

func Files(ctx context.Context) error {
	if !utils.IsGitRepo() {
		wd, _ := os.Getwd()
		return fmt.Errorf("current directory (%s) is not a git repository: %w", wd, os.ErrInvalid)
	}

	commonDotFiles, err := utils.CommonDotFiles()
	if err != nil {
		return fmt.Errorf("failed to get common dotfiles list: %w", err)
	}

	var merr error
	for _, repoDot := range commonDotFiles {
		// Is the file in the home directory?
		// If not, show the diff as deleted.
		homeDotName := repoDot
		if len(repoDot) > 4 && repoDot[:4] == "dot_" {
			homeDotName = "." + repoDot[4:]
		}

		// Find the users home directory path.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			merr = multierr.Append(merr, fmt.Errorf("error getting user home directory: %w", err))
			continue
		}

		homeDotPath := filepath.Join(homeDir, homeDotName)
		if _, err := os.Stat(homeDotPath); os.IsNotExist(err) {
			fmt.Println("File does not exist in home directory:", homeDotPath)
			continue
		} else if err != nil {
			merr = multierr.Append(merr, fmt.Errorf("error stating file %s: %w", homeDotPath, err))
			continue
		}

		// Are there any differences between the two files?
		diff, err := utils.GetFileDiff(repoDot, homeDotPath)
		if err != nil {
			merr = multierr.Append(merr, fmt.Errorf("error getting diff for %s and %s: %w", repoDot, homeDotPath, err))
			continue
		} else if diff == "" {
			// No differences, skip to the next file.
			continue
		}

		// Copy the file from the home directory to the repository.
		if err := utils.CopyFile(homeDotPath, repoDot); err != nil {
			merr = multierr.Append(merr, fmt.Errorf("error copying file from %s to %s: %w", homeDotPath, repoDot, err))
			continue
		}

		fmt.Printf("Pulled changes from %s to %s\n", homeDotPath, repoDot)
	}
	if merr != nil {
		return merr
	}

	return nil
}
