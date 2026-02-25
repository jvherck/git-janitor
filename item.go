package main

// item represents a single selectable branch row within the UI list.
// It implements the bubbles/list.Item interface.
type item struct {
	name        string
	selected    bool
	isProtected bool
	isMerged    bool
}

// Title returns the primary formatted text for the list item.
func (i item) Title() string {
	if i.isProtected {
		return SymbolProtected + i.name
	}
	if i.selected {
		return SymbolSelected + i.name
	}
	return SymbolUnselected + i.name
}

// Description provides contextual secondary text rendered below the title.
func (i item) Description() string {
	if i.isProtected {
		return "Protected or active branch (cannot be deleted)"
	}

	desc := "Press space to toggle selection"
	if i.isMerged {
		desc += " | Merged"
	}
	return desc
}

// FilterValue determines the string that the list's fuzzy finder evaluates against.
func (i item) FilterValue() string {
	return i.name
}
