package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"qomoboro/internal/models"
	"qomoboro/internal/storage"
)

// ViewMode represents different application views
type ViewMode int

const (
	ViewModeMain ViewMode = iota
	ViewModeTaskList
	ViewModeSchedule
	ViewModeStats
	ViewModeCreateTask
	ViewModeTaskDetail
	ViewModeSettings
)

// Colors and styles
var (
	primaryColor   = lipgloss.AdaptiveColor{Light: "#688060", Dark: "#688060"}
	secondaryColor = lipgloss.AdaptiveColor{Light: "#8CD0D3", Dark: "#8CD0D3"}
	accentColor    = lipgloss.AdaptiveColor{Light: "#F0DFAF", Dark: "#F0DFAF"}
	errorColor     = lipgloss.AdaptiveColor{Light: "#DCA3A3", Dark: "#DCA3A3"}
	mutedColor     = lipgloss.AdaptiveColor{Light: "#7F7F7F", Dark: "#7F7F7F"}
)

// Styles holds all the lipgloss styles
type Styles struct {
	Base,
	Header,
	Title,
	Subtitle,
	Status,
	StatusActive,
	StatusCompleted,
	StatusPending,
	Score,
	ScoreWork,
	ScorePlay,
	ScoreLearn,
	Help,
	Error,
	Border,
	Highlight,
	Muted lipgloss.Style
}

// NewStyles creates a new set of styles
func NewStyles() *Styles {
	s := &Styles{}

	s.Base = lipgloss.NewStyle().
		Padding(1, 2)

	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(primaryColor).
		PaddingBottom(1).
		MarginBottom(1)

	s.Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(secondaryColor)

	s.Status = lipgloss.NewStyle().
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor)

	s.StatusActive = s.Status.Copy().
		Foreground(accentColor).
		BorderForeground(accentColor)

	s.StatusCompleted = s.Status.Copy().
		Foreground(primaryColor).
		BorderForeground(primaryColor)

	s.StatusPending = s.Status.Copy().
		Foreground(mutedColor).
		BorderForeground(mutedColor)

	s.Score = lipgloss.NewStyle().
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder())

	s.ScoreWork = s.Score.Copy().
		Foreground(primaryColor).
		BorderForeground(primaryColor)

	s.ScorePlay = s.Score.Copy().
		Foreground(accentColor).
		BorderForeground(accentColor)

	s.ScoreLearn = s.Score.Copy().
		Foreground(secondaryColor).
		BorderForeground(secondaryColor)

	s.Help = lipgloss.NewStyle().
		Foreground(mutedColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(mutedColor).
		PaddingTop(1).
		MarginTop(1)

	s.Error = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)

	s.Border = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor)

	s.Highlight = lipgloss.NewStyle().
		Background(primaryColor).
		Foreground(lipgloss.Color("0"))

	s.Muted = lipgloss.NewStyle().
		Foreground(mutedColor)

	return s
}

// App represents the main TUI application
type App struct {
	storage         storage.Storage
	styles          *Styles
	currentView     ViewMode
	currentTask     *models.Task
	currentSchedule *models.Schedule
	tasks           []*models.Task
	selectedIndex   int
	form            *huh.Form
	width           int
	height          int
	error           error
	message         string
	quitting        bool
}

// NewApp creates a new TUI application
func NewApp(storage storage.Storage) *App {
	app := &App{
		storage:     storage,
		styles:      NewStyles(),
		currentView: ViewModeMain,
		width:       80,
		height:      24,
	}

	// Load initial data
	app.loadData()

	return app
}

// loadData loads initial data from storage
func (a *App) loadData() {
	var err error

	// Load tasks
	a.tasks, err = a.storage.ListTasks()
	if err != nil {
		a.error = fmt.Errorf("failed to load tasks: %w", err)
		return
	}

	// Load schedule
	a.currentSchedule, err = a.storage.GetSchedule()
	if err != nil {
		a.error = fmt.Errorf("failed to load schedule: %w", err)
		return
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		switch a.currentView {
		case ViewModeMain:
			return a.updateMain(msg)
		case ViewModeTaskList:
			return a.updateTaskList(msg)
		case ViewModeSchedule:
			return a.updateSchedule(msg)
		case ViewModeStats:
			return a.updateStats(msg)
		case ViewModeCreateTask:
			return a.updateCreateTask(msg)
		case ViewModeTaskDetail:
			return a.updateTaskDetail(msg)
		case ViewModeSettings:
			return a.updateSettings(msg)
		}

	case tea.QuitMsg:
		a.quitting = true
		return a, tea.Quit
	}

	return a, nil
}

