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

## Requirements

- Git 2.25+

## How it works

Gitsub uses sparse checkout with partial clone:

```bash
git clone --filter=blob:none --sparse <repo>
git sparse-checkout set <directory>
```
