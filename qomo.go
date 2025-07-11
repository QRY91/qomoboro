package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"qomoboro/internal/models"
	"qomoboro/internal/storage"
)

const (
	appName = "qomoboro"
	version = "2.0.0"
	ascii   = `
 ██████  ██████  ███    ███  ██████  ██████   ██████  ██████   ██████
██    ██ ██    ██ ████  ████ ██    ██ ██   ██ ██    ██ ██   ██ ██    ██
██    ██ ██    ██ ██ ████ ██ ██    ██ ██████  ██    ██ ██████  ██    ██
██ ▄▄ ██ ██    ██ ██  ██  ██ ██    ██ ██   ██ ██    ██ ██   ██ ██    ██
 ██████   ██████  ██      ██  ██████  ██████   ██████  ██   ██  ██████
    ▀▀
        Canonical Hours Task Manager - CLI Edition
`
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

	// Parse command line arguments
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	cmd := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	switch cmd {
	case "add", "task", "new":
		handleAddTask(store, args)
	case "list", "ls":
		handleListTasks(store, args)
	case "complete", "done":
		handleCompleteTask(store, args)
	case "delete", "rm":
		handleDeleteTask(store, args)
	case "status", "stat":
		handleStatus(store, args)
	case "schedule", "sched":
		handleSchedule(store, args)
	case "stats":
		handleStats(store, args)
	case "tui", "interactive":
		fmt.Println("TUI mode not implemented yet. Use CLI commands for now.")
	case "version", "--version", "-v":
		fmt.Printf("%s %s\n", appName, version)
	case "help", "--help", "-h":
		showHelp()
	case "backup":
		handleBackup(store)
	case "data-dir":
		fmt.Printf("Data directory: %s\n", dataDir)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		showHelp()
		os.Exit(1)
	}
}

func handleAddTask(store storage.Storage, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: qomoboro add \"<title>\" [description] [work] [play] [learn]")
		fmt.Println("Example: qomoboro add \"Fix bug in API\" \"Memory leak in handler\" 4 1 3")
		return
	}

	title := args[0]
	description := ""
	work, play, learn := 0, 0, 0

	// Parse optional arguments
	if len(args) > 1 {
		description = args[1]
	}
	if len(args) > 2 {
		if w, err := strconv.Atoi(args[2]); err == nil && w >= 0 && w <= 5 {
			work = w
		}
	}
	if len(args) > 3 {
		if p, err := strconv.Atoi(args[3]); err == nil && p >= 0 && p <= 5 {
			play = p
		}
	}
	if len(args) > 4 {
		if l, err := strconv.Atoi(args[4]); err == nil && l >= 0 && l <= 5 {
			learn = l
		}
	}

	task := &models.Task{
		ID:          fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Title:       title,
		Description: description,
		Score: models.Score{
			Work:  work,
			Play:  play,
			Learn: learn,
		},
		Status:    models.TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := store.CreateTask(task); err != nil {
		fmt.Printf("Error creating task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Task created: %s\n", title)
	fmt.Printf("   Scores: Work %d, Play %d, Learn %d\n", work, play, learn)
}

func handleListTasks(store storage.Storage, args []string) {
	tasks, err := store.ListTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks yet. Create one with: qomoboro add \"Task title\"")
		return
	}

	fmt.Printf("%s\n", ascii)
	fmt.Printf("📋 Your Tasks (%d total)\n", len(tasks))
	fmt.Println(strings.Repeat("─", 60))

	pending := 0
	completed := 0

	for i, task := range tasks {
		status := getStatusEmoji(task.Status)
		scores := fmt.Sprintf("W:%d P:%d L:%d", task.Score.Work, task.Score.Play, task.Score.Learn)

		fmt.Printf("%2d. %s %s %s\n", i+1, status, task.Title, colorize(scores, "dim"))

		if task.Description != "" {
			fmt.Printf("     %s\n", colorize(task.Description, "dim"))
		}

		if task.Status == models.TaskStatusPending {
			pending++
		} else if task.Status == models.TaskStatusCompleted {
			completed++
		}
	}

	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("📊 Status: %d pending, %d completed\n", pending, completed)
}

