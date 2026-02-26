package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// appState defines the distinct UI views within the application's lifecycle.
type appState int

const (
	stateList appState = iota
	stateConfirm
)

// Global UI styling definitions utilizing configurable constants.
var (
	docStyle     = lipgloss.NewStyle().Margin(DocMarginVertical, DocMarginHorizontal)
	footerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextMuted))
	confirmStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPrimary)).
			Padding(DialogPaddingVertical, ConfirmPaddingHorizontal).
			Align(lipgloss.Center)
)

// model encapsulates the complete application state.
type model struct {
	list    list.Model
	deleted []string
	errs    []string
	state   appState
	width   int
	height  int
	dryRun  bool
}

// Init handles background tasks upon application startup.
func (m model) Init() tea.Cmd {
	return nil
}

// Update acts as the central event loop, processing keypresses and window resizes.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-2)

	case tea.KeyMsg:
		if m.state == stateConfirm {
			return m.handleConfirmUpdate(msg)
		}

		if m.state == stateList {
			var handled bool
			m, cmd, handled = m.handleCustomListKeys(msg)
			if handled {
				return m, cmd
			}
		}
	}

	// Passes unhandled messages down to the list component to ensure standard
	// navigation (up/down/filtering) continues to function.
	if m.state == stateList {
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleCustomListKeys intercepts keys specific to branch operations before the list handles them.
// The boolean return value determines if the main event loop should halt propagation.
func (m model) handleCustomListKeys(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true

	case " ":
		if i, ok := m.list.SelectedItem().(item); ok {
			if !i.isProtected {
				i.selected = !i.selected
				m.list.SetItem(m.list.Index(), i)
			}
		}
		return m, nil, true

	case "a", "m", "c":
		var newItems []list.Item
		for _, listItem := range m.list.Items() {
			i := listItem.(item)

			if msg.String() == "c" {
				i.selected = false
			} else if !i.isProtected && (msg.String() == "a" || (msg.String() == "m" && i.isMerged)) {
				i.selected = true
			}
			newItems = append(newItems, i)
		}
		cmd := m.list.SetItems(newItems)
		return m, cmd, true

	case "enter":
		hasSelection := false
		for _, listItem := range m.list.Items() {
			if i, ok := listItem.(item); ok && i.selected {
				hasSelection = true
				break
			}
		}
		if hasSelection {
			m.state = stateConfirm
		}
		return m, nil, true
	}

	return m, nil, false
}

// handleConfirmUpdate processes keypresses specifically for the deletion confirmation dialog.
func (m model) handleConfirmUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch strings.ToLower(msg.String()) {
	case "ctrl+c", "q", "n", "esc", "enter":
		m.state = stateList
		return m, nil

	case "y":
		for _, listItem := range m.list.Items() {
			i, ok := listItem.(item)
			if ok && i.selected && !i.isProtected {
				if m.dryRun {
					m.deleted = append(m.deleted, i.name)
				} else {
					cmd := exec.Command("git", "branch", "-D", i.name)
					if err := cmd.Run(); err != nil {
						m.errs = append(m.errs, fmt.Sprintf("Failed to delete %s", i.name))
					} else {
						m.deleted = append(m.deleted, i.name)
					}
				}
			}
		}
		return m, tea.Quit
	}
	return m, nil
}

// View evaluates the current application state and renders the corresponding UI layout.
func (m model) View() string {
	if m.state == stateConfirm {
		selectedCount := 0
		for _, listItem := range m.list.Items() {
			if i, ok := listItem.(item); ok && i.selected {
				selectedCount++
			}
		}

		prompt := fmt.Sprintf("Are you sure you want to force delete %d branches?\n\n(y/N)", selectedCount)
		if m.dryRun {
			prompt = fmt.Sprintf("DRY RUN: Would delete %d branches.\nProceed?\n\n(y/N)", selectedCount)
		}

		confirmBox := confirmStyle.Render(prompt)

		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, confirmBox)
	}

	footer := footerStyle.Render("  a: all • m: merged • c: clear")
	return docStyle.Render(m.list.View() + "\n" + footer)
}
