package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// docStyle defines the overall application margin.
var docStyle = lipgloss.NewStyle().Margin(1, 2)

// item represents a single row in the list.
type item struct {
	name     string
	selected bool
}

// Title returns the primary text for the list item.
func (i item) Title() string {
	if i.selected {
		return "[x] " + i.name
	}
	return "[ ] " + i.name
}

// Description provides secondary text below the title.
func (i item) Description() string { return "Press space to toggle selection" }

// FilterValue determines what text the list's fuzzy finder searches against.
func (i item) FilterValue() string { return i.name }

// model represents the application state.
type model struct {
	list    list.Model
	deleted []string
	errs    []string
}

// getLocalBranches executes the git CLI to retrieve local branches and formats them as list items.
func getLocalBranches() ([]list.Item, error) {
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rawOutput := strings.TrimSpace(string(out))
	if rawOutput == "" {
		return []list.Item{}, nil
	}

	branchNames := strings.Split(rawOutput, "\n")
	items := make([]list.Item, len(branchNames))
	for i, name := range branchNames {
		items[i] = item{name: name, selected: false}
	}

	return items, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case " ":
			if i, ok := m.list.SelectedItem().(item); ok {
				i.selected = !i.selected
				m.list.SetItem(m.list.Index(), i)
			}

		case "enter":
			for _, listItem := range m.list.Items() {
				i, ok := listItem.(item)
				if ok && i.selected {
					// Prevent accidental deletion of primary branches.
					if i.name == "main" || i.name == "master" {
						m.errs = append(m.errs, fmt.Sprintf("Skipped %s (protected)", i.name))
						continue
					}

					// Execute git branch -D <branch_name> to force delete the branch.
					cmd := exec.Command("git", "branch", "-D", i.name)
					if err := cmd.Run(); err != nil {
						m.errs = append(m.errs, fmt.Sprintf("Failed to delete %s", i.name))
					} else {
						m.deleted = append(m.deleted, i.name)
					}
				}
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	items, err := getLocalBranches()
	if err != nil {
		fmt.Printf("Error fetching branches: %v\nAre you in a git repository?\n", err)
		os.Exit(1)
	}

	initialModel := model{
		list:    list.New(items, list.NewDefaultDelegate(), 0, 0),
		deleted: []string{},
		errs:    []string{},
	}
	initialModel.list.Title = "Git Branch Janitor"

	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	// Capture the final state of the model after the program quits.
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	// Assert the final model back to our custom model type to access its fields.
	if m, ok := finalModel.(model); ok {
		if len(m.deleted) > 0 {
			fmt.Println("Successfully deleted branches:")
			for _, b := range m.deleted {
				fmt.Printf("  - %s\n", b)
			}
		}
		if len(m.errs) > 0 {
			fmt.Println("\nWarnings/Errors:")
			for _, e := range m.errs {
				fmt.Printf("  - %s\n", e)
			}
		}
	}
}