func handleCompleteTask(store storage.Storage, args []string) {
	tasks, err := store.ListTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	// Filter to pending tasks only
	var pendingTasks []*models.Task
	for _, task := range tasks {
		if task.Status == models.TaskStatusPending {
			pendingTasks = append(pendingTasks, task)
		}
	}

	if len(pendingTasks) == 0 {
		fmt.Println("🎉 No pending tasks! All done.")
		return
	}

	var task *models.Task

	if len(args) == 0 {
		// Interactive selection
		fmt.Println("📋 Pending tasks:")
		for i, t := range pendingTasks {
			scores := fmt.Sprintf("W:%d P:%d L:%d", t.Score.Work, t.Score.Play, t.Score.Learn)
			fmt.Printf("%2d. %s %s\n", i+1, t.Title, colorize(scores, "dim"))
		}
		fmt.Print("\nWhich task to complete? (1-" + fmt.Sprintf("%d", len(pendingTasks)) + "): ")

		var choice int
		if _, err := fmt.Scanf("%d", &choice); err != nil || choice < 1 || choice > len(pendingTasks) {
			fmt.Println("❌ Invalid selection")
			return
		}
		task = pendingTasks[choice-1]
	} else {
		// Try to parse as number first
		if taskNum, err := strconv.Atoi(args[0]); err == nil {
			if taskNum < 1 || taskNum > len(pendingTasks) {
				fmt.Printf("❌ Task number %d out of range (1-%d pending tasks)\n", taskNum, len(pendingTasks))
				return
			}
			task = pendingTasks[taskNum-1]
		} else {
			// Try partial title matching
			query := strings.ToLower(args[0])
			var matches []*models.Task
			for _, t := range pendingTasks {
				if strings.Contains(strings.ToLower(t.Title), query) {
					matches = append(matches, t)
				}
			}

			if len(matches) == 0 {
				fmt.Printf("❌ No pending tasks match '%s'\n", args[0])
				return
			} else if len(matches) == 1 {
				task = matches[0]
			} else {
				fmt.Printf("🤔 Multiple tasks match '%s':\n", args[0])
				for i, t := range matches {
					fmt.Printf("%2d. %s\n", i+1, t.Title)
				}
				fmt.Print("Which one? (1-" + fmt.Sprintf("%d", len(matches)) + "): ")

				var choice int
				if _, err := fmt.Scanf("%d", &choice); err != nil || choice < 1 || choice > len(matches) {
					fmt.Println("❌ Invalid selection")
					return
				}
				task = matches[choice-1]
			}
		}
	}

	// Confirm completion
	scores := fmt.Sprintf("W:%d P:%d L:%d", task.Score.Work, task.Score.Play, task.Score.Learn)
	fmt.Printf("✅ Completing: %s %s\n", task.Title, colorize(scores, "dim"))
	if task.Description != "" {
		fmt.Printf("   %s\n", colorize(task.Description, "dim"))
	}

	task.Complete()

	if err := store.UpdateTask(task); err != nil {
		fmt.Printf("Error updating task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🎉 Done! %s\n", task.Title)
}

func handleDeleteTask(store storage.Storage, args []string) {
	tasks, err := store.ListTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("📝 No tasks to delete")
		return
	}

	var task *models.Task

	if len(args) == 0 {
		// Interactive selection
		fmt.Println("📋 All tasks:")
		for i, t := range tasks {
			status := getStatusEmoji(t.Status)
			scores := fmt.Sprintf("W:%d P:%d L:%d", t.Score.Work, t.Score.Play, t.Score.Learn)
			fmt.Printf("%2d. %s %s %s\n", i+1, status, t.Title, colorize(scores, "dim"))
		}
		fmt.Print("\nWhich task to delete? (1-" + fmt.Sprintf("%d", len(tasks)) + "): ")

		var choice int
		if _, err := fmt.Scanf("%d", &choice); err != nil || choice < 1 || choice > len(tasks) {
			fmt.Println("❌ Invalid selection")
			return
		}
		task = tasks[choice-1]
	} else {
		// Try to parse as number first
		if taskNum, err := strconv.Atoi(args[0]); err == nil {
			if taskNum < 1 || taskNum > len(tasks) {
				fmt.Printf("❌ Task number %d out of range (1-%d)\n", taskNum, len(tasks))
				return
			}
			task = tasks[taskNum-1]
		} else {
			// Try partial title matching
			query := strings.ToLower(args[0])
			var matches []*models.Task
			for _, t := range tasks {
				if strings.Contains(strings.ToLower(t.Title), query) {
					matches = append(matches, t)
				}
			}

			if len(matches) == 0 {
				fmt.Printf("❌ No tasks match '%s'\n", args[0])
				return
			} else if len(matches) == 1 {
				task = matches[0]
			} else {
				fmt.Printf("🤔 Multiple tasks match '%s':\n", args[0])
				for i, t := range matches {
					status := getStatusEmoji(t.Status)
					fmt.Printf("%2d. %s %s\n", i+1, status, t.Title)
				}
				fmt.Print("Which one to delete? (1-" + fmt.Sprintf("%d", len(matches)) + "): ")

				var choice int
				if _, err := fmt.Scanf("%d", &choice); err != nil || choice < 1 || choice > len(matches) {
					fmt.Println("❌ Invalid selection")
					return
				}
				task = matches[choice-1]
			}
		}
	}

	// Confirm deletion
	status := getStatusEmoji(task.Status)
	scores := fmt.Sprintf("W:%d P:%d L:%d", task.Score.Work, task.Score.Play, task.Score.Learn)
	fmt.Printf("🗑️  Deleting: %s %s %s\n", status, task.Title, colorize(scores, "dim"))
	if task.Description != "" {
		fmt.Printf("   %s\n", colorize(task.Description, "dim"))
	}

	fmt.Print("Are you sure? (y/N): ")
	var confirm string
	fmt.Scanf("%s", &confirm)
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("❌ Cancelled")
		return
	}

	if err := store.DeleteTask(task.ID); err != nil {
		fmt.Printf("Error deleting task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🗑️  Deleted: %s\n", task.Title)
}

