package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/jacobbrewer1/dotmanager/cmd/add"
	"github.com/jacobbrewer1/dotmanager/cmd/diff"
	"github.com/jacobbrewer1/dotmanager/cmd/pull"
	"github.com/jacobbrewer1/dotmanager/cmd/push"
	"github.com/jacobbrewer1/dotmanager/pkg/utils"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a new dotfile to be tracked",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return add.Files(
						ctx,
					)
				},
			},
			{
				Name:  "diff",
				Usage: "Diff your local dotfiles with the ones in the repository",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return diff.PrintDiff(
						ctx,
					)
				},
			},
			{
				Name:  "push",
				Usage: "Push changes from your local dotfiles to the repository",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return push.Files(
						ctx,
					)
				},
			},
			{
				Name:  "pull",
				Usage: "Pull changes from the repository to your local dotfiles",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return pull.Files(
						ctx,
					)
				},
			},
		},
	}

	ctx, cancel := utils.CoreContext()
	defer cancel()

	if err := cmd.Run(ctx, os.Args); err != nil {
		cancel()
		log.Fatal(err) // nolint:gocritic // Calling cancel() before exiting on the line above
	}
}