// updateMain handles the main menu navigation
func (a *App) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return a, tea.Quit
	case "t":
		a.currentView = ViewModeTaskList
		a.selectedIndex = 0
	case "s":
		a.currentView = ViewModeSchedule
	case "d":
		a.currentView = ViewModeStats
	case "c":
		a.currentView = ViewModeCreateTask
		a.initCreateTaskForm()
	case "g":
		a.currentView = ViewModeSettings
	case "r":
		a.loadData()
		a.message = "Data reloaded"
	}
	return a, nil
}

// updateTaskList handles task list navigation
func (a *App) updateTaskList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.currentView = ViewModeMain
	case "up", "k":
		if a.selectedIndex > 0 {
			a.selectedIndex--
		}
	case "down", "j":
		if a.selectedIndex < len(a.tasks)-1 {
			a.selectedIndex++
		}
	case "enter":
		if a.selectedIndex < len(a.tasks) {
			a.currentTask = a.tasks[a.selectedIndex]
			a.currentView = ViewModeTaskDetail
		}
	case "c":
		a.currentView = ViewModeCreateTask
		a.initCreateTaskForm()
	case "d":
		if a.selectedIndex < len(a.tasks) {
			task := a.tasks[a.selectedIndex]
			if err := a.storage.DeleteTask(task.ID); err != nil {
				a.error = err
			} else {
				a.loadData()
				a.message = "Task deleted"
				if a.selectedIndex >= len(a.tasks) {
					a.selectedIndex = len(a.tasks) - 1
				}
				if a.selectedIndex < 0 {
					a.selectedIndex = 0
				}
			}
		}
	case " ":
		if a.selectedIndex < len(a.tasks) {
			task := a.tasks[a.selectedIndex]
			if task.Status == models.TaskStatusCompleted {
				task.Status = models.TaskStatusPending
			} else {
				task.Complete()
			}
			if err := a.storage.UpdateTask(task); err != nil {
				a.error = err
			} else {
				a.loadData()
				a.message = "Task updated"
			}
		}
	}
	return a, nil
}

// updateSchedule handles schedule view navigation
func (a *App) updateSchedule(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.currentView = ViewModeMain
	}
	return a, nil
}

// updateStats handles statistics view navigation
func (a *App) updateStats(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.currentView = ViewModeMain
	}
	return a, nil
}

// updateCreateTask handles task creation form
func (a *App) updateCreateTask(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.currentView = ViewModeTaskList
		return a, nil
	}

	if a.form != nil {
		form, cmd := a.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			a.form = f

			if a.form.State == huh.StateCompleted {
				a.createTaskFromForm()
				a.currentView = ViewModeTaskList
			}
		}
		return a, cmd
	}

	return a, nil
}

// updateTaskDetail handles task detail view
func (a *App) updateTaskDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.currentView = ViewModeTaskList
	case " ":
		if a.currentTask != nil {
			if a.currentTask.Status == models.TaskStatusCompleted {
				a.currentTask.Status = models.TaskStatusPending
			} else {
				a.currentTask.Complete()
			}
			if err := a.storage.UpdateTask(a.currentTask); err != nil {
				a.error = err
			} else {
				a.loadData()
				a.message = "Task updated"
			}
		}
	case "d":
		if a.currentTask != nil {
			if err := a.storage.DeleteTask(a.currentTask.ID); err != nil {
				a.error = err
			} else {
				a.loadData()
				a.message = "Task deleted"
				a.currentView = ViewModeTaskList
			}
		}
	}
	return a, nil
}

// updateSettings handles settings view
func (a *App) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		a.currentView = ViewModeMain
	}
	return a, nil
}

// View renders the current view
func (a *App) View() string {
	if a.quitting {
		return ""
	}

	var content string

	switch a.currentView {
	case ViewModeMain:
		content = a.viewMain()
	case ViewModeTaskList:
		content = a.viewTaskList()
	case ViewModeSchedule:
		content = a.viewSchedule()
	case ViewModeStats:
		content = a.viewStats()
	case ViewModeCreateTask:
		content = a.viewCreateTask()
	case ViewModeTaskDetail:
		content = a.viewTaskDetail()
	case ViewModeSettings:
		content = a.viewSettings()
	}

	// Add error message if present
	if a.error != nil {
		content += "\n\n" + a.styles.Error.Render(fmt.Sprintf("Error: %v", a.error))
		a.error = nil // Clear error after displaying
	}

	// Add success message if present
	if a.message != "" {
		content += "\n\n" + a.styles.Subtitle.Render(a.message)
		a.message = "" // Clear message after displaying
	}

	return a.styles.Base.Render(content)
}

