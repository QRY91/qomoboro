package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"qomoboro/internal/models"
)

// Storage defines the interface for data persistence
type Storage interface {
	// Task operations
	CreateTask(task *models.Task) error
	GetTask(id string) (*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id string) error
	ListTasks() ([]*models.Task, error)
	ListTasksByDate(date time.Time) ([]*models.Task, error)
	ListTasksByStatus(status models.TaskStatus) ([]*models.Task, error)

	// Schedule operations
	SaveSchedule(schedule *models.Schedule) error
	GetSchedule() (*models.Schedule, error)

	// Statistics operations
	GetDailyStats(date time.Time) (*models.DailyStats, error)
	GetWeeklyStats(startDate time.Time) (*models.WeeklyStats, error)
	SaveDailyStats(stats *models.DailyStats) error

	// Utility operations
	Close() error
	Backup() error
}

// FileStorage implements Storage interface using file-based persistence
type FileStorage struct {
	dataDir   string
	tasksFile string
	schedFile string
	statsDir  string
	mu        sync.RWMutex
}

// NewFileStorage creates a new file-based storage instance
func NewFileStorage(dataDir string) (*FileStorage, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create stats directory
	statsDir := filepath.Join(dataDir, "stats")
	if err := os.MkdirAll(statsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create stats directory: %w", err)
	}

	fs := &FileStorage{
		dataDir:   dataDir,
		tasksFile: filepath.Join(dataDir, "tasks.json"),
		schedFile: filepath.Join(dataDir, "schedule.json"),
		statsDir:  statsDir,
	}

	// Initialize files if they don't exist
	if err := fs.initFiles(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage files: %w", err)
	}

	return fs, nil
}

// initFiles creates initial storage files if they don't exist
func (fs *FileStorage) initFiles() error {
	// Initialize tasks file
	if _, err := os.Stat(fs.tasksFile); os.IsNotExist(err) {
		tasks := make([]*models.Task, 0)
		if err := fs.writeJSON(fs.tasksFile, tasks); err != nil {
			return fmt.Errorf("failed to initialize tasks file: %w", err)
		}
	}

	// Initialize schedule file with default schedule
	if _, err := os.Stat(fs.schedFile); os.IsNotExist(err) {
		schedule := models.GetDefaultSchedule()
		if err := fs.writeJSON(fs.schedFile, schedule); err != nil {
			return fmt.Errorf("failed to initialize schedule file: %w", err)
		}
	}

	return nil
}

// writeJSON writes data to a JSON file
func (fs *FileStorage) writeJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// readJSON reads data from a JSON file
func (fs *FileStorage) readJSON(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(data)
}

// loadTasks loads all tasks from storage
func (fs *FileStorage) loadTasks() ([]*models.Task, error) {
	var tasks []*models.Task
	if err := fs.readJSON(fs.tasksFile, &tasks); err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}
	return tasks, nil
}

// saveTasks saves all tasks to storage
func (fs *FileStorage) saveTasks(tasks []*models.Task) error {
	return fs.writeJSON(fs.tasksFile, tasks)
}

// CreateTask creates a new task
func (fs *FileStorage) CreateTask(task *models.Task) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	tasks, err := fs.loadTasks()
	if err != nil {
		return err
	}

	// Check for duplicate ID
	for _, existing := range tasks {
		if existing.ID == task.ID {
			return fmt.Errorf("task with ID %s already exists", task.ID)
		}
	}

	tasks = append(tasks, task)
	return fs.saveTasks(tasks)
}

// GetTask retrieves a task by ID
func (fs *FileStorage) GetTask(id string) (*models.Task, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	tasks, err := fs.loadTasks()
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, fmt.Errorf("task with ID %s not found", id)
}

// UpdateTask updates an existing task
func (fs *FileStorage) UpdateTask(task *models.Task) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	tasks, err := fs.loadTasks()
	if err != nil {
		return err
	}

	for i, existing := range tasks {
		if existing.ID == task.ID {
			task.UpdatedAt = time.Now()
			tasks[i] = task
			return fs.saveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %s not found", task.ID)
}

// DeleteTask deletes a task by ID
func (fs *FileStorage) DeleteTask(id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	tasks, err := fs.loadTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return fs.saveTasks(tasks)
		}
	}

	return fmt.Errorf("task with ID %s not found", id)
}

// ListTasks returns all tasks
func (fs *FileStorage) ListTasks() ([]*models.Task, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	return fs.loadTasks()
}

// ListTasksByDate returns tasks for a specific date
func (fs *FileStorage) ListTasksByDate(date time.Time) ([]*models.Task, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	tasks, err := fs.loadTasks()
	if err != nil {
		return nil, err
	}

	var result []*models.Task
	dateStr := date.Format("2006-01-02")

	for _, task := range tasks {
		// Check if task was created on the date
		if task.CreatedAt.Format("2006-01-02") == dateStr {
			result = append(result, task)
			continue
		}

		// Check if task was scheduled for the date
		if task.ScheduledTime != nil && task.ScheduledTime.Format("2006-01-02") == dateStr {
			result = append(result, task)
			continue
		}

		// Check if task was completed on the date
		if task.CompletedAt != nil && task.CompletedAt.Format("2006-01-02") == dateStr {
			result = append(result, task)
			continue
		}
	}

	// Sort by created time
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}

