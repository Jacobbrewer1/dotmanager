package push

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		homeDotName := strings.Replace(repoDot, "dot_", ".", 1)

		// Find the users home directory path.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			merr = multierr.Append(merr, fmt.Errorf("error getting user home directory: %w", err))
			continue
		}

		homeDotPath := filepath.Join(homeDir, homeDotName)

		diff, err := utils.GetFileDiff(repoDot, homeDotPath)
		if err != nil {
			merr = multierr.Append(merr, fmt.Errorf("error getting diff for %s and %s: %w", repoDot, homeDotPath, err))
			continue
		} else if diff == "" {
			continue
		}

		// Copy the file from the repository to the home directory.
		if err := utils.CopyFile(repoDot, homeDotPath); err != nil {
			merr = multierr.Append(merr, fmt.Errorf("error copying file from %s to %s: %w", repoDot, homeDotPath, err))
			continue
		}

		fmt.Println("Pushed changes from", repoDot, "to", homeDotPath)
	}

	if merr != nil {
		return merr
	}

	return nil
}
