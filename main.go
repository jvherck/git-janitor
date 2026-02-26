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

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Variables for the ldflags to overwrite during Github Action build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Intercept help flags immediately before executing Git checks.
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "-h" || arg == "--help" || arg == "help" {
				printHelp()
			}
			if arg == "-v" || arg == "--version" || arg == "version" {
				printVersion()
			}
		}
	}

	items, err := getLocalBranches()
	if err != nil {
		fmt.Printf("Error fetching branches: %v\nAre you in a git repository?\n", err)
		os.Exit(1)
	}

	// Initialize a custom delegate to override the default styling.
	delegate := list.NewDefaultDelegate()

	// Override the selected item colors.
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color(ColorPrimary)).
		BorderLeftForeground(lipgloss.Color(ColorPrimary))

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color(ColorSecondary)).
		BorderLeftForeground(lipgloss.Color(ColorPrimary))

	// Instantiate the list with the custom delegate.
	l := list.New(items, delegate, 0, 0)
	l.Title = "Git Janitor"

	// Override the default list title and filter styling.
	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorPrimary)).
		Foreground(lipgloss.Color("#0F172A")). // Dark slate for high contrast readability
		Padding(0, 1)
	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle))
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))

	initialModel := model{
		list:    l,
		deleted: []string{},
		errs:    []string{},
		state:   stateList,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok {
		printSummary(m)
	}
}

// printSummary processes the final model state and outputs a formatted
// results table to standard out after the TUI has closed.
func printSummary(m model) {
	if len(m.deleted) == 0 && len(m.errs) == 0 {
		return
	}

	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle)).Bold(true)
	sb.WriteString(titleStyle.Render("* Git Janitor Summary *") + "\n\n")

	if len(m.deleted) > 0 {
		successHeader := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).Bold(true).Render("Deleted Branches:")
		sb.WriteString(successHeader + "\n")
		for _, b := range m.deleted {
			sb.WriteString(fmt.Sprintf("  - %s\n", b))
		}
		sb.WriteString("\n")
	}

	if len(m.errs) > 0 {
		errHeader := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorWarning)).Bold(true).Render("Errors & Warnings:")
		sb.WriteString(errHeader + "\n")
		for _, e := range m.errs {
			sb.WriteString(fmt.Sprintf("  ! %s\n", e))
		}
	}

	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorSecondary)).
		Padding(DialogPaddingVertical, SummaryPaddingHorizontal).
		Width(SummaryBoxWidth).
		Render(strings.TrimSpace(sb.String()))

	fmt.Println(summaryBox)
}
