package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// printVersion outputs the application's version, commit hash, and build date,
// then exits the application successfully.
func printVersion() {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle)).Bold(true)
	fmt.Println(titleStyle.Render("🧹  Git Janitor") + fmt.Sprintf(" version %s", version))

	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextMuted))
	fmt.Println(mutedStyle.Render(fmt.Sprintf("  Commit:  %s", commit)))
	fmt.Println(mutedStyle.Render(fmt.Sprintf("  Built:   %s", date)))
	os.Exit(0)
}
