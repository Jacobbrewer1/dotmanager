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

type addCmd struct {
	// localName is the name of the local file.
	localName string
}

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

func (a *addCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&a.localName, "local", "", "The name of the local file within your home directory.")
}

func (a *addCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Get the dotfiles in the local directory that have the `dot_` prefix.
	// We are assuming that the user is running this command from the root of their repository.

	// Is the current directory a git repository?
	// If not, return an error.
	if !isGitRepo() {
		slog.Error("Current directory is not a git repository")
		return subcommands.ExitFailure
	}

	if a.localName == "" {
		slog.Error("The local flag is required")
		return subcommands.ExitFailure
	}

	// See if the file is already being tracked.
	// If it is, return an error.
	if stat, err := os.Stat(strings.Replace(a.localName, ".", "dot_", 1)); !os.IsNotExist(err) {
		slog.Error("The file is already being tracked", slog.String(loggingKeyFile, a.localName))
		return subcommands.ExitFailure
	} else if err != nil {
		slog.Error("Error checking if the file is already being tracked", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	} else if stat.Name() != "" {
		slog.Error("The file is already being tracked", slog.String(loggingKeyFile, a.localName))
		return subcommands.ExitFailure
	}

	// Find the users home directory path.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Error getting user home directory", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	// Get the absolute path of the file in the home directory.
	homeDotPath := filepath.Join(homeDir, a.localName)

	// Does the file exist?
	// If not, return an error.
	if _, err := os.Stat(homeDotPath); os.IsNotExist(err) {
		slog.Error("The file does not exist", slog.String(loggingKeyFile, homeDotPath))
		return subcommands.ExitFailure
	} else if err != nil {
		slog.Error("Error checking if the file exists", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	// Add the file to the repository.
	if err := addFile(homeDotPath); err != nil {
		slog.Error("Error adding file to the repository", slog.String(loggingKeyError, err.Error()))
		return subcommands.ExitFailure
	}

	slog.Info("File added to the repository", slog.String(loggingKeyFile, homeDotPath))

	return subcommands.ExitSuccess
}
