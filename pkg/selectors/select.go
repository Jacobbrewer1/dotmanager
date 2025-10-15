package selectors

import "github.com/charmbracelet/huh"

func UserSelectionForm(title string, options []string) ([]string, error) {
	selected := make([]string, 0)

	opts := make([]huh.Option[string], 0)
	for _, f := range options {
		opts = append(opts, huh.NewOption(f, f))
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(title).
				Options(opts...).
				Value(&selected),
		),
	).WithShowHelp(true).Run()
	if err != nil {
		return nil, err
	}

	return selected, nil
}
