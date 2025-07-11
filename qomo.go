package main

import (
	"fmt"
	"os"
	"path/filepath"

	"qomoboro/internal/storage"
	"qomoboro/internal/ui"
)

const (
	appName = "qomoboro"
	version = "2.0.0"
)

func main() {
	// Get data directory
	dataDir, err := getDataDir()
	if err != nil {
		fmt.Printf("Error getting data directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize storage
	store, err := storage.NewFileStorage(dataDir)
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// Handle command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("%s %s\n", appName, version)
			return
		case "help", "--help", "-h":
			showHelp()
			return
		case "data-dir", "--data-dir":
			fmt.Printf("Data directory: %s\n", dataDir)
			return
		case "backup", "--backup":
			if err := store.Backup(); err != nil {
				fmt.Printf("Error creating backup: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Backup created successfully")
			return
		}
	}

	// Create and run the TUI application
	app := ui.NewApp(store)
	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}

// getDataDir returns the data directory for the application
func getDataDir() (string, error) {
	// Try XDG_DATA_HOME first
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, appName), nil
	}

	// Fall back to ~/.local/share
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".local", "share", appName), nil
}

// showHelp displays usage information
func showHelp() {
	fmt.Printf(`%s %s - Canonical Hours Task Manager

A systematic task management tool built with Go and Charmbracelet TUI components,
designed around canonical hours and gamified productivity tracking.

USAGE:
    %s [command]

COMMANDS:
    help, --help, -h     Show this help message
    version, --version, -v Show version information
    data-dir, --data-dir Show data directory path
    backup, --backup     Create a backup of all data

NAVIGATION (in TUI):
    [t] Tasks       - Manage your tasks
    [s] Schedule    - View canonical hours
    [d] Statistics  - View productivity stats
    [c] Create      - Add new task
    [g] Settings    - Configure app
    [r] Refresh     - Reload data
    [q] Quit        - Exit application

TASK MANAGEMENT:
    ↑/↓ or k/j      - Navigate lists
    Enter           - View details
    Space           - Toggle task status
    c               - Create new task
    d               - Delete task
    q/Esc           - Go back

SCORING SYSTEM:
    Each task is scored on three dimensions (0-5):
    - Work: Business/professional productivity value
    - Play: Recreation/leisure/enjoyment value
    - Learn: Educational/skill development value

CANONICAL HOURS:
    Traditional time blocks adapted for modern productivity:
    - Matins (06:00-07:30): Deep work, planning
    - Lauds (07:30-09:00): Administrative tasks
    - Prime (09:00-12:00): High-focus work blocks
    - Terce (12:00-13:30): Meetings, collaboration
    - Sext (13:30-15:00): Lunch, recovery
    - None (15:00-16:30): Creative work, experimentation
    - Vespers (16:30-18:00): Learning, documentation
    - Compline (18:00-20:00): Planning, reflection

DATA LOCATION:
    Configuration and data are stored in:
    - Linux/macOS: ~/.local/share/qomoboro/
    - Or $XDG_DATA_HOME/qomoboro/ if XDG_DATA_HOME is set

EXAMPLES:
    %s                    # Start the TUI application
    %s --version          # Show version
    %s --data-dir         # Show data directory
    %s --backup           # Create backup

For more information, visit: https://github.com/QRY91/qomoboro
`, appName, version, appName, appName, appName, appName, appName)
}
