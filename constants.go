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

// Color definitions for the application's UI.
// These accept standard ANSI color codes (e.g., "62") or Hex codes (e.g., "#8839ef").
const (
	ColorPrimary   = "#38BDF8" // Used for primary borders and the confirm dialog
	ColorSecondary = "#7DD3FC" // Used for the summary box border
	ColorTextMuted = "#64748B" // Used for footer help text
	ColorSuccess   = "#4ADE80" // Used for success message headers
	ColorWarning   = "#FB923C" // Used for error and warning message headers
	ColorTitle     = "#FBBF24" // Used for the summary title text
)

// Layout and sizing constraints to maintain consistent spacing.
const (
	DocMarginVertical        = 1
	DocMarginHorizontal      = 2
	DialogPaddingVertical    = 1
	ConfirmPaddingHorizontal = 4
	SummaryPaddingHorizontal = 3
	SummaryBoxWidth          = 50
)

// UI text symbols for the interactive list.
const (
	SymbolProtected  = "🔒  "
	SymbolSelected   = "[x] "
	SymbolUnselected = "[ ] "
)

// Standard branch names used for fallback checks and default protections.
const (
	BranchMain   = "main"
	BranchMaster = "master"
	BranchDev    = "dev"
)
