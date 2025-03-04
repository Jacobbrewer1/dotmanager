package main

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type fileSelector struct {
	title   string
	choices []string

	cursor int
	choice string
}

func (m *fileSelector) Init() tea.Cmd {
	return nil
}

func (m *fileSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			return m, tea.Quit
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}
	return m, nil
}

func (m *fileSelector) View() string {
	s := strings.Builder{}
	s.WriteString(m.title + "\n\n")

	for i := range len(m.choices) {
		str := "( ) "
		if m.cursor == i {
			str = "(•) "
		}
		s.WriteString(str)
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}

	s.WriteString("\n(press q to quit)\n")
	return s.String()
}

func (m *fileSelector) Exec() (string, error) {
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return "", fmt.Errorf("running program: %w", err)
	}
	return m.choice, nil
}

func newFileSelector(title string, files []string) fileSelector {
	return fileSelector{
		title:   title,
		choices: files,
	}
}
