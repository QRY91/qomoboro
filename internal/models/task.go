package models

import (
	"time"
)

// Score represents the triple-metric scoring system
type Score struct {
	Work  int `json:"work" yaml:"work"`   // Business/professional productivity (0-5)
	Play  int `json:"play" yaml:"play"`   // Recreation/leisure/enjoyment (0-5)
	Learn int `json:"learn" yaml:"learn"` // Educational/skill development (0-5)
}

// IsValid checks if all scores are within valid range (0-5)
func (s Score) IsValid() bool {
	return s.Work >= 0 && s.Work <= 5 &&
		s.Play >= 0 && s.Play <= 5 &&
		s.Learn >= 0 && s.Learn <= 5
}

// Total returns the sum of all scores
func (s Score) Total() int {
	return s.Work + s.Play + s.Learn
}

// Average returns the average score across all metrics
func (s Score) Average() float64 {
	return float64(s.Total()) / 3.0
}

// TaskStatus represents the current state of a task
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusActive
	TaskStatusCompleted
	TaskStatusCancelled
	TaskStatusPaused
)

// String returns the string representation of TaskStatus
func (ts TaskStatus) String() string {
	switch ts {
	case TaskStatusPending:
		return "pending"
	case TaskStatusActive:
		return "active"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusCancelled:
		return "cancelled"
	case TaskStatusPaused:
		return "paused"
	default:
		return "unknown"
	}
}

// Task represents a single activity or work item
type Task struct {
	ID          string     `json:"id" yaml:"id"`
	Title       string     `json:"title" yaml:"title"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	Score       Score      `json:"score" yaml:"score"`
	Status      TaskStatus `json:"status" yaml:"status"`

	// Timing information
	EstimatedDuration time.Duration `json:"estimated_duration" yaml:"estimated_duration"`
	ActualDuration    time.Duration `json:"actual_duration,omitempty" yaml:"actual_duration,omitempty"`
	StartTime         *time.Time    `json:"start_time,omitempty" yaml:"start_time,omitempty"`
	EndTime           *time.Time    `json:"end_time,omitempty" yaml:"end_time,omitempty"`

	// Scheduling
	ScheduledTime *time.Time `json:"scheduled_time,omitempty" yaml:"scheduled_time,omitempty"`
	CanonicalHour string     `json:"canonical_hour,omitempty" yaml:"canonical_hour,omitempty"`

	// Metadata
	Tags        []string   `json:"tags,omitempty" yaml:"tags,omitempty"`
	CreatedAt   time.Time  `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" yaml:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`

	// Notes and reflection
	Notes      string `json:"notes,omitempty" yaml:"notes,omitempty"`
	Reflection string `json:"reflection,omitempty" yaml:"reflection,omitempty"`
}

// IsActive returns true if the task is currently being worked on
func (t *Task) IsActive() bool {
	return t.Status == TaskStatusActive
}

// IsCompleted returns true if the task has been completed
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusCompleted
}

// Start begins work on the task
func (t *Task) Start() {
	now := time.Now()
	t.StartTime = &now
	t.Status = TaskStatusActive
	t.UpdatedAt = now
}

// Complete finishes the task
func (t *Task) Complete() {
	now := time.Now()
	t.EndTime = &now
	t.CompletedAt = &now
	t.Status = TaskStatusCompleted
	t.UpdatedAt = now

	if t.StartTime != nil {
		t.ActualDuration = now.Sub(*t.StartTime)
	}
}

// Pause temporarily stops work on the task
func (t *Task) Pause() {
	if t.Status == TaskStatusActive && t.StartTime != nil {
		now := time.Now()
		t.ActualDuration += now.Sub(*t.StartTime)
		t.StartTime = nil
		t.Status = TaskStatusPaused
		t.UpdatedAt = now
	}
}

// Resume continues work on a paused task
func (t *Task) Resume() {
	if t.Status == TaskStatusPaused {
		now := time.Now()
		t.StartTime = &now
		t.Status = TaskStatusActive
		t.UpdatedAt = now
	}
}

