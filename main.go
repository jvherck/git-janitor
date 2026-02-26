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
This file provides the entry point for Git Janitor, an interactive TUI tool designed to help developers clean up their
local Git branches efficiently.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Variables for the ldflags to overwrite during Github Action build.
// These are used to provide versioning information to the user.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// main is the entry point of the application. It parses command-line flags,
// initializes the branch list, and starts the Bubble Tea TUI.
func main() {
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "Shows this help menu")
	flag.BoolVar(&showHelp, "help", false, "Shows this help menu")

	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Shows the version of Git Janitor")
	flag.BoolVar(&showVersion, "version", false, "Shows the version of Git Janitor")

	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "Simulate deletion without removing branches")

	var protectFlag string
	flag.StringVar(&protectFlag, "protect", "", "Comma-separated list of branches to protect (supports wildcards like 'release-*')")

	var staleDays float64
	flag.Float64Var(&staleDays, "stale-days", 30, "Number of days before a branch is considered stale")

	// Override default usage to show our custom help menu
	flag.Usage = printHelp
	flag.Parse()

	// Check for positional commands like 'help' or 'version'
	if flag.NArg() > 0 {
		cmd := flag.Arg(0)
		if cmd == "help" {
			showHelp = true
		} else if cmd == "version" {
			showVersion = true
		}
	}

	if showHelp {
		printHelp()
	}
	if showVersion {
		printVersion()
	}

	// Fetch local branches based on protection and staleness criteria
	items, err := getLocalBranches(protectFlag, staleDays)
	if err != nil {
		fmt.Printf("Error fetching branches: %v\nAre you in a git repository?\n", err)
		os.Exit(1)
	}

	// Configure the list delegate for consistent styling
	delegate := list.NewDefaultDelegate()

	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color(ColorPrimary)).
		BorderLeftForeground(lipgloss.Color(ColorPrimary))

	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color(ColorSecondary)).
		BorderLeftForeground(lipgloss.Color(ColorPrimary))

	// Initialize the list component
	l := list.New(items, delegate, 0, 0)
	l.Title = "Git Janitor"
	if dryRun {
		l.Title += " (DRY RUN)"
	}

	// Apply custom styling to the list
	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorPrimary)).
		Foreground(lipgloss.Color("#0F172A")).
		Padding(0, 1)
	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle))
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary))

	// Set up the initial application model
	initialModel := model{
		list:    l,
		deleted: []string{},
		errs:    []string{},
		state:   stateList,
		dryRun:  dryRun,
	}

	// Start the Bubble Tea program with the alternate screen buffer
	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	// Print a summary of actions taken after the TUI exits
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
		headerText := "Deleted Branches:"
		if m.dryRun {
			headerText = "Dry Run - Would Delete (but not actually):"
		}

		successHeader := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).Bold(true).Render(headerText)
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
