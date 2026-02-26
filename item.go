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
