package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"
)

type addCmd struct{}

func (a *addCmd) Name() string {
	return "add"
}

func (a *addCmd) Synopsis() string {
	return "Start tracking a new dotfile"
}

func (a *addCmd) Usage() string {
	return `add:
  Start tracking a new dotfile.
`
}

func (a *addCmd) SetFlags(f *flag.FlagSet) {}

func (a *addCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Get the dotfiles in the local directory that have the `dot_` prefix.
	// We are assuming that the user is running this command from the root of their repository.

	// Is the current directory a git repository?
	// If not, return an error.
	if !isGitRepo() {
		slog.Error("Current directory is not a git repository")
		return subcommands.ExitFailure
	}

	// Find the users home directory path.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Error getting user home directory", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	// Get all files in the home directory.
	files, err := os.ReadDir(homeDir)
	if err != nil {
		slog.Error("Error reading the home directory", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	availableFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else if !strings.HasPrefix(file.Name(), ".") {
			continue
		} else if strings.HasSuffix(file.Name(), ".tmp") {
			continue
		}

		availableFiles = append(availableFiles, file.Name())
	}

	selector := newFileSelector(availableFiles)
	choice, err := selector.Exec()
	if err != nil {
		slog.Error("Error selecting file", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	// See if the file is already being tracked.
	// If it is, return an error.
	if stat, err := os.Stat(strings.Replace(choice, ".", "dot_", 1)); !os.IsNotExist(err) {
		slog.Error("The file is already being tracked", slog.String(loggingKeyFile, choice))
		return subcommands.ExitFailure
	} else if err != nil && !os.IsNotExist(err) {
		slog.Error("Error checking if the file is already being tracked", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	} else if stat != nil && stat.Name() != "" {
		slog.Error("The file is already being tracked", slog.String(loggingKeyFile, choice))
		return subcommands.ExitFailure
	}

	// Get the absolute path of the file in the home directory.
	homeDotPath := filepath.Join(homeDir, choice)

	// Add the file to the repository.
	if err := addFile(homeDotPath); err != nil {
		slog.Error("Error adding file to the repository", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	slog.Info("File added to the repository", slog.String(loggingKeyFile, homeDotPath))

	return subcommands.ExitSuccess
}