// viewMain renders the main menu
func (a *App) viewMain() string {
	title := a.styles.Header.Render("QOMOBORO - Canonical Hours Task Manager")

	// Get current canonical hour
	currentHour := ""
	if a.currentSchedule != nil {
		if hour := a.currentSchedule.GetCurrentHour(time.Now()); hour != nil {
			currentHour = fmt.Sprintf("Current Hour: %s (%s)", hour.Name, hour.Description)
		}
	}

	// Task summary
	taskSummary := a.getTaskSummary()

	menu := `
Navigation:
  [t] Tasks       - Manage your tasks
  [s] Schedule    - View canonical hours
  [d] Statistics  - View productivity stats
  [c] Create      - Add new task
  [g] Settings    - Configure app
  [r] Refresh     - Reload data
  [q] Quit        - Exit application
`

	content := []string{
		title,
		"",
		a.styles.Subtitle.Render(currentHour),
		"",
		taskSummary,
		"",
		menu,
	}

	return strings.Join(content, "\n")
}

// getTaskSummary returns a summary of current tasks
func (a *App) getTaskSummary() string {
	if len(a.tasks) == 0 {
		return a.styles.Muted.Render("No tasks yet. Press 'c' to create your first task!")
	}

	var pending, active, completed int
	for _, task := range a.tasks {
		switch task.Status {
		case models.TaskStatusPending:
			pending++
		case models.TaskStatusActive:
			active++
		case models.TaskStatusCompleted:
			completed++
		}
	}

	return fmt.Sprintf("Tasks: %s pending, %s active, %s completed",
		a.styles.StatusPending.Render(fmt.Sprintf("%d", pending)),
		a.styles.StatusActive.Render(fmt.Sprintf("%d", active)),
		a.styles.StatusCompleted.Render(fmt.Sprintf("%d", completed)))
}

// viewTaskList renders the task list
func (a *App) viewTaskList() string {
	title := a.styles.Header.Render("Task List")

	if len(a.tasks) == 0 {
		return title + "\n\n" + a.styles.Muted.Render("No tasks yet. Press 'c' to create a task.")
	}

	var taskList []string
	for i, task := range a.tasks {
		style := a.styles.Base
		if i == a.selectedIndex {
			style = a.styles.Highlight
		}

		statusStyle := a.styles.StatusPending
		switch task.Status {
		case models.TaskStatusActive:
			statusStyle = a.styles.StatusActive
		case models.TaskStatusCompleted:
			statusStyle = a.styles.StatusCompleted
		}

		scores := fmt.Sprintf("W:%d P:%d L:%d", task.Score.Work, task.Score.Play, task.Score.Learn)

		line := fmt.Sprintf("%s %s %s",
			statusStyle.Render(fmt.Sprintf("[%s]", task.Status.String()[:1])),
			task.Title,
			a.styles.Muted.Render(scores))

		taskList = append(taskList, style.Render(line))
	}

	help := a.styles.Help.Render("↑/↓: navigate, Enter: details, Space: toggle, c: create, d: delete, q: back")

	return title + "\n\n" + strings.Join(taskList, "\n") + "\n\n" + help
}

// viewSchedule renders the schedule view
func (a *App) viewSchedule() string {
	title := a.styles.Header.Render("Canonical Hours Schedule")

	if a.currentSchedule == nil {
		return title + "\n\n" + a.styles.Error.Render("No schedule loaded")
	}

	var scheduleLines []string
	now := time.Now()

	for _, hour := range a.currentSchedule.Hours {
		style := a.styles.Base
		if hour.IsActive(now) {
			style = a.styles.Highlight
		}

		line := fmt.Sprintf("%s - %s: %s",
			hour.StartTime,
			hour.EndTime,
			hour.Name)

		if hour.Description != "" {
			line += fmt.Sprintf(" (%s)", hour.Description)
		}

		scheduleLines = append(scheduleLines, style.Render(line))
	}

	help := a.styles.Help.Render("q: back to main menu")

	return title + "\n\n" + strings.Join(scheduleLines, "\n") + "\n\n" + help
}

// viewStats renders the statistics view
func (a *App) viewStats() string {
	title := a.styles.Header.Render("Statistics")

	// Get today's stats
	today := time.Now()
	stats, err := a.storage.GetDailyStats(today)
	if err != nil {
		return title + "\n\n" + a.styles.Error.Render(fmt.Sprintf("Failed to load stats: %v", err))
	}

	content := []string{
		title,
		"",
		fmt.Sprintf("Today (%s):", today.Format("2006-01-02")),
		fmt.Sprintf("  Tasks: %d total, %d completed", stats.TotalTasks, stats.CompletedTasks),
		fmt.Sprintf("  Scores: Work %d, Play %d, Learn %d", stats.TotalScore.Work, stats.TotalScore.Play, stats.TotalScore.Learn),
		fmt.Sprintf("  Time: %s", stats.TimeSpent.String()),
		"",
		a.styles.Help.Render("q: back to main menu"),
	}

	return strings.Join(content, "\n")
}