func handleStatus(store storage.Storage, args []string) {
	schedule, err := store.GetSchedule()
	if err != nil {
		fmt.Printf("Error loading schedule: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	currentHour := schedule.GetCurrentHour(now)

	fmt.Printf("%s\n", ascii)
	fmt.Printf("🕐 Current Time: %s\n", now.Format("15:04 Monday, Jan 2"))

	if currentHour != nil {
		fmt.Printf("📍 Current Hour: %s (%s)\n", currentHour.Name, currentHour.Description)
		fmt.Printf("   Suggested Focus: %s\n", currentHour.Purpose)
	} else {
		fmt.Printf("📍 Outside canonical hours\n")
	}

	// Show task summary
	tasks, err := store.ListTasks()
	if err == nil {
		pending := 0
		completed := 0
		totalWork, totalPlay, totalLearn := 0, 0, 0

		for _, task := range tasks {
			if task.Status == models.TaskStatusPending {
				pending++
			} else if task.Status == models.TaskStatusCompleted {
				completed++
				totalWork += task.Score.Work
				totalPlay += task.Score.Play
				totalLearn += task.Score.Learn
			}
		}

		fmt.Printf("\n📊 Today's Progress:\n")
		fmt.Printf("   Tasks: %d pending, %d completed\n", pending, completed)
		fmt.Printf("   Scores: Work %d, Play %d, Learn %d\n", totalWork, totalPlay, totalLearn)
	}
}

func handleSchedule(store storage.Storage, args []string) {
	schedule, err := store.GetSchedule()
	if err != nil {
		fmt.Printf("Error loading schedule: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", ascii)
	fmt.Printf("🗓️  Canonical Hours Schedule\n")
	fmt.Println(strings.Repeat("─", 60))

	now := time.Now()
	for _, hour := range schedule.Hours {
		marker := "  "
		if hour.IsActive(now) {
			marker = "👉"
		}

		fmt.Printf("%s %s - %s: %s\n", marker, hour.StartTime, hour.EndTime, hour.Name)
		fmt.Printf("     %s\n", hour.Description)
		if hour.Purpose != "" {
			fmt.Printf("     Focus: %s\n", colorize(hour.Purpose, "dim"))
		}
		fmt.Println()
	}
}

func handleStats(store storage.Storage, args []string) {
	today := time.Now()
	stats, err := store.GetDailyStats(today)
	if err != nil {
		fmt.Printf("Error loading stats: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", ascii)
	fmt.Printf("📈 Statistics for %s\n", today.Format("Monday, Jan 2, 2006"))
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("Tasks: %d total, %d completed (%.1f%%)\n",
		stats.TotalTasks, stats.CompletedTasks, stats.CompletionRate())
	fmt.Printf("Scores: Work %d, Play %d, Learn %d\n",
		stats.TotalScore.Work, stats.TotalScore.Play, stats.TotalScore.Learn)
	fmt.Printf("Time: %s\n", stats.TimeSpent.String())
}

func handleBackup(store storage.Storage) {
	if err := store.Backup(); err != nil {
		fmt.Printf("Error creating backup: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("💾 Backup created successfully")
}

func getStatusEmoji(status models.TaskStatus) string {
	switch status {
	case models.TaskStatusPending:
		return "⏳"
	case models.TaskStatusActive:
		return "🔄"
	case models.TaskStatusCompleted:
		return "✅"
	case models.TaskStatusCancelled:
		return "❌"
	case models.TaskStatusPaused:
		return "⏸️"
	default:
		return "❓"
	}
}

func colorize(text, style string) string {
	// Simple color codes - can be enhanced later
	switch style {
	case "dim":
		return fmt.Sprintf("\033[2m%s\033[0m", text)
	case "bold":
		return fmt.Sprintf("\033[1m%s\033[0m", text)
	case "green":
		return fmt.Sprintf("\033[32m%s\033[0m", text)
	default:
		return text
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
	fmt.Printf(`%s

%s %s - Canonical Hours Task Manager

USAGE:
    %s <command> [arguments]

COMMANDS:
    add <title> [description] [work] [play] [learn]
        Create a new task with optional scores (0-5)
        Example: %s add "Fix API bug" "Memory leak" 4 1 3

    list
        Show all tasks with their status and scores

    complete <number>
        Mark a task as completed (use task number from list)

    delete <number>
        Delete a task (use task number from list)

    status
        Show current canonical hour and task summary

    schedule
        Display the canonical hours schedule

    stats
        Show today's productivity statistics

    backup
        Create a backup of all data

    data-dir
        Show data directory location

    version
        Show version information

    help
        Show this help message

EXAMPLES:
    %s add "Review code" "PR #123" 5 2 3
    %s list
    %s complete              # Interactive selection
    %s complete bug          # Complete task matching "bug"
    %s complete 1            # Complete task #1
    %s delete old            # Delete task matching "old"
    %s status

CANONICAL HOURS:
    Matins    06:00-07:30  Deep work, planning
    Lauds     07:30-09:00  Administrative tasks
    Prime     09:00-12:00  High-focus work blocks
    Terce     12:00-13:30  Meetings, collaboration
    Sext      13:30-15:00  Lunch, recovery
    None      15:00-16:30  Creative work, experimentation
    Vespers   16:30-18:00  Learning, documentation
    Compline  18:00-20:00  Planning, reflection

SCORING SYSTEM:
    Each task gets three scores (0-5):
    Work:  Business/professional productivity value
    Play:  Recreation/leisure/enjoyment value
    Learn: Educational/skill development value

DATA LOCATION:
    Linux/macOS: ~/.local/share/qomoboro/
    Or: $XDG_DATA_HOME/qomoboro/ if XDG_DATA_HOME is set

For more information, visit: https://github.com/QRY91/qomoboro
`, ascii, appName, version, appName, appName, appName, appName, appName, appName)
}
