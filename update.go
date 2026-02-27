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
This file provides auto-update functionality for Git Janitor. It detects how the tool was installed,
checks for newer releases on GitHub, and can upgrade the tool using the appropriate package manager.
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	githubRepo         = "jvherck/git-janitor"
	githubReleasesAPI  = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
	updateCheckTimeout = 3 * time.Second
)

// packageManager represents a detected package manager and the command needed to upgrade git-janitor.
type packageManager struct {
	name    string
	command []string
}

// githubRelease is a minimal representation of the GitHub releases API response.
type githubRelease struct {
	TagName string `json:"tag_name"`
}

// detectPackageManager tries to figure out how git-janitor was installed by checking
// which package managers are available and whether they manage this binary.
func detectPackageManager() *packageManager {
	exe, err := os.Executable()
	if err != nil {
		exe = ""
	}
	exeLower := strings.ToLower(exe)

	switch runtime.GOOS {
	case "darwin", "linux":
		// Homebrew — check if the executable lives under the Homebrew prefix
		if brewPrefix, err := exec.LookPath("brew"); err == nil {
			brewPrefixDir := strings.TrimSuffix(brewPrefix, "/bin/brew")
			if strings.HasPrefix(exeLower, strings.ToLower(brewPrefixDir)) {
				return &packageManager{name: "Homebrew", command: []string{"brew", "upgrade", "jvherck/tap/git-janitor"}}
			}
			// Also ask brew where it thinks git-janitor lives
			out, err2 := exec.Command("brew", "--prefix", "git-janitor").Output()
			if err2 == nil && strings.TrimSpace(string(out)) != "" {
				return &packageManager{name: "Homebrew", command: []string{"brew", "upgrade", "jvherck/tap/git-janitor"}}
			}
		}

	case "windows":
		// Scoop — binary typically lives under %USERPROFILE%\scoop
		if strings.Contains(exeLower, "scoop") {
			return &packageManager{name: "Scoop", command: []string{"scoop", "update", "git-janitor"}}
		}
	}

	// Fallback: go install
	if _, err := exec.LookPath("go"); err == nil {
		return &packageManager{name: "go install", command: []string{"go", "install", "github.com/" + githubRepo + "@latest"}}
	}

	return nil
}

// fetchLatestVersion queries the GitHub releases API and returns the latest tag name.
// Returns an empty string and an error on failure.
func fetchLatestVersion() (string, error) {
	client := &http.Client{Timeout: updateCheckTimeout}
	req, err := http.NewRequest("GET", githubReleasesAPI, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "git-janitor/"+version)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	return strings.TrimPrefix(release.TagName, "v"), nil
}

// normalizeVersion strips a leading "v" so we can compare "v1.2.3" with "1.2.3".
func normalizeVersion(v string) string {
	return strings.TrimPrefix(v, "v")
}

// checkForUpdateNotification fetches the latest version and, if a newer version exists,
// prints a styled one-line notice to stdout. It is intentionally fire-and-forget so it
// never blocks the normal TUI startup.
func checkForUpdateNotification() {
	// Skip the check entirely for dev builds that have no real version.
	if version == "dev" || version == "" {
		return
	}

	latest, err := fetchLatestVersion()
	if err != nil || latest == "" {
		return
	}

	if normalizeVersion(latest) != normalizeVersion(version) {
		infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle)).Bold(true)
		mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextMuted))
		fmt.Println(
			infoStyle.Render("  ✦ Update available!") +
				mutedStyle.Render(fmt.Sprintf(" %s → %s  ", normalizeVersion(version), latest)) +
				mutedStyle.Render("Run: git-janitor update"),
		)
		fmt.Println()
	}
}

// runUpdate detects the package manager and runs the appropriate upgrade command.
// It streams the child process output directly to the user's terminal.
func runUpdate() {
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTitle)).Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTextMuted))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ColorWarning)).Bold(true)

	fmt.Println(titleStyle.Render("🧹  Git Janitor - Update"))
	fmt.Println()

	// Check what's available upstream first
	fmt.Print(mutedStyle.Render("  Checking for updates..."))
	latest, err := fetchLatestVersion()
	if err != nil {
		fmt.Println()
		fmt.Println(warnStyle.Render("  ✗ Could not reach GitHub: ") + err.Error())
		os.Exit(1)
	}

	current := normalizeVersion(version)
	latestClean := normalizeVersion(latest)

	fmt.Println() // newline after "Checking..."

	if current == latestClean && version != "dev" {
		fmt.Println(successStyle.Render("  ✓ You are already on the latest version (" + current + ")"))
		os.Exit(0)
	}

	if version != "dev" {
		fmt.Println(mutedStyle.Render(fmt.Sprintf("  Current version: %s", current)))
		fmt.Println(mutedStyle.Render(fmt.Sprintf("  Latest version:  %s", latestClean)))
		fmt.Println()
	}

	pm := detectPackageManager()
	if pm == nil {
		fmt.Println(warnStyle.Render("  ✗ Could not detect your package manager."))
		fmt.Println(mutedStyle.Render("  Please upgrade manually:"))
		fmt.Println(mutedStyle.Render("    go install github.com/" + githubRepo + "@latest"))
		os.Exit(1)
	}

	fmt.Println(mutedStyle.Render(fmt.Sprintf("  Detected package manager: %s", pm.name)))
	fmt.Println(mutedStyle.Render(fmt.Sprintf("  Running: %s", strings.Join(pm.command, " "))))
	fmt.Println()

	cmd := exec.Command(pm.command[0], pm.command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println()
		fmt.Println(warnStyle.Render("  ✗ Update failed: ") + err.Error())
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(successStyle.Render("  ✓ Git Janitor updated successfully!"))
	os.Exit(0)
}
