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
