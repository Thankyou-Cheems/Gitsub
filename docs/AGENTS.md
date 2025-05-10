# Repository Guidelines

This repository currently contains the MVP design document for **Gitsub**, a planned Go CLI that clones only selected directories from large GitHub repositories using sparse checkout. The codebase is not yet present; use this guide to keep contributions aligned with the documented MVP scope.

## Project Structure & Module Organization

Planned structure (per `Gitsub.md`):

```
gitsub/
├── main.go
├── clone.go
├── git.go
├── utils.go
├── go.mod
└── README.md
```

- `main.go` handles CLI parsing.
- `clone.go` implements the sparse-clone workflow.
- `git.go` wraps git command execution.
- `utils.go` includes validation and path helpers.

## Build, Test, and Development Commands

When the Go code exists, use:

- `go mod init github.com/<you>/gitsub` to initialize modules (once).
- `go run . clone <repo-url> <dir...>` to run the CLI locally.
- `go build -o gitsub` to produce a local binary.
- Cross-compile examples:
  - `GOOS=linux GOARCH=amd64 go build -o gitsub-linux`
  - `GOOS=darwin GOARCH=amd64 go build -o gitsub-macos`
  - `GOOS=windows GOARCH=amd64 go build -o gitsub.exe`

## Coding Style & Naming Conventions

- Language: Go (standard library preferred).
- Indentation: tabs (Go default).
- Naming: `CamelCase` for exported, `camelCase` for unexported, `ALL_CAPS` for constants.
- Keep files small and focused (see module responsibilities above).
- Run `gofmt` on all Go files before committing.

## Testing Guidelines

No test framework is defined yet. When adding tests:

- Use Go’s standard `testing` package.
- Name files `*_test.go`.
- Prefer table-driven tests and explicit error messages.
- Add at least one test for URL validation and path normalization.

## Commit & Pull Request Guidelines

No commit convention is established yet. Until one is defined:

- Use short, imperative messages (e.g., `Add sparse-checkout clone flow`).
- Keep commits scoped to a single change.
- PRs should include:
  - A brief description of behavior changes.
  - Example commands used for validation.
  - Any known limitations or follow-up work.

## Security & Configuration Tips

- The tool expects Git 2.25+ with network access.
- Avoid embedding tokens or secrets in command history or logs.
- Validate repository URLs before running git commands.
