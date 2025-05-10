package gitsub

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Clone(repoURL string, directories []string, branch, output string) error {
	if len(directories) == 0 {
		return fmt.Errorf("no directories provided")
	}

	normalized, err := normalizeDirectories(directories)
	if err != nil {
		return err
	}

	if output == "" {
		output = extractRepoName(repoURL)
	}
	if output == "" {
		return fmt.Errorf("unable to determine output directory")
	}

	if branch == "" {
		branch = "main"
	}

	output = filepath.Clean(output)

	if dirExists(output) {
		return fmt.Errorf("output directory already exists: %s", output)
	}

	fmt.Printf("Cloning repository: %s\n", repoURL)
	fmt.Printf("Output directory: %s\n", output)
	fmt.Printf("Branch: %s\n", branch)

	if err := RunGit("init", output); err != nil {
		return err
	}

	if err := RunGitInDir(output, "remote", "add", "origin", repoURL); err != nil {
		return err
	}

	if err := RunGitInDir(output, "config", "core.sparseCheckout", "true"); err != nil {
		return err
	}

	args := append([]string{"sparse-checkout", "set", "--cone"}, normalized...)
	if err := RunGitInDir(output, args...); err != nil {
		return err
	}

	fmt.Println("Fetching required files...")
	if err := RunGitInDir(output, "fetch", "--filter=blob:none", "--depth=1", "origin", branch); err != nil {
		return err
	}

	if err := RunGitInDir(output, "checkout", branch); err != nil {
		return err
	}

	fmt.Println("Done.")
	return nil
}

func extractRepoName(url string) string {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func normalizeDirectories(directories []string) ([]string, error) {
	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(directories))

	for _, dir := range directories {
		nd := NormalizePath(dir)
		if nd == "" {
			return nil, fmt.Errorf("invalid directory path: %q", dir)
		}
		if _, ok := seen[nd]; ok {
			continue
		}
		seen[nd] = struct{}{}
		normalized = append(normalized, nd)
	}

	return normalized, nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
