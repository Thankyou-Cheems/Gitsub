package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "clone":
		if err := runClone(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func runClone(args []string) error {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	branchShort := fs.String("b", "", "Branch name (default: main)")
	branchLong := fs.String("branch", "", "Branch name (default: main)")
	outputShort := fs.String("o", "", "Output directory (default: repo name)")
	outputLong := fs.String("output", "", "Output directory (default: repo name)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) < 1 {
		printCloneUsage()
		return fmt.Errorf("missing required arguments")
	}

	branch := firstNonEmpty(*branchShort, *branchLong)
	output := firstNonEmpty(*outputShort, *outputLong)

	if err := CheckGitVersion(); err != nil {
		return err
	}

	repoURL, directories, parsedBranch, ok := ParseGitHubDirURL(rest)
	if ok {
		if branch == "" {
			branch = parsedBranch
		}
	} else {
		if len(rest) < 2 {
			printCloneUsage()
			return fmt.Errorf("missing required arguments")
		}
		repoURL = rest[0]
		directories = rest[1:]
	}

	if err := ValidateRepoURL(repoURL); err != nil {
		return err
	}

	return Clone(repoURL, directories, branch, output)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  gitsub clone <repo-url> <directory> [directory...]")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  gitsub clone https://github.com/tensorflow/tensorflow tensorflow/core")
}

func printCloneUsage() {
	fmt.Println("Usage:")
	fmt.Println("  gitsub clone <repo-url> <directory> [directory...]")
	fmt.Println("Options:")
	fmt.Println("  -b, --branch <branch>   Specify branch (default: main)")
	fmt.Println("  -o, --output <path>     Output directory (default: repo name)")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
