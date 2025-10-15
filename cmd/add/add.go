package add

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"

	"github.com/jacobbrewer1/dotmanager/pkg/selectors"
	"github.com/jacobbrewer1/dotmanager/pkg/utils"
)

func Files(ctx context.Context) error {
	if !utils.IsGitRepo() {
		wd, _ := os.Getwd()
		return fmt.Errorf("current directory (%s) is not a git repository: %w", wd, os.ErrInvalid)
	}

	// Find the users home directory path.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	// Get all files in the home directory.
	files, err := os.ReadDir(homeDir)
	if err != nil {
		return fmt.Errorf("reading home directory: %w", err)
	}

	availableFiles := make([]string, 0)
	for _, file := range files {
		switch {
		case file.IsDir():
			continue
		case !strings.HasPrefix(file.Name(), "."):
			continue
		case strings.HasSuffix(file.Name(), ".tmp"):
			continue
		default:
			availableFiles = append(availableFiles, file.Name())
		}
	}

	choices, err := selectors.UserSelectionForm("Select files to add", availableFiles)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("Operation cancelled")
			return nil
		}
		return fmt.Errorf("getting user selection: %w", err)
	}

	var merr error
	for _, choice := range choices {
		// See if the file is already being tracked.
		// If it is, return an error.
		if stat, err := os.Stat(strings.Replace(choice, ".", "dot_", 1)); !os.IsNotExist(err) {
			merr = multierr.Append(merr, fmt.Errorf("the file %s is already being tracked", choice))
			continue
		} else if err != nil && !os.IsNotExist(err) {
			merr = multierr.Append(merr, fmt.Errorf("checking if the file %s is already being tracked: %w", choice, err))
			continue
		} else if stat != nil && stat.Name() != "" {
			merr = multierr.Append(merr, fmt.Errorf("the file %s is already being tracked", choice))
			continue
		}

		// Get the absolute path of the file in the home directory.
		homeDotPath := filepath.Join(homeDir, choice)

		// Add the file to the repository.
		if err := utils.AddFileToWd(homeDotPath); err != nil {
			merr = multierr.Append(merr, fmt.Errorf("adding file %s to repository: %w", choice, err))
			continue
		}

		fmt.Println("Added", choice)
	}

	if merr != nil {
		return merr
	}

	return nil
}
