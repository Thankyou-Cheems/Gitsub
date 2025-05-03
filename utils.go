package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

func ValidateRepoURL(url string) error {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "git@") {
		return nil
	}
	return fmt.Errorf("invalid repository URL: %s", url)
}

func NormalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "\\", "/")
	return path
}

// ParseGitHubDirURL parses a GitHub tree/blob URL into repo URL, directories, and branch.
// It only accepts a single argument and returns ok=false for non-matching inputs.
func ParseGitHubDirURL(args []string) (string, []string, string, bool) {
	if len(args) != 1 {
		return "", nil, "", false
	}

	raw := strings.TrimSpace(args[0])
	if raw == "" {
		return "", nil, "", false
	}

	u, err := url.Parse(raw)
	if err != nil || u.Host != "github.com" {
		return "", nil, "", false
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 5 {
		return "", nil, "", false
	}

	owner := parts[0]
	repo := strings.TrimSuffix(parts[1], ".git")
	mode := parts[2]
	branch := parts[3]
	subPath := path.Clean(strings.Join(parts[4:], "/"))

	if owner == "" || repo == "" || branch == "" || subPath == "." {
		return "", nil, "", false
	}

	if mode != "tree" && mode != "blob" {
		return "", nil, "", false
	}

	// If a file URL is provided, sparse checkout can still use its parent directory.
	if mode == "blob" {
		subPath = path.Dir(subPath)
		if subPath == "." {
			return "", nil, "", false
		}
	}

	repoURL := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
	return repoURL, []string{subPath}, branch, true
}
