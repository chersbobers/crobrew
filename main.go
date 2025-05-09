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
	remove  string
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
			remove:  "sudo apt-get remove",
		},
		{
			name:    "dnf",
			search:  "dnf search",
			update:  "sudo dnf check-update",
			install: "sudo dnf install",
			remove:  "sudo dnf remove",
		},
		// Example: Add pacman package manager
		{
			name:    "pacman",
			search:  "pacman -Ss",
			update:  "sudo pacman -Sy",
			install: "sudo pacman -S",
			remove:  "sudo pacman -R",
		},
		{
			name:    "snap",
			search:  "snap find",
			update:  "sudo snap refresh",
			install: "sudo snap install",
			remove:  "sudo snap remove",
		},
	},
	"windows": {
		{
			name:    "wsl-apt",
			search:  "wsl apt-cache search",
			update:  "wsl sudo apt-get update",
			install: "wsl sudo apt-get install",
			remove:  "wsl sudo apt-get remove",
		},
		{
			name:    "choco",
			search:  "choco search",
			update:  "choco upgrade all -y",
			install: "choco install",
			remove:  "choco uninstall",
		},
	},
	"darwin": {
		{
			name:    "brew",
			search:  "brew search",
			update:  "brew update",
			install: "brew install",
			remove:  "brew uninstall",
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

func installPackage(packageName string) error {
	if defaultManager == nil {
		detectPackageManager()
	}
	cmdParts := strings.Split(defaultManager.install, " ")
	cmdParts = append(cmdParts, packageName)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error installing package: %v\nOutput: %s\nThis might be because:\n1. You don't have sudo permissions\n2. The package doesn't exist\n3. The package manager is not available", err, string(output))
	}
	return nil
}

func removePackage(packageName string) error {
	if defaultManager == nil {
		detectPackageManager()
	}
	cmdParts := strings.Split(defaultManager.remove, " ")
	cmdParts = append(cmdParts, packageName)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error removing package: %v\nOutput: %s\nThis might be because:\n1. You don't have sudo permissions\n2. The package doesn't exist\n3. The package manager is not available", err, string(output))
	}
	return nil
}

func showHelp() {
	fmt.Println("Crobrew - ChromeOS Package Manager")
	fmt.Println("\nUsage:")
	fmt.Println("  cro <command> [package]")
	fmt.Println("\nCommands:")
	fmt.Println("  update             Update package list")
	fmt.Println("  search <query>     Search for packages")
	fmt.Println("  install <package>  Install a package")
	fmt.Println("  remove <package>   Remove a package")
	fmt.Println("  help              Show this help message")
	fmt.Println("  interactive       Start interactive mode")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// If no arguments, start interactive mode
		startInteractiveMode()
		return
	}

	command := args[0]
	switch command {
	case "update":
		fmt.Println("Updating package list...")
		if err := updatePackageList(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Package list updated successfully!")

	case "search":
		if len(args) < 2 {
			fmt.Println("Usage: cro search <query>")
			os.Exit(1)
		}
		query := args[1]
		fmt.Println("Searching packages...")
		results, err := searchPackages(query)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("\nAvailable packages:")
		fmt.Println(results)

	case "install":
		if len(args) < 2 {
			fmt.Println("Usage: cro install <package>")
			os.Exit(1)
		}
		packageName := args[1]
		fmt.Printf("Installing package %s...\n", packageName)
		if err := installPackage(packageName); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Package %s installed successfully!\n", packageName)

	case "remove":
		if len(args) < 2 {
			fmt.Println("Usage: cro remove <package>")
			os.Exit(1)
		}
		packageName := args[1]
		fmt.Printf("Removing package %s...\n", packageName)
		if err := removePackage(packageName); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Package %s removed successfully!\n", packageName)

	case "help":
		showHelp()

	case "interactive":
		startInteractiveMode()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func startInteractiveMode() {
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
		fmt.Println("3. Install package")
		fmt.Println("4. Remove package")
		fmt.Println("5. Exit")
		fmt.Print("\nChoose an option (1-5): ")

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
			fmt.Print("Enter package name to install: ")
			packageName, _ := reader.ReadString('\n')
			packageName = strings.TrimSpace(packageName)

			if packageName == "" {
				fmt.Println("Package name cannot be empty")
				continue
			}

			fmt.Printf("Installing package %s...\n", packageName)
			if err := installPackage(packageName); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("Package %s installed successfully!\n", packageName)

		case "4":
			fmt.Print("Enter package name to remove: ")
			packageName, _ := reader.ReadString('\n')
			packageName = strings.TrimSpace(packageName)

			if packageName == "" {
				fmt.Println("Package name cannot be empty")
				continue
			}

			fmt.Printf("Removing package %s...\n", packageName)
			if err := removePackage(packageName); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("Package %s removed successfully!\n", packageName)

		case "5":
			fmt.Println("Thank you for using Crobrew!")
			os.Exit(0)

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