// CanonicalHour represents a traditional canonical hour time block
type CanonicalHour struct {
	Name        string        `json:"name" yaml:"name"`
	StartTime   string        `json:"start_time" yaml:"start_time"` // Format: "15:04"
	EndTime     string        `json:"end_time" yaml:"end_time"`     // Format: "15:04"
	Duration    time.Duration `json:"duration" yaml:"duration"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Purpose     string        `json:"purpose,omitempty" yaml:"purpose,omitempty"`

	// Default scoring suggestions for this hour
	DefaultScore Score `json:"default_score,omitempty" yaml:"default_score,omitempty"`
}

// IsActive checks if the given time falls within this canonical hour
func (ch *CanonicalHour) IsActive(t time.Time) bool {
	timeStr := t.Format("15:04")
	return timeStr >= ch.StartTime && timeStr < ch.EndTime
}

// Schedule represents a collection of canonical hours
type Schedule struct {
	Name  string          `json:"name" yaml:"name"`
	Hours []CanonicalHour `json:"hours" yaml:"hours"`
}

// GetCurrentHour returns the canonical hour for the given time
func (s *Schedule) GetCurrentHour(t time.Time) *CanonicalHour {
	for _, hour := range s.Hours {
		if hour.IsActive(t) {
			return &hour
		}
	}
	return nil
}

// GetHourByName returns the canonical hour with the given name
func (s *Schedule) GetHourByName(name string) *CanonicalHour {
	for _, hour := range s.Hours {
		if hour.Name == name {
			return &hour
		}
	}
	return nil
}

// DailyStats represents aggregated statistics for a day
type DailyStats struct {
	Date           time.Time     `json:"date" yaml:"date"`
	TotalTasks     int           `json:"total_tasks" yaml:"total_tasks"`
	CompletedTasks int           `json:"completed_tasks" yaml:"completed_tasks"`
	TotalScore     Score         `json:"total_score" yaml:"total_score"`
	AverageScore   Score         `json:"average_score" yaml:"average_score"`
	TimeSpent      time.Duration `json:"time_spent" yaml:"time_spent"`

	// Breakdown by canonical hour
	HourlyBreakdown map[string]Score `json:"hourly_breakdown,omitempty" yaml:"hourly_breakdown,omitempty"`
}

// CompletionRate returns the percentage of tasks completed
func (ds *DailyStats) CompletionRate() float64 {
	if ds.TotalTasks == 0 {
		return 0.0
	}
	return float64(ds.CompletedTasks) / float64(ds.TotalTasks) * 100.0
}

// WeeklyStats represents aggregated statistics for a week
type WeeklyStats struct {
	StartDate      time.Time     `json:"start_date" yaml:"start_date"`
	EndDate        time.Time     `json:"end_date" yaml:"end_date"`
	DailyStats     []DailyStats  `json:"daily_stats" yaml:"daily_stats"`
	WeeklyTotal    Score         `json:"weekly_total" yaml:"weekly_total"`
	WeeklyAverage  Score         `json:"weekly_average" yaml:"weekly_average"`
	TotalTimeSpent time.Duration `json:"total_time_spent" yaml:"total_time_spent"`
}

// GetDefaultSchedule returns the standard canonical hours schedule
func GetDefaultSchedule() Schedule {
	return Schedule{
		Name: "Traditional Canonical Hours",
		Hours: []CanonicalHour{
			{
				Name:         "Matins",
				StartTime:    "06:00",
				EndTime:      "07:30",
				Duration:     time.Hour + 30*time.Minute,
				Description:  "Deep work, planning, and preparation",
				Purpose:      "High-focus work when mind is fresh",
				DefaultScore: Score{Work: 4, Play: 1, Learn: 3},
			},
			{
				Name:         "Lauds",
				StartTime:    "07:30",
				EndTime:      "09:00",
				Duration:     time.Hour + 30*time.Minute,
				Description:  "Administrative tasks and organization",
				Purpose:      "Handle communications and planning",
				DefaultScore: Score{Work: 3, Play: 1, Learn: 2},
			},
			{
				Name:         "Prime",
				StartTime:    "09:00",
				EndTime:      "12:00",
				Duration:     3 * time.Hour,
				Description:  "Primary work blocks and major tasks",
				Purpose:      "Core productive work period",
				DefaultScore: Score{Work: 5, Play: 1, Learn: 3},
			},
			{
				Name:         "Terce",
				StartTime:    "12:00",
				EndTime:      "13:30",
				Duration:     time.Hour + 30*time.Minute,
				Description:  "Meetings, collaboration, and communication",
				Purpose:      "Social and collaborative work",
				DefaultScore: Score{Work: 3, Play: 2, Learn: 2},
			},
			{
				Name:         "Sext",
				StartTime:    "13:30",
				EndTime:      "15:00",
				Duration:     time.Hour + 30*time.Minute,
				Description:  "Lunch, recovery, and personal time",
				Purpose:      "Rest and recharge",
				DefaultScore: Score{Work: 1, Play: 4, Learn: 1},
			},
			{
				Name:         "None",
				StartTime:    "15:00",
				EndTime:      "16:30",
				Duration:     time.Hour + 30*time.Minute,
				Description:  "Creative work and experimentation",
				Purpose:      "Innovation and creative problem-solving",
				DefaultScore: Score{Work: 3, Play: 3, Learn: 4},
			},
			{
				Name:         "Vespers",
				StartTime:    "16:30",
				EndTime:      "18:00",
				Duration:     time.Hour + 30*time.Minute,
				Description:  "Learning, documentation, and skill development",
				Purpose:      "Knowledge acquisition and sharing",
				DefaultScore: Score{Work: 2, Play: 2, Learn: 5},
			},
			{
				Name:         "Compline",
				StartTime:    "18:00",
				EndTime:      "20:00",
				Duration:     2 * time.Hour,
				Description:  "Planning, reflection, and personal projects",
				Purpose:      "Review and prepare for tomorrow",
				DefaultScore: Score{Work: 2, Play: 3, Learn: 3},
			},
		},
	}
}
