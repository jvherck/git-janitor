package main

// Color definitions for the application's UI.
// These accept standard ANSI color codes (e.g., "62") or Hex codes (e.g., "#8839ef").
const (
	ColorPrimary   = "62"  // Used for primary borders and the confirm dialog
	ColorSecondary = "63"  // Used for the summary box border
	ColorTextMuted = "241" // Used for footer help text
	ColorSuccess   = "42"  // Used for success message headers
	ColorWarning   = "204" // Used for error and warning message headers
	ColorTitle     = "205" // Used for the summary title text
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
