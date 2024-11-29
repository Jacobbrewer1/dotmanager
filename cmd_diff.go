package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"
)

type diffCmd struct{}

func (d *diffCmd) Name() string {
	return "diff"
}

func (d *diffCmd) Synopsis() string {
	return "Diff your local dotfiles with the ones in the repository"
}

func (d *diffCmd) Usage() string {
	return `diff:
  Diff your local dotfiles with the ones in the repository.
`
}

func (d *diffCmd) SetFlags(f *flag.FlagSet) {}

func (d *diffCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Get the dotfiles in the local directory that have the `dot_` prefix.
	// We are assuming that the user is running this command from the root of their repository.

	// Is the current directory a git repository?
	// If not, return an error.
	if !isGitRepo() {
		slog.Error("Current directory is not a git repository")
		return subcommands.ExitFailure
	}

	repoDotFiles, err := getRepositoryDotFiles()
	if err != nil {
		slog.Error("Error getting repository dotfiles", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	for _, repoDot := range repoDotFiles {
		// Is the file in the home directory?
		// If not, show the diff as deleted.
		homeDotName := strings.Replace(repoDot, "dot_", ".", 1)

		// Find the users home directory path.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			slog.Error("Error getting user home directory", slog.String(loggingKeyError, err.Error()))
			return subcommands.ExitFailure
		}

		homeDotPath := filepath.Join(homeDir, homeDotName)
		if _, err := os.Stat(homeDotPath); os.IsNotExist(err) {
			slog.Info("File not found in home directory", slog.String(loggingKeyFile, homeDotPath))
			continue
		}

		// Show the diff between the repository dotfile and the local dotfile.
		// If there is no diff, then skip to the next file.
		// If there is a diff, show the diff.

		// Get the diff between the two files.
		diff, err := getDiff(repoDot, homeDotPath)
		if err != nil {
			slog.Error("Error getting diff", slog.String(loggingKeyError, err.Error()))
			return subcommands.ExitFailure
		}

		if diff == "" {
			slog.Info("No diff found", slog.String(loggingKeyFile, homeDotPath))
			continue
		}

		fmt.Printf("Diff for %s:\n%s\n", homeDotPath, diff)
	}

	return subcommands.ExitSuccess
}