// ListTasksByStatus returns tasks with a specific status
func (fs *FileStorage) ListTasksByStatus(status models.TaskStatus) ([]*models.Task, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	tasks, err := fs.loadTasks()
	if err != nil {
		return nil, err
	}

	var result []*models.Task
	for _, task := range tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}

	return result, nil
}

// SaveSchedule saves the schedule configuration
func (fs *FileStorage) SaveSchedule(schedule *models.Schedule) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	return fs.writeJSON(fs.schedFile, schedule)
}

// GetSchedule returns the current schedule configuration
func (fs *FileStorage) GetSchedule() (*models.Schedule, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var schedule models.Schedule
	if err := fs.readJSON(fs.schedFile, &schedule); err != nil {
		return nil, fmt.Errorf("failed to load schedule: %w", err)
	}

	return &schedule, nil
}

// GetDailyStats returns statistics for a specific date
func (fs *FileStorage) GetDailyStats(date time.Time) (*models.DailyStats, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	statsFile := filepath.Join(fs.statsDir, fmt.Sprintf("%s.json", date.Format("2006-01-02")))

	var stats models.DailyStats
	if err := fs.readJSON(statsFile, &stats); err != nil {
		if os.IsNotExist(err) {
			// Return empty stats if file doesn't exist
			return &models.DailyStats{
				Date:            date,
				HourlyBreakdown: make(map[string]models.Score),
			}, nil
		}
		return nil, fmt.Errorf("failed to load daily stats: %w", err)
	}

	return &stats, nil
}

// SaveDailyStats saves statistics for a specific date
func (fs *FileStorage) SaveDailyStats(stats *models.DailyStats) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	statsFile := filepath.Join(fs.statsDir, fmt.Sprintf("%s.json", stats.Date.Format("2006-01-02")))
	return fs.writeJSON(statsFile, stats)
}

// GetWeeklyStats returns aggregated statistics for a week
func (fs *FileStorage) GetWeeklyStats(startDate time.Time) (*models.WeeklyStats, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	// Calculate week boundaries
	endDate := startDate.AddDate(0, 0, 6)

	var dailyStats []models.DailyStats
	var weeklyTotal models.Score
	var totalTime time.Duration

	// Load daily stats for each day of the week
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dayStats, err := fs.GetDailyStats(d)
		if err != nil {
			return nil, fmt.Errorf("failed to load stats for %s: %w", d.Format("2006-01-02"), err)
		}

		dailyStats = append(dailyStats, *dayStats)
		weeklyTotal.Work += dayStats.TotalScore.Work
		weeklyTotal.Play += dayStats.TotalScore.Play
		weeklyTotal.Learn += dayStats.TotalScore.Learn
		totalTime += dayStats.TimeSpent
	}

	// Calculate weekly average
	weeklyAverage := models.Score{
		Work:  weeklyTotal.Work / 7,
		Play:  weeklyTotal.Play / 7,
		Learn: weeklyTotal.Learn / 7,
	}

	return &models.WeeklyStats{
		StartDate:      startDate,
		EndDate:        endDate,
		DailyStats:     dailyStats,
		WeeklyTotal:    weeklyTotal,
		WeeklyAverage:  weeklyAverage,
		TotalTimeSpent: totalTime,
	}, nil
}

// Backup creates a backup of all data
func (fs *FileStorage) Backup() error {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	backupDir := filepath.Join(fs.dataDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")

	// For now, just copy the JSON files
	// In a more sophisticated implementation, we'd create a proper archive
	backupDataDir := filepath.Join(backupDir, fmt.Sprintf("data_%s", timestamp))
	if err := os.MkdirAll(backupDataDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup data directory: %w", err)
	}

	// Copy tasks file
	if err := fs.copyFile(fs.tasksFile, filepath.Join(backupDataDir, "tasks.json")); err != nil {
		return fmt.Errorf("failed to backup tasks: %w", err)
	}

	// Copy schedule file
	if err := fs.copyFile(fs.schedFile, filepath.Join(backupDataDir, "schedule.json")); err != nil {
		return fmt.Errorf("failed to backup schedule: %w", err)
	}

	// Copy stats directory
	if err := fs.copyDir(fs.statsDir, filepath.Join(backupDataDir, "stats")); err != nil {
		return fmt.Errorf("failed to backup stats: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func (fs *FileStorage) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// copyDir copies a directory from src to dst
func (fs *FileStorage) copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := fs.copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := fs.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// Close closes the storage (no-op for file storage)
func (fs *FileStorage) Close() error {
	return nil
}

// GetDataDir returns the data directory path
func (fs *FileStorage) GetDataDir() string {
	return fs.dataDir
}
