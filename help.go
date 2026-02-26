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

	"github.com/charmbracelet/lipgloss"
)

// printHelp constructs and renders a styled help menu to standard output,
// then exits the application successfully.
func printHelp() {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle)).Bold(true)
	sb.WriteString(titleStyle.Render("Git Janitor") + "\n\n")

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	sb.WriteString(descStyle.Render("A fast, interactive TUI for cleaning up local Git branches.") + "\n\n")

	headingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPrimary)).Bold(true).Underline(true)

	sb.WriteString(headingStyle.Render("USAGE:") + "\n")
	sb.WriteString("  git-janitor [flags]\n\n")

	sb.WriteString(headingStyle.Render("KEYBINDINGS:") + "\n")

	keys := [][]string{
		{"↑ / k", "Move cursor up"},
		{"↓ / j", "Move cursor down"},
		{"Space", "Toggle selection for the current branch"},
		{"a", "Select ALL unprotected branches"},
		{"m", "Select only MERGED branches"},
		{"c", "CLEAR all selections"},
		{"Enter", "Proceed to deletion confirmation"},
		{"q / Ctrl+C", "Quit"},
	}

	keyColumnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).Width(18)

	for _, row := range keys {
		sb.WriteString(fmt.Sprintf("  %s %s\n", keyColumnStyle.Render(row[0]), row[1]))
	}
	sb.WriteString("\n")

	sb.WriteString(headingStyle.Render("FLAGS:") + "\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", keyColumnStyle.Render("-h, --help"), "Show this help menu"))

	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorSecondary)).
		Padding(1, 4).
		Render(sb.String())

	fmt.Println(helpBox)
	os.Exit(0)
}
