/*
MIT License

Copyright (c) 2026 Jan Van Herck (https://github.com/jvherck)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
This file implements the User Interface (UI) using the Bubble Tea framework. It handles user interactions, rendering
the branch list, and managing the confirmation dialog for branch deletion.
*/
package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// appState defines the distinct UI views within the application's lifecycle.
type appState int

const (
	stateLoading appState = iota // Initial loading view
	stateList                    // Primary view: selecting branches to delete
	stateConfirm                 // Secondary view: confirming the deletion of selected branches
)

var (
	// docStyle defines the basic layout margins for the main application window.
	docStyle = lipgloss.NewStyle().Margin(DocMarginVertical, DocMarginHorizontal)

	// footerStyle defines the appearance of the keybinding hints at the bottom of the list.
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextMuted))

	// confirmStyle defines the appearance of the deletion confirmation dialog.
	confirmStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorPrimary)).
			Padding(DialogPaddingVertical, ConfirmPaddingHorizontal).
			Align(lipgloss.Center)
)

// model encapsulates the complete application state for the Bubble Tea framework.
type model struct {
	list         list.Model    // The core list component for displaying branches
	spinner      spinner.Model // The loading spinner
	deleted      []string      // Names of branches successfully deleted (or marked for deletion in dry-run)
	errs         []string      // Error messages encountered during branch deletion
	state        appState      // The current UI view
	width        int           // Current terminal window width
	height       int           // Current terminal window height
	dryRun       bool          // If true, no actual deletion occurs
	sortMode     SortMode      // Current sorting order of the branch list
	protectFlag  string        // Protection patterns to use when fetching branches
	staleDays    float64       // Stale days threshold
	updateNotice string        // Any notice about available updates
}

type branchesMsg struct {
	items []list.Item
	err   error
}

func fetchBranchesCmd(protectFlag string, staleDays float64) tea.Cmd {
	return func() tea.Msg {
		items, err := getLocalBranches(protectFlag, staleDays)
		return branchesMsg{items: items, err: err}
	}
}

// Init handles background tasks upon application startup.
// It returns a tea.Cmd which is executed when the program starts.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchBranchesCmd(m.protectFlag, m.staleDays),
		checkUpdateCmd(),
	)
}

// Update acts as the central event loop, processing keypresses, window resizes, and other messages.
// It updates the model and returns a tea.Cmd for any side effects.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resizing by updating list dimensions
		m.width = msg.Width
		m.height = msg.Height
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-2)

	case tea.KeyMsg:
		// Logic depends on the current window state
		if m.state == stateLoading {
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
			return m, nil
		}

		if m.state == stateConfirm {
			return m.handleConfirmUpdate(msg)
		}

		if m.state == stateList {
			// Only intercept custom selection keys if the user is not actively typing in the filter input.
			if m.list.FilterState() != list.Filtering {
				var handled bool
				m, cmd, handled = m.handleCustomListKeys(msg)
				if handled {
					return m, cmd
				}
			}
		}

	case branchesMsg:
		if msg.err != nil {
			m.errs = append(m.errs, fmt.Sprintf("Error fetching branches: %v", msg.err))
			return m, tea.Quit
		}
		cmd = m.list.SetItems(msg.items)
		cmds = append(cmds, cmd)
		m.state = stateList

	case checkUpdateMsg:
		m.updateNotice = msg.notice

	case spinner.TickMsg:
		if m.state == stateLoading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Passes unhandled messages down to the list component to ensure standard
	// navigation (up/down/filtering) continues to function.
	if m.state == stateList {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleCustomListKeys intercepts keys specific to branch operations (selection, filtering, sorting)
// before the list component handles them.
// The boolean return value determines if the main event loop should halt propagation.
func (m model) handleCustomListKeys(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true

	case " ":
		// Toggle selection for the currently highlighted branch
		if i, ok := m.list.SelectedItem().(item); ok {
			if !i.isProtected {
				i.selected = !i.selected
				m.list.SetItem(m.list.Index(), i)
			}
		}
		return m, nil, true

	case "a", "m", "g", "s", "c":
		// Handle batch selection/deselection operations.
		// These operations only affect the currently visible items (respecting filters).
		visible := m.list.VisibleItems()
		isVisible := make(map[string]bool)
		for _, v := range visible {
			if i, ok := v.(item); ok {
				isVisible[i.name] = true
			}
		}

		var newItems []list.Item
		for _, listItem := range m.list.Items() {
			i := listItem.(item)

			// Only modify the selection if the item is currently visible in the filtered list
			if isVisible[i.name] {
				if msg.String() == "c" {
					// Clear selection for visible items
					i.selected = false
				} else if !i.isProtected {
					// Select based on specific criteria for visible items
					if msg.String() == "a" ||
						(msg.String() == "m" && i.isMerged) ||
						(msg.String() == "g" && i.isGone) ||
						(msg.String() == "s" && i.isStale) {
						i.selected = true
					}
				}
			}
			newItems = append(newItems, i)
		}
		cmd := m.list.SetItems(newItems)
		return m, cmd, true

	case "o":
		// Cycle through sorting modes: Alphabetical -> Newest First -> Oldest First
		m.sortMode = (m.sortMode + 1) % 3

		items := m.list.Items()
		sort.Slice(items, func(i, j int) bool {
			a := items[i].(item)
			b := items[j].(item)

			switch m.sortMode {
			case SortLatestCommits:
				return a.lastCommitUnix > b.lastCommitUnix
			case SortOldestCommits:
				return a.lastCommitUnix < b.lastCommitUnix
			default:
				return a.name < b.name
			}
		})

		cmd := m.list.SetItems(items)
		return m, cmd, true

	case "enter":
		// Transition to confirmation state if any branches are selected
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
		// Abort deletion and return to the list
		m.state = stateList
		return m, nil

	case "y":
		// Execute deletion for all selected branches
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
	if m.state == stateLoading {
		s := fmt.Sprintf("\n\n   %s Fetching git branches...\n\n", m.spinner.View())
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
	}

	if m.state == stateConfirm {
		// Render the confirmation dialog
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

	// Render the primary branch list with footer hints
	var sortHint string
	switch m.sortMode {
	case SortLatestCommits:
		sortHint = "Latest Commits"
	case SortOldestCommits:
		sortHint = "Oldest Commits"
	default:
		sortHint = "Alphabetical"
	}

	footer := footerStyle.Render(fmt.Sprintf("  a: all • m: merged • g: gone • s: stale • c: clear • o: sort (%s)", sortHint))
	return docStyle.Render(m.list.View() + "\n" + footer)
}
