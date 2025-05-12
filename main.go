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
	// version field removed as it's currently unused
}

func (pm *PackageManager) checkAvailable() bool {
	var cmd *exec.Cmd

	if strings.HasPrefix(pm.name, "wsl") {
		cmd = exec.Command("wsl", "command", "-v", strings.Split(pm.search, " ")[1])
	} else {
		cmdParts := strings.Split(pm.search, " ")
		cmd = exec.Command(cmdParts[0], "--version")
	}

	return cmd.Run() == nil
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
	// Detect the package manager during initialization
	detectPackageManager()
	if defaultManager == nil {
		fmt.Println("Error: No valid package manager detected. Please ensure a supported package manager is installed.")
		os.Exit(1)
	}
}

func detectPackageManager() {
	os := runtime.GOOS
	fmt.Printf("Detected OS: %s\n", os)
	managers, ok := packageManagers[os]

	if !ok || len(managers) == 0 {
		fmt.Println("No package managers found for OS, defaulting to Linux")
		managers = packageManagers["linux"]
	}

	for i := range managers {
		pm := &managers[i]
		fmt.Printf("Checking for %s...\n", pm.name)

		if pm.checkAvailable() {
			defaultManager = pm
			fmt.Printf("Selected package manager: %s\n", pm.name)
			// Update package list after detecting package manager
			if err := updatePackageList(); err != nil {
				fmt.Printf("Warning: Failed to update package list: %v\n", err)
			}
			return
		}
	}

	fmt.Println("No working package manager found, defaulting to apt")
	defaultManager = &packageManagers["linux"][0]
}

func searchPackages(query string) (string, error) {
	if defaultManager == nil {
		detectPackageManager()
	}
	cmdParts := strings.Split(defaultManager.search, " ")
	cmdParts = append(cmdParts, query)

	return executeCommand(cmdParts)
}

func updatePackageList() error {
	if defaultManager == nil {
		detectPackageManager()
	}
	cmdParts := strings.Split(defaultManager.update, " ")

	output, err := executeCommand(cmdParts)
	if err != nil {
		return fmt.Errorf("failed to update package list: %v", err)
	}
	fmt.Print(output)
	return nil
}

func updateAllPackageManagers() error {
	goos := runtime.GOOS
	managers := packageManagers[goos]

	if len(managers) == 0 {
		return fmt.Errorf("no package managers found for OS: %s", goos)
	}

	var errors []string
	updated := false

	for _, pm := range managers {
		fmt.Printf("\nUpdating %s...\n", pm.name)
		cmdParts := strings.Split(pm.update, " ")
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", pm.name, err))
			continue
		}
		updated = true
		fmt.Printf("%s updated successfully!\n", pm.name)
	}

	if !updated {
		return fmt.Errorf("failed to update any package managers:\n%s", strings.Join(errors, "\n"))
	}

	if len(errors) > 0 {
		fmt.Printf("\nWarning: Some updates failed:\n%s\n", strings.Join(errors, "\n"))
	}

	return nil
}

func executeCommand(cmdParts []string) (string, error) {
	if len(cmdParts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" && strings.HasPrefix(cmdParts[0], "wsl") {
		// Remove "wsl" prefix and add -e flag
		cmdParts = cmdParts[1:] // Remove "wsl"
		wslArgs := append([]string{"-e"}, cmdParts...)
		cmd = exec.Command("wsl", wslArgs...)
	} else {
		program := cmdParts[0]
		args := []string{}
		if len(cmdParts) > 1 {
			args = cmdParts[1:]
		}
		cmd = exec.Command(program, args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("%v: %s", err, output)
	}

	return string(output), nil
}

func installPackage(packageName string) error {
    if defaultManager == nil {
        fmt.Println("Error: No valid package manager detected.")
        os.Exit(1)
    }

    cmdParts := strings.Split(defaultManager.install, " ")
    cmdParts = append(cmdParts, packageName)

    fmt.Printf("Executing install command: %s\n", strings.Join(cmdParts, " "))
    output, err := executeCommand(cmdParts)
    if err != nil {
        return fmt.Errorf("failed to install %s: %v", packageName, err)
    }
    fmt.Print(output)
    return nil
}

func removePackage(packageName string) error {
	if defaultManager == nil {
		detectPackageManager()
	}
	cmdParts := strings.Split(defaultManager.remove, " ")
	cmdParts = append(cmdParts, packageName)

	output, err := executeCommand(cmdParts)
	if err != nil {
		return fmt.Errorf("failed to remove %s: %v", packageName, err)
	}
	fmt.Print(output)
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
		fmt.Println("Updating all package managers...")
		if err := updateAllPackageManagers(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All package managers updated successfully!")

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
			fmt.Println("Updating all package managers...")
			if err := updateAllPackageManagers(); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Println("All package managers updated successfully!")

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
