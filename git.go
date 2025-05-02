package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunGitInDir(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CheckGitVersion() error {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git not found or not installed")
	}

	version := string(output)
	if !strings.Contains(version, "git version") {
		return fmt.Errorf("unable to detect git version")
	}

	return nil
}
