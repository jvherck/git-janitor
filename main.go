package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// model represents the state of the application.
type model struct {
	cursor int
}

// Init is called when the application starts.
// It can return a command to run background tasks, but returns nil for now.
func (m model) Init() tea.Cmd {
	return nil
}

// Update processes incoming messages and updates the application state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			m.cursor++
		}
	}

	return m, nil
}

// View renders the UI based on the current data in the model.
func (m model) View() string {
	s := "Git Branch Janitor (Skeleton)\n\n"
	s += fmt.Sprintf("Cursor position: %d\n\n", m.cursor)
	s += "Press 'j/k' or 'up/down' to move. Press 'q' to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
