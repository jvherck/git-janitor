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

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

// getCurrentBranch returns the name of the currently checked-out local branch.
func getCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// getDefaultBranch attempts to identify the primary branch of the repository.
// It checks the remote origin's HEAD reference first, falling back to local standard
// branches defined in constants, and finally defaulting to the standard master branch.
func getDefaultBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "origin/HEAD")
	out, err := cmd.Output()
	if err == nil {
		ref := strings.TrimSpace(string(out))
		parts := strings.Split(ref, "/")
		if len(parts) == 2 {
			return parts[1]
		}
	}

	if err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+BranchMain).Run(); err == nil {
		return BranchMain
	}

	return BranchMaster
}

// getMergedBranches returns a map acting as a set of branch names that have
// already been merged into the specified target branch.
func getMergedBranches(targetBranch string) map[string]struct{} {
	merged := make(map[string]struct{})

	cmd := exec.Command("git", "branch", "--format=%(refname:short)", "--merged", targetBranch)
	out, err := cmd.Output()
	if err != nil {
		return merged
	}

	rawOutput := strings.ReplaceAll(string(out), "\r\n", "\n")
	for _, name := range strings.Split(rawOutput, "\n") {
		name = strings.TrimSpace(name)
		if name != "" {
			merged[name] = struct{}{}
		}
	}

	return merged
}

// getLocalBranches interacts with the Git CLI to retrieve all local branches,
// evaluates their protected and merged status, and formats them as UI list items.
func getLocalBranches() ([]list.Item, error) {
	currentBranch := getCurrentBranch()
	defaultBranch := getDefaultBranch()
	mergedBranches := getMergedBranches(defaultBranch)

	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rawOutput := strings.ReplaceAll(string(out), "\r\n", "\n")
	rawOutput = strings.TrimSpace(rawOutput)

	if rawOutput == "" {
		return []list.Item{}, nil
	}

	branchNames := strings.Split(rawOutput, "\n")
	var items []list.Item

	for _, name := range branchNames {
		name = strings.TrimSpace(name)

		isProtected := name == BranchDev || name == currentBranch || name == defaultBranch
		_, isMerged := mergedBranches[name]

		items = append(items, item{
			name:        name,
			selected:    false,
			isProtected: isProtected,
			isMerged:    isMerged,
		})
	}

	return items, nil
}
