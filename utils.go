package main

import (
	"fmt"
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
