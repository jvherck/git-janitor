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
This file defines the core constants used throughout the Git Janitor application, including UI colors, layout
dimensions, interactive symbols, and default branch names.
*/
package main

// SortMode defines the available sorting orders for the branch list.
type SortMode int

const (
	SortAlphabetical SortMode = iota
	SortLatestCommits
	SortOldestCommits
)

// UI color definitions for the application's interface.
// These accept standard ANSI color codes (e.g., "62") or Hex codes (e.g., "#8839ef").
const (
	ColorPrimary   = "#38BDF8" // Used for primary borders, selected item titles, and the confirm dialog
	ColorSecondary = "#7DD3FC" // Used for selected item descriptions and the summary box border
	ColorTextMuted = "#64748B" // Used for footer help text and version details
	ColorSuccess   = "#4ADE80" // Used for success message headers and help key columns
	ColorWarning   = "#FB923C" // Used for error and warning message headers
	ColorTitle     = "#FBBF24" // Used for the summary and help menu title text
)

// Layout and sizing constraints to maintain consistent spacing across different views.
const (
	DocMarginVertical        = 1  // Vertical margin around the main list
	DocMarginHorizontal      = 2  // Horizontal margin around the main list
	DialogPaddingVertical    = 1  // Vertical padding inside dialog boxes
	ConfirmPaddingHorizontal = 4  // Horizontal padding inside the confirmation dialog
	SummaryPaddingHorizontal = 3  // Horizontal padding inside the final summary box
	SummaryBoxWidth          = 50 // Fixed width for the final summary box
)

// UI text symbols for the interactive branch list.
const (
	SymbolProtected  = "🔒  "  // Prefix for branches that cannot be deleted
	SymbolSelected   = "[x] " // Prefix for branches marked for deletion
	SymbolUnselected = "[ ] " // Prefix for branches not currently selected
)

// Standard branch names used for fallback checks and default protections.
const (
	BranchMain   = "main"   // Standard primary branch name
	BranchMaster = "master" // Legacy/Standard primary branch name
	BranchDev    = "dev"    // Standard development branch name
)
