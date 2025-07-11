package models

import (
	"testing"
	"time"
)

func TestScore_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		score Score
		want  bool
	}{
		{
			name:  "valid scores",
			score: Score{Work: 3, Play: 4, Learn: 2},
			want:  true,
		},
		{
			name:  "minimum valid scores",
			score: Score{Work: 0, Play: 0, Learn: 0},
			want:  true,
		},
		{
			name:  "maximum valid scores",
			score: Score{Work: 5, Play: 5, Learn: 5},
			want:  true,
		},
		{
			name:  "work score too high",
			score: Score{Work: 6, Play: 3, Learn: 2},
			want:  false,
		},
		{
			name:  "play score negative",
			score: Score{Work: 3, Play: -1, Learn: 2},
			want:  false,
		},
		{
			name:  "learn score too high",
			score: Score{Work: 3, Play: 4, Learn: 10},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.score.IsValid(); got != tt.want {
				t.Errorf("Score.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScore_Total(t *testing.T) {
	tests := []struct {
		name  string
		score Score
		want  int
	}{
		{
			name:  "normal scores",
			score: Score{Work: 3, Play: 4, Learn: 2},
			want:  9,
		},
		{
			name:  "zero scores",
			score: Score{Work: 0, Play: 0, Learn: 0},
			want:  0,
		},
		{
			name:  "maximum scores",
			score: Score{Work: 5, Play: 5, Learn: 5},
			want:  15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.score.Total(); got != tt.want {
				t.Errorf("Score.Total() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScore_Average(t *testing.T) {
	tests := []struct {
		name  string
		score Score
		want  float64
	}{
		{
			name:  "normal scores",
			score: Score{Work: 3, Play: 6, Learn: 0},
			want:  3.0,
		},
		{
			name:  "zero scores",
			score: Score{Work: 0, Play: 0, Learn: 0},
			want:  0.0,
		},
		{
			name:  "maximum scores",
			score: Score{Work: 5, Play: 5, Learn: 5},
			want:  5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.score.Average(); got != tt.want {
				t.Errorf("Score.Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status TaskStatus
		want   string
	}{
		{
			name:   "pending",
			status: TaskStatusPending,
			want:   "pending",
		},
		{
			name:   "active",
			status: TaskStatusActive,
			want:   "active",
		},
		{
			name:   "completed",
			status: TaskStatusCompleted,
			want:   "completed",
		},
		{
			name:   "cancelled",
			status: TaskStatusCancelled,
			want:   "cancelled",
		},
		{
			name:   "paused",
			status: TaskStatusPaused,
			want:   "paused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("TaskStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_IsActive(t *testing.T) {
	task := &Task{Status: TaskStatusActive}
	if !task.IsActive() {
		t.Errorf("Task.IsActive() = false, want true for active task")
	}

	task.Status = TaskStatusPending
	if task.IsActive() {
		t.Errorf("Task.IsActive() = true, want false for pending task")
	}
}

func TestTask_IsCompleted(t *testing.T) {
	task := &Task{Status: TaskStatusCompleted}
	if !task.IsCompleted() {
		t.Errorf("Task.IsCompleted() = false, want true for completed task")
	}

	task.Status = TaskStatusActive
	if task.IsCompleted() {
		t.Errorf("Task.IsCompleted() = true, want false for active task")
	}
}

func TestTask_Start(t *testing.T) {
	task := &Task{
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	task.Start()

	if task.Status != TaskStatusActive {
		t.Errorf("Task.Start() status = %v, want %v", task.Status, TaskStatusActive)
	}

	if task.StartTime == nil {
		t.Errorf("Task.Start() StartTime = nil, want non-nil")
	}
}

func TestTask_Complete(t *testing.T) {
	task := &Task{
		Status:    TaskStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	startTime := time.Now()
	task.StartTime = &startTime

	task.Complete()

	if task.Status != TaskStatusCompleted {
		t.Errorf("Task.Complete() status = %v, want %v", task.Status, TaskStatusCompleted)
	}

	if task.EndTime == nil {
		t.Errorf("Task.Complete() EndTime = nil, want non-nil")
	}

	if task.CompletedAt == nil {
		t.Errorf("Task.Complete() CompletedAt = nil, want non-nil")
	}

	if task.ActualDuration == 0 {
		t.Errorf("Task.Complete() ActualDuration = 0, want > 0")
	}
}

func TestTask_Pause(t *testing.T) {
	task := &Task{
		Status:    TaskStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	startTime := time.Now().Add(-time.Minute)
	task.StartTime = &startTime

	task.Pause()

	if task.Status != TaskStatusPaused {
		t.Errorf("Task.Pause() status = %v, want %v", task.Status, TaskStatusPaused)
	}

	if task.StartTime != nil {
		t.Errorf("Task.Pause() StartTime = %v, want nil", task.StartTime)
	}

	if task.ActualDuration == 0 {
		t.Errorf("Task.Pause() ActualDuration = 0, want > 0")
	}
}

func TestTask_Resume(t *testing.T) {
	task := &Task{
		Status:         TaskStatusPaused,
		ActualDuration: time.Minute,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	task.Resume()

	if task.Status != TaskStatusActive {
		t.Errorf("Task.Resume() status = %v, want %v", task.Status, TaskStatusActive)
	}

	if task.StartTime == nil {
		t.Errorf("Task.Resume() StartTime = nil, want non-nil")
	}
}

func TestCanonicalHour_IsActive(t *testing.T) {
	hour := &CanonicalHour{
		StartTime: "09:00",
		EndTime:   "12:00",
	}

	tests := []struct {
		name string
		time string
		want bool
	}{
		{
			name: "within range",
			time: "10:30",
			want: true,
		},
		{
			name: "start time",
			time: "09:00",
			want: true,
		},
		{
			name: "end time",
			time: "12:00",
			want: false, // End time is exclusive
		},
		{
			name: "before range",
			time: "08:30",
			want: false,
		},
		{
			name: "after range",
			time: "12:30",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse time for testing
			testTime, err := time.Parse("15:04", tt.time)
			if err != nil {
				t.Fatalf("Failed to parse test time: %v", err)
			}

			// Use today's date with the test time
			now := time.Now()
			testDateTime := time.Date(now.Year(), now.Month(), now.Day(),
				testTime.Hour(), testTime.Minute(), 0, 0, now.Location())

			if got := hour.IsActive(testDateTime); got != tt.want {
				t.Errorf("CanonicalHour.IsActive(%v) = %v, want %v", tt.time, got, tt.want)
			}
		})
	}
}

func TestSchedule_GetCurrentHour(t *testing.T) {
	schedule := &Schedule{
		Hours: []CanonicalHour{
			{
				Name:      "Morning",
				StartTime: "06:00",
				EndTime:   "09:00",
			},
			{
				Name:      "Afternoon",
				StartTime: "12:00",
				EndTime:   "15:00",
			},
		},
	}

	// Test time within first hour
	testTime, _ := time.Parse("15:04", "07:30")
	now := time.Now()
	testDateTime := time.Date(now.Year(), now.Month(), now.Day(),
		testTime.Hour(), testTime.Minute(), 0, 0, now.Location())

	currentHour := schedule.GetCurrentHour(testDateTime)
	if currentHour == nil {
		t.Errorf("Schedule.GetCurrentHour() = nil, want non-nil")
	} else if currentHour.Name != "Morning" {
		t.Errorf("Schedule.GetCurrentHour() name = %v, want %v", currentHour.Name, "Morning")
	}

	// Test time outside any hour
	testTime2, _ := time.Parse("15:04", "10:30")
	testDateTime2 := time.Date(now.Year(), now.Month(), now.Day(),
		testTime2.Hour(), testTime2.Minute(), 0, 0, now.Location())

	currentHour2 := schedule.GetCurrentHour(testDateTime2)
	if currentHour2 != nil {
		t.Errorf("Schedule.GetCurrentHour() = %v, want nil for time outside hours", currentHour2)
	}
}

func TestSchedule_GetHourByName(t *testing.T) {
	schedule := &Schedule{
		Hours: []CanonicalHour{
			{
				Name:      "Morning",
				StartTime: "06:00",
				EndTime:   "09:00",
			},
			{
				Name:      "Afternoon",
				StartTime: "12:00",
				EndTime:   "15:00",
			},
		},
	}

	hour := schedule.GetHourByName("Morning")
	if hour == nil {
		t.Errorf("Schedule.GetHourByName() = nil, want non-nil")
	} else if hour.Name != "Morning" {
		t.Errorf("Schedule.GetHourByName() name = %v, want %v", hour.Name, "Morning")
	}

	hour2 := schedule.GetHourByName("Evening")
	if hour2 != nil {
		t.Errorf("Schedule.GetHourByName() = %v, want nil for non-existent hour", hour2)
	}
}

func TestDailyStats_CompletionRate(t *testing.T) {
	tests := []struct {
		name           string
		totalTasks     int
		completedTasks int
		want           float64
	}{
		{
			name:           "normal completion",
			totalTasks:     10,
			completedTasks: 7,
			want:           70.0,
		},
		{
			name:           "perfect completion",
			totalTasks:     5,
			completedTasks: 5,
			want:           100.0,
		},
		{
			name:           "no completion",
			totalTasks:     3,
			completedTasks: 0,
			want:           0.0,
		},
		{
			name:           "no tasks",
			totalTasks:     0,
			completedTasks: 0,
			want:           0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &DailyStats{
				TotalTasks:     tt.totalTasks,
				CompletedTasks: tt.completedTasks,
			}

			if got := stats.CompletionRate(); got != tt.want {
				t.Errorf("DailyStats.CompletionRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultSchedule(t *testing.T) {
	schedule := GetDefaultSchedule()

	if schedule.Name == "" {
		t.Errorf("GetDefaultSchedule() name is empty")
	}

	if len(schedule.Hours) == 0 {
		t.Errorf("GetDefaultSchedule() has no hours")
	}

	// Check that we have the expected canonical hours
	expectedHours := []string{"Matins", "Lauds", "Prime", "Terce", "Sext", "None", "Vespers", "Compline"}
	if len(schedule.Hours) != len(expectedHours) {
		t.Errorf("GetDefaultSchedule() has %d hours, want %d", len(schedule.Hours), len(expectedHours))
	}

	for i, expectedName := range expectedHours {
		if i < len(schedule.Hours) && schedule.Hours[i].Name != expectedName {
			t.Errorf("GetDefaultSchedule() hour %d name = %v, want %v", i, schedule.Hours[i].Name, expectedName)
		}
	}
}
