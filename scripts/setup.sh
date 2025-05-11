#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${GREEN}[*]${NC} $1"
}

# Function to print error messages
print_error() {
    echo -e "${RED}[!]${NC} $1"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_error "Please run as root (use sudo)"
    exit 1
fi

print_status "Setting up Crobrew..."

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
else
    OS=$(uname -s)
fi

# Install Go if not present
if ! command -v go &> /dev/null; then
    print_status "Installing Go..."
    apt-get update
    apt-get install -y golang-go
fi

# Install Crobrew
print_status "Installing Crobrew..."
GOBIN=/usr/local/bin go install github.com/chersbobers/crobrew@latest

# Create symlink if needed
if [ ! -f /usr/local/bin/cro ]; then
    ln -s $(which cro) /usr/local/bin/cro
fi

# Set permissions
chmod +x /usr/local/bin/cro

print_status "Installation complete! You can now use 'cro' commands."
print_status "Try 'cro help' to get started."
