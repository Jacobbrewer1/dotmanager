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

type pushCmd struct{}

func (p *pushCmd) Name() string {
	return "push"
}

func (p *pushCmd) Synopsis() string {
	return "Push the local dotfiles into the repository"
}

func (p *pushCmd) Usage() string {
	return `push:
  Push the repository dotfiles into the local directory.
`
}

func (p *pushCmd) SetFlags(f *flag.FlagSet) {}

func (p *pushCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

		// Get the absolute path of the file in the home directory.
		homeDotPath := filepath.Join(homeDir, homeDotName)

		// Is there any difference between the home dotfile and the repository dotfile?
		diff, err := getDiff(repoDot, homeDotPath)
		if err != nil {
			slog.Error("Error getting diff", slog.String(loggingKeyError, err.Error()))
			return subcommands.ExitFailure
		} else if diff == "" {
			continue
		}

		slog.Info("Pushing file", slog.String(loggingKeyFile, homeDotPath))

		// Copy the file from the repository to the home directory.
		if err := copyFile(repoDot, homeDotPath); err != nil {
			slog.Error("Error copying file", slog.String(loggingKeyError, err.Error()))
			return subcommands.ExitFailure
		}
	}

	slog.Info("Push complete")

	return subcommands.ExitSuccess
}
