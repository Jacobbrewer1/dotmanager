package main

import (
	"github.com/charmbracelet/huh"
)

func getFileFromUser(availableFiles []string) ([]string, error) {
	const title = "Select files to start tracking"

	selected := make([]string, 0)

	fileOpts := make([]huh.Option[string], 0)
	for _, f := range availableFiles {
		fileOpts = append(fileOpts, huh.NewOption(f, f))
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(title).
				Options(fileOpts...).
				Value(&selected),
		),
	).WithShowHelp(true).Run()
	if err != nil {
		return nil, err
	}

	return selected, nil
}
