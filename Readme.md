# Go Version Control System

> NOTE: This repo is currently work in progress. Please check back in
> a couple of days.

This is an educational project with no practical purpose in mind.
The primary goal here is to write a simple version control system.

## Features

- Create repository with no configuration
- Check current status
- Add files and commit changes
- List Commits

## Usage

```sh
govcs init # init empty repository
govcs add <file> # add file to staging area
govcs remove <file> # remove file from staging area
govcs stat # check status
govcs commit <message> # commit changes
govcs list-commits # list commits
```

## Planned

- End-to-end tests
- Checkout commit
- Implement branching
