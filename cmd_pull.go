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

type pullCmd struct{}

func (p *pullCmd) Name() string {
	return "pull"
}

func (p *pullCmd) Synopsis() string {
	return "Pull the local dotfile changes into the repository"
}

func (p *pullCmd) Usage() string {
	return `pull:
  Pull the local dotfile changes into the repository.
`
}

func (p *pullCmd) SetFlags(f *flag.FlagSet) {}

func (p *pullCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

		// Is there any difference between the home dotfile and the repository dotfile?
		diff, err := getDiff(repoDot, homeDotPath)
		if err != nil {
			slog.Error("Error getting diff", slog.String(loggingKeyError, err.Error()))
			return subcommands.ExitFailure
		} else if diff == "" {
			continue
		}

		slog.Info("Copying file", slog.String(loggingKeyFile, homeDotName))

		// Copy the contents of the home dotfile to the repository dotfile.
		if err := copyFile(homeDotPath, repoDot); err != nil {
			slog.Error("Error copying file", slog.String(loggingKeyError, err.Error()))
			return subcommands.ExitFailure
		}
	}

	slog.Info("Pull complete")

	return subcommands.ExitSuccess
}
