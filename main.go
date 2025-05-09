package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// PackageManager defines the structure for package management commands
type PackageManager struct {
	search  string
	update  string
	install string
	name    string
}

// Package manager commands for different platforms
var packageManagers = map[string][]PackageManager{
	"linux": {
		{
			name:    "apt",
			search:  "apt-cache search",
			update:  "sudo apt-get update",
			install: "sudo apt-get install",
		},
		{
			name:    "dnf",
			search:  "dnf search",
			update:  "sudo dnf check-update",
			install: "sudo dnf install",
		},
	},
	"windows": {
		{
			name:    "wsl-apt",
			search:  "wsl apt-cache search",
			update:  "wsl sudo apt-get update",
			install: "wsl sudo apt-get install",
		},
		{
			name:    "choco",
			search:  "choco search",
			update:  "choco upgrade all -y",
			install: "choco install",
		},
	},
	"darwin": {
		{
			name:    "brew",
			search:  "brew search",
			update:  "brew update",
			install: "brew install",
		},
	},
}

var defaultManager *PackageManager

func init() {
	// Detect the package manager
	detectPackageManager()
}

func detectPackageManager() {
	os := runtime.GOOS
	managers := packageManagers[os]

	if len(managers) == 0 {
		// Default to Linux package managers for ChromeOS and others
		managers = packageManagers["linux"]
	}

	// Try each package manager until we find one that works
	for _, pm := range managers {
		cmdParts := strings.Split(pm.search, " ")
		cmd := exec.Command(cmdParts[0], "--version")
		if err := cmd.Run(); err == nil {
			defaultManager = &pm
			return
		}
	}

	// If no package manager is found, default to apt for ChromeOS
	defaultManager = &packageManagers["linux"][0]
}

func getPackageCommands() (string, string) {
	if defaultManager == nil {
		detectPackageManager()
	}
	return defaultManager.search, defaultManager.update
}

func searchPackages(query string) (string, error) {
	searchCmd, _ := getPackageCommands()
	cmdParts := strings.Split(searchCmd, " ")
	cmdParts = append(cmdParts, query)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error searching packages: %v\nThis might be because:\n1. You're not in a Linux environment\n2. The package manager is not available\n3. You don't have the required permissions", err)
	}
	return string(output), nil
}

func updatePackageList() error {
	_, updateCmd := getPackageCommands()
	cmdParts := strings.Split(updateCmd, " ")

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error updating package list: %v\nThis might be because:\n1. You're not in a Linux environment\n2. You don't have sudo permissions\n3. The package manager is not available", err)
	}
	return nil
}

func main() {
	fmt.Println("Welcome to Crobrew - ChromeOS Package Manager")
	fmt.Println(`
     ____           _
    / ___|_ __ ___ | |__  _ __ _____      __
   | |   | '__/ _ \| '_ \| '__/ _ \ \ /\ / /
   | |___| | | (_) | |_) | | |  __/\ V  V /
    \____|_|  \___/|_.__/|_|  \___| \_/\_/
    `)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Update package list")
		fmt.Println("2. Search packages")
		fmt.Println("3. Exit")
		fmt.Print("\nChoose an option (1-3): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Println("Updating package list...")
			if err := updatePackageList(); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("Package list updated successfully!")

		case "2":
			fmt.Print("Enter search term (or press Enter to list all): ")
			query, _ := reader.ReadString('\n')
			query = strings.TrimSpace(query)

			fmt.Println("Searching packages...")
			results, err := searchPackages(query)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			fmt.Println("\nAvailable packages:")
			fmt.Println(results)

		case "3":
			fmt.Println("Thank you for using Crobrew!")
			os.Exit(0)

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
