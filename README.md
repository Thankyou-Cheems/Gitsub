# Gitsub

Clone only the directories you need from large GitHub repositories.

## Usage

```bash
gitsub clone <repo-url> <directory> [directory...]
```

Example:

```bash
gitsub clone https://github.com/tensorflow/tensorflow tensorflow/core
```

Or pass a GitHub directory URL directly:

```bash
gitsub clone https://github.com/tensorflow/tensorflow/tree/master/tensorflow/core
```

## Requirements

- Git 2.25+
- Go 1.21+ (to build)

## Build

```bash
go build ./cmd/gitsub
```

## Web Tool (GitHub Pages)

The static web tool lives in `docs/`. To publish via GitHub Pages, set the
repository's Pages source to the `docs/` folder on the `main` branch.

## How it works

Gitsub uses sparse checkout with partial clone:

```bash
git clone --filter=blob:none --sparse <repo>
git sparse-checkout set <directory>
```
