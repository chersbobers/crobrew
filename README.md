# Crobrew

A Simple, Lightweight. package manger
Crobrew is still in its early stages some things might not work.

## Installation

You can install Crobrew directly using Go:

```bash
go install github.com/chersbobers/crobrew@latest
```
You may need to use
```bash
echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> ~/.bashrc && source ~/.bashrc
```

### Prerequisites

- Go (any version)

#### ChromeOS Setup
1. Enable Linux (Beta) in ChromeOS Settings
2. Open Linux Terminal
3. Install Go:
   ```bash
   sudo apt-get update && sudo apt-get install golang-go
   ```
4. Add Go to your PATH:
   ```bash
   echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> ~/.bashrc && source ~/.bashrc
   ```

#### Linux Setup
1. Install Go using your distribution's package manager:
   - For Ubuntu/Debian:
     ```bash
     sudo apt-get update && sudo apt-get install golang-go
     ```
   - For Fedora:
     ```bash
     sudo dnf install golang
     ```
   - For Arch Linux:
     ```bash
     sudo pacman -S go
     ```
2. Add Go to your PATH:
   ```bash
   echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> ~/.bashrc && source ~/.bashrc
   ```

#### macOS Setup
1. Install Go using Homebrew:
   ```bash
   brew install go
   ```
2. Add Go to your PATH:
   ```bash
   echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> ~/.zshrc && source ~/.zshrc
   ```

#### Windows Setup
1. Download Go from [golang.org](https://golang.org/dl/)
2. Run the installer
3. Open Command Prompt and add Go to your PATH:
   ```cmd
   setx PATH "%PATH%;%USERPROFILE%\go\bin"
   ```

## Usage

Once installed, you can run `crobrew` from your terminal. The application will automatically detect your system's package manager:

- ChromeOS/Linux: Uses apt (Debian/Ubuntu) or dnf (Fedora)
- macOS: Uses Homebrew
- Windows: Uses Chocolatey (needs to be installed) or WSL
## Features

- Automatic package manager detection for all major platforms
- Update package lists
- Search for packages
- Install packages
- Remove packages
- Supports multiple package managers:
  - apt (Debian/Ubuntu/ChromeOS)
  - dnf (Fedora)
  - pacman (Arch Linux)
  - brew (macOS)
  - chocolatey (Windows)
  - snap (Linux)
  - WSL (Windows Subsystem for Linux)

## Supported Package Managers

| OS | Package Managers |
|----|-----------------|
| ChromeOS | apt, snap |
| Linux | apt, dnf, pacman, snap |
| macOS | brew |
| Windows | chocolatey, WSL-apt |
