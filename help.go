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
This file contains the logic for rendering the command-line help menu. It uses lipgloss for styling the output to
provide a clean and readable experience.
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// commands lists the available positional commands and their descriptions.
var commands = [][]string{
	{"help", "Shows this help menu"},
	{"version", "Shows the version of Git Janitor"},
}

// flags lists the available CLI flags and their descriptions.
var flags = [][]string{
	{"-h, --help", "Shows this help menu"},
	{"-v, --version", "Shows the version of Git Janitor"},
	{"--dry-run", "Simulate deletion without actually removing branches"},
	{"--protect", "Comma-separated list of branches to protect, supports * wildcards (e.g. 'qa,release-*')"},
	{"--stale-days", "Number of days before a branch is considered stale (default 30)"},
}

// keyBindings lists the TUI interactive keybindings and their functions.
var keyBindings = [][]string{
	{"↑ / k", "Move cursor up"},
	{"↓ / j", "Move cursor down"},
	{"Space", "Toggle selection for the current branch"},
	{"a", "Select ALL unprotected branches"},
	{"m", "Select MERGED branches (those already merged into the default branch)"},
	{"g", "Select GONE branches (those whose upstream remote branch was deleted)"},
	{"s", "Select STALE branches (no commits within the stale-days threshold)"},
	{"c", "CLEAR all current selections"},
	{"o", "Cycle through sort orders (Alphabetical, Latest Commits, Oldest Commits)"},
	{"Enter", "Proceed to deletion confirmation screen"},
	{"q / Ctrl+C", "Quit the application"},
}

// printHelp constructs and renders a styled help menu to standard output,
// then exits the application successfully.
func printHelp() {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle)).Bold(true)
	sb.WriteString(titleStyle.Render("🧹 Git Janitor") + "\n\n")

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	sb.WriteString(descStyle.Render("A fast, interactive TUI for cleaning up local Git branches.") + "\n\n")

	headingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary)).Bold(true).Underline(true)
	keyColumnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).Width(18)

	// Usage section
	sb.WriteString(headingStyle.Render("USAGE:") + "\n")
	sb.WriteString("  git-janitor [command] [flags]\n\n")

	// Commands section
	sb.WriteString(headingStyle.Render("COMMANDS:") + "\n")
	for _, row := range commands {
		sb.WriteString(fmt.Sprintf("  %s %s\n", keyColumnStyle.Render(row[0]), row[1]))
	}

	// Flags section
	sb.WriteString("\n" + headingStyle.Render("FLAGS:") + "\n")
	for _, row := range flags {
		sb.WriteString(fmt.Sprintf("  %s %s\n", keyColumnStyle.Render(row[0]), row[1]))
	}

	// Keybindings section
	sb.WriteString("\n" + headingStyle.Render("INTERACTIVE KEYBINDINGS:") + "\n")
	for _, row := range keyBindings {
		sb.WriteString(fmt.Sprintf("  %s %s\n", keyColumnStyle.Render(row[0]), row[1]))
	}

	// Examples section for better clarity
	sb.WriteString("\n" + headingStyle.Render("EXAMPLES:") + "\n")
	sb.WriteString("  git-janitor --dry-run\n")
	sb.WriteString("  git-janitor --protect \"production,release-*\"\n")
	sb.WriteString("  git-janitor --stale-days 14\n")

	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorSecondary)).
		Padding(1, 4).
		Render(sb.String())

	fmt.Println(helpBox)
	os.Exit(0)
}
