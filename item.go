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
This file defines the 'item' struct, which implements the list.Item interface from the charmbracelet/bubbles/list
package. It represents a single branch entry in the TUI list.
*/
package main

import "strings"

// item represents a single selectable branch row within the UI list.
// It tracks the branch name, user selection state, and various metadata
// flags used for filtering and displaying status information.
type item struct {
	name           string // The Git branch name
	selected       bool   // Whether the user has marked this branch for deletion
	isProtected    bool   // Whether the branch is protected (cannot be deleted)
	isMerged       bool   // Whether the branch has been merged into the default branch
	isGone         bool   // Whether the upstream tracking branch no longer exists
	isStale        bool   // Whether the branch exceeds the stale days threshold
	age            string // Human-readable relative time of the last commit
	lastCommitUnix int64  // Unix timestamp of the last commit (used for sorting)
}

// Title returns the primary formatted text for the list item.
// It includes a status symbol indicating its state (e.g.: protected, selected, unselected).
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
// It lists status tags (like Merged, Gone, Stale) and the last active date.
func (i item) Description() string {
	if i.isProtected {
		d := "Protected or active branch (cannot be deleted)"
		if i.age != "" {
			d += " | Last active: " + i.age
		}
		return d
	}

	var tags []string
	if i.isMerged {
		tags = append(tags, "Merged")
	}
	if i.isGone {
		tags = append(tags, "Gone (Upstream deleted)")
	}
	if i.isStale {
		tags = append(tags, "Stale")
	}

	desc := "Press space to toggle"
	if len(tags) > 0 {
		desc += " | " + strings.Join(tags, ", ")
	}

	if i.age != "" {
		desc += " | Last active: " + i.age
	}

	return desc
}

// FilterValue determines the string that the list's fuzzy finder evaluates against.
// In this case, we filter by the branch name.
func (i item) FilterValue() string {
	return i.name
}
