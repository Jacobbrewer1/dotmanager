# Dot Manager

## Overview

Dot Manager is a simple tool to manage dotfiles. It allows you to easily add, remove, and update dotfiles in your home
directory by syncing them with a git repository.

## Installation

To install Dot Manager, clone the repository and run the install script:

```bash
go install github.com/jacobbrewer1/dotmanager/cmd@latest
```

## Usage

### Add

To add a dotfile to the repository, use the `add` command:

```bash
dotmanager add
```

This will then direct you to select the file you would like to add through a TUI.

### Remove

To stop tracking a dotfile, simply delete it from the repository.

### Pull

To pull the latest changes from you local machine into the repository, use the `pull` command:

```bash
dotmanager pull
```

This will then update all the tracked files in the repository with the latest changes from your local machine.

### Push

To push the latest changes from the repository to your local machine, use the `push` command:

```bash
dotmanager push
```

This will then update all the tracked files in your local machine with the latest changes from the repository.

### Version

To check the version of Dot Manager, use the `version` command:

```bash
dotmanager version
```