// viewCreateTask renders the task creation form
func (a *App) viewCreateTask() string {
	title := a.styles.Header.Render("Create New Task")

	if a.form == nil {
		return title + "\n\n" + a.styles.Error.Render("Form not initialized")
	}

	return title + "\n\n" + a.form.View()
}

// viewTaskDetail renders the task detail view
func (a *App) viewTaskDetail() string {
	title := a.styles.Header.Render("Task Details")

	if a.currentTask == nil {
		return title + "\n\n" + a.styles.Error.Render("No task selected")
	}

	task := a.currentTask

	statusStyle := a.styles.StatusPending
	switch task.Status {
	case models.TaskStatusActive:
		statusStyle = a.styles.StatusActive
	case models.TaskStatusCompleted:
		statusStyle = a.styles.StatusCompleted
	}

	content := []string{
		title,
		"",
		a.styles.Title.Render(task.Title),
		statusStyle.Render(fmt.Sprintf("Status: %s", task.Status.String())),
		"",
		fmt.Sprintf("Scores: %s %s %s",
			a.styles.ScoreWork.Render(fmt.Sprintf("Work: %d", task.Score.Work)),
			a.styles.ScorePlay.Render(fmt.Sprintf("Play: %d", task.Score.Play)),
			a.styles.ScoreLearn.Render(fmt.Sprintf("Learn: %d", task.Score.Learn))),
		"",
		fmt.Sprintf("Created: %s", task.CreatedAt.Format("2006-01-02 15:04")),
		fmt.Sprintf("Updated: %s", task.UpdatedAt.Format("2006-01-02 15:04")),
	}

	if task.Description != "" {
		content = append(content, "", "Description:", task.Description)
	}

	if task.Notes != "" {
		content = append(content, "", "Notes:", task.Notes)
	}

	content = append(content, "", a.styles.Help.Render("Space: toggle status, d: delete, q: back"))

	return strings.Join(content, "\n")
}

// viewSettings renders the settings view
func (a *App) viewSettings() string {
	title := a.styles.Header.Render("Settings")

	content := []string{
		title,
		"",
		"Settings panel coming soon...",
		"",
		a.styles.Help.Render("q: back to main menu"),
	}

	return strings.Join(content, "\n")
}

// initCreateTaskForm initializes the task creation form
func (a *App) initCreateTaskForm() {
	a.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("title").
				Title("Task Title").
				Description("What needs to be done?").
				Placeholder("Enter task title..."),

			huh.NewText().
				Key("description").
				Title("Description").
				Description("Additional details (optional)").
				Placeholder("Enter description..."),

			huh.NewSelect[int]().
				Key("work").
				Title("Work Score (0-5)").
				Description("How much does this contribute to work/business?").
				Options(
					huh.NewOption("0", 0),
					huh.NewOption("1", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3", 3),
					huh.NewOption("4", 4),
					huh.NewOption("5", 5),
				).
				Value(new(int)),

			huh.NewSelect[int]().
				Key("play").
				Title("Play Score (0-5)").
				Description("How enjoyable/recreational is this?").
				Options(
					huh.NewOption("0", 0),
					huh.NewOption("1", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3", 3),
					huh.NewOption("4", 4),
					huh.NewOption("5", 5),
				).
				Value(new(int)),

			huh.NewSelect[int]().
				Key("learn").
				Title("Learn Score (0-5)").
				Description("How much will you learn from this?").
				Options(
					huh.NewOption("0", 0),
					huh.NewOption("1", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3", 3),
					huh.NewOption("4", 4),
					huh.NewOption("5", 5),
				).
				Value(new(int)),
		),
	)
}

// createTaskFromForm creates a task from the form data
func (a *App) createTaskFromForm() {
	if a.form == nil {
		return
	}

	task := &models.Task{
		ID:          fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Title:       a.form.GetString("title"),
		Description: a.form.GetString("description"),
		Score: models.Score{
			Work:  a.form.Get("work").(int),
			Play:  a.form.Get("play").(int),
			Learn: a.form.Get("learn").(int),
		},
		Status:    models.TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.storage.CreateTask(task); err != nil {
		a.error = err
	} else {
		a.loadData()
		a.message = "Task created successfully"
	}
}

// Run starts the TUI application
func (a *App) Run() error {
	p := tea.NewProgram(a, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
