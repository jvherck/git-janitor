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
This file provides the git-related functionality for identifying and filtering branches based on various criteria.
*/
package main

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	// Try to get the default branch from origin/HEAD
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "origin/HEAD")
	out, err := cmd.Output()
	if err == nil {
		ref := strings.TrimSpace(string(out))
		parts := strings.Split(ref, "/")
		if len(parts) == 2 {
			return parts[1]
		}
	}

	// Fallback: check if 'main' exists locally
	if err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+BranchMain).Run(); err == nil {
		return BranchMain
	}

	// Final fallback: assume 'master'
	return BranchMaster
}

// getMergedBranches returns a map acting as a set of branch names that have
// already been merged into the specified target branch.
func getMergedBranches(targetBranch string) map[string]struct{} {
	merged := make(map[string]struct{})

	// List branches that are merged into the target branch
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

// isProtectedBranch evaluates a branch name against default safe branches
// (dev, current, default) and user-provided protection patterns (e.g., release-*).
func isProtectedBranch(name, currentBranch, defaultBranch, customPatterns string) bool {
	// Always protect dev, the current branch, and the default branch
	if name == BranchDev || name == currentBranch || name == defaultBranch {
		return true
	}

	// Check against user-defined patterns
	if customPatterns != "" {
		patterns := strings.Split(customPatterns, ",")
		for _, p := range patterns {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			// Supports shell-style wildcards via filepath.Match
			matched, err := filepath.Match(p, name)
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}

// getLocalBranches interacts with the Git CLI to retrieve all local branches
// and evaluates their detailed status parameters like age, upstream status, and merge state.
func getLocalBranches(protectFlag string, staleDays float64) ([]list.Item, error) {
	currentBranch := getCurrentBranch()
	defaultBranch := getDefaultBranch()
	mergedBranches := getMergedBranches(defaultBranch)

	// Prune remote-tracking branches to ensure upstream status is accurate
	// This will mark local branches without an upstream as "gone"
	err := exec.Command("git", "fetch", "--prune").Run()
	if err != nil {
		// Ignore error as user probably is not connected to the internet, not the end of the world...
	}

	// Format used to extract branch details efficiently in one call
	format := "%(refname:short)|||%(upstream:track)|||%(committerdate:relative)|||%(committerdate:unix)"
	cmd := exec.Command("git", "for-each-ref", "--format="+format, "refs/heads/")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rawOutput := strings.ReplaceAll(string(out), "\r\n", "\n")
	rawOutput = strings.TrimSpace(rawOutput)

	if rawOutput == "" {
		return []list.Item{}, nil
	}

	lines := strings.Split(rawOutput, "\n")
	var items []list.Item
	now := time.Now().Unix()
	staleSeconds := int64(staleDays * 24 * 60 * 60)

	for _, line := range lines {
		parts := strings.Split(line, "|||")
		if len(parts) != 4 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		upstreamStatus := parts[1] // Contains info if upstream is [gone]
		relativeDate := parts[2]   // Human-readable age (e.g., "2 weeks ago")
		unixStr := strings.TrimSpace(parts[3])

		if name == "" {
			continue
		}

		isProtected := isProtectedBranch(name, currentBranch, defaultBranch, protectFlag)
		_, isMerged := mergedBranches[name]
		isGone := strings.Contains(upstreamStatus, "[gone]")

		// Evaluate staleness based on the provided threshold
		isStale := false
		var unixTime int64
		if u, err := strconv.ParseInt(unixStr, 10, 64); err == nil {
			unixTime = u
			if now-unixTime > staleSeconds {
				isStale = true
			}
		}

		items = append(items, item{
			name:           name,
			selected:       false,
			isProtected:    isProtected,
			isMerged:       isMerged,
			isGone:         isGone,
			isStale:        isStale,
			age:            relativeDate,
			lastCommitUnix: unixTime,
		})
	}

	return items, nil
}
