package main

import (
	"testing"
)

func TestIsProtectedBranch(t *testing.T) {
	tests := []struct {
		name           string
		branchName     string
		currentBranch  string
		defaultBranch  string
		customPatterns string
		expected       bool
	}{
		{
			name:           "protect dev branch by default",
			branchName:     "dev",
			currentBranch:  "feature",
			defaultBranch:  "main",
			customPatterns: "",
			expected:       true,
		},
		{
			name:           "protect current branch",
			branchName:     "feature",
			currentBranch:  "feature",
			defaultBranch:  "main",
			customPatterns: "",
			expected:       true,
		},
		{
			name:           "protect default branch",
			branchName:     "main",
			currentBranch:  "feature",
			defaultBranch:  "main",
			customPatterns: "",
			expected:       true,
		},
		{
			name:           "protect exact custom pattern",
			branchName:     "qa",
			currentBranch:  "feature",
			defaultBranch:  "main",
			customPatterns: "qa",
			expected:       true,
		},
		{
			name:           "protect wildcard custom pattern",
			branchName:     "release-1.0",
			currentBranch:  "feature",
			defaultBranch:  "main",
			customPatterns: "release-*",
			expected:       true,
		},
		{
			name:           "do not protect other branches",
			branchName:     "fix-bug",
			currentBranch:  "feature",
			defaultBranch:  "main",
			customPatterns: "release-*",
			expected:       false,
		},
		{
			name:           "multiple custom patterns",
			branchName:     "staging",
			currentBranch:  "feature",
			defaultBranch:  "main",
			customPatterns: "qa,staging,release-*",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isProtectedBranch(tt.branchName, tt.currentBranch, tt.defaultBranch, tt.customPatterns)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for branch %s", tt.expected, result, tt.branchName)
			}
		})
	}
}
