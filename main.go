package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(new(versionCmd), "")
	subcommands.Register(new(diffCmd), "")
	subcommands.Register(new(pullCmd), "")
	subcommands.Register(new(pushCmd), "")
	subcommands.Register(new(addCmd), "")

	flag.Parse()

	// Listen for ctrl+c and kill signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		got := <-sig
		slog.Info("Received signal, shutting down", slog.String(loggingKeySignal, got.String()))
		cancel()
	}()

	os.Exit(int(subcommands.Execute(ctx))) // nolint:gocritic // This is the main entry point for the application
}
