# LazyLinux

LazyLinux is a command-line tool that simplifies Linux package management across different distributions. It provides a unified interface for installing, removing, updating, and managing packages through your system's native package manager and optional Flatpak support.

## Features

- **Unified Commands** - Simple `install`, `remove`, `update`, and `clean` commands across all distros
- **Multi-Source Support** - Works with native package managers (DNF, APT, Pacman) and Flatpak
- **Interactive Source Selection** - Automatically prompts when packages are available from multiple sources
- **Auto-Detection** - Detects your distribution's package manager during setup
- **Clean Output** - Human-readable console messages with clear status indicators

## Supported Package Managers

- **DNF** (Fedora, RHEL, CentOS)
- **APT** (Debian, Ubuntu, Linux Mint)
- **Pacman** (Arch Linux, Manjaro)
- **Flatpak** (optional, cross-distribution)

## Installation

### Step 1: Clone and Install
git clone https://github.com/<your-username>/lazylinux.git
cd lazylinux
./install.sh

### Step 2: Configure

Run the setup command to auto-detect your package manager:

lazylinux setup

This creates `~/.config/lazylinux/config.yaml` with your system's configuration.

### Step 3: (Optional) Add Alias

For quicker usage, add an alias to your shell:

alias lzl="lazylinux"
