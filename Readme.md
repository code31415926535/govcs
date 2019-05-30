# Go Version Control System

This is an educational project with no practical purpose in mind.
The primary goal here is to write a simple version control system.

## Features

- Create repository with no configuration
- Check current status

## Usage

```sh
govcs init # init empty repository
govcs add <file> # add file to staging area
govcs stat # check status
```

## Planned

- Create commit from changed files
  - Store file diffs
  - Store commit information
- List commits
- Checkout specific commit
- End-to-end tests
- Abstract file system