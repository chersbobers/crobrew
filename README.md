# Crobrew

A package manger for chromeos. Simple, Lightweight.
Crobrew is still in its early stages some things might not work.

## Installation

You can install Crobrew directly using Go:

```bash
go install github.com/chersbobers/crobrew@latest
```

### Prerequisites

- Go 1.16 or later
- For ChromeOS users:
  1. Enable Linux (Beta) in ChromeOS Settings
  2. Install Go in the Linux container: `sudo apt-get update && sudo apt-get install golang-go`

## Usage

Once installed, you can run `crobrew` from your terminal. The application will automatically detect your system's package manager:

- ChromeOS: Uses apt or dnf
## Features

- Update package lists
- Search for packages
- Automatic package manager detection
