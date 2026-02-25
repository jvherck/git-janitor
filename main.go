package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
		state:   stateList,
	}
	initialModel.list.Title = "Git Branch Janitor"

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
