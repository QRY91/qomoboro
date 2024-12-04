package main

import (
	"fmt"
	"log"
	"os"
	employee "qomoboro/pomodoro"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// fresh
type Qomoboro struct {
	Type     string
	Duration int64
}

// un-fresh

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func main() {
	// e := employee.Employee{
	// 	FirstName:   "Sam",
	// 	LastName:    "Adolf",
	// 	TotalLeaves: 30,
	// 	LeavesTaken: 20,
	// }
	// e.LeavesRemaining()

	t := employee.Pomodoro{
		Type:     "Kos",
		Duration: 10,
	}
	t.TimeRemaining()
	// p := tea.NewProgram(initialModel())
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
	// pomodoro()
}

type (
	errMsg error
)

type Model struct {
	progress  progress.Model
	textInput textinput.Model
	err       error
}

func initialModel() Model {
	pr := progress.New(progress.WithDefaultGradient())
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return Model{
		progress:  pr,
		textInput: ti,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
func pomodoro() {
	workDuration := 1 * time.Second
	breakDuration := 1 * time.Second
	numPomodoros := 2

	for i := 0; i < numPomodoros; i++ {
		fmt.Println("Starting work session...")
		playSound("assets/sounds/start_beep.wav")
		startTimer(workDuration)
		playSound("assets/sounds/stop_beep.wav")

		fmt.Println("Starting break session...")
		playSound("assets/sounds/start_beep.wav")
		startTimer(breakDuration)
		playSound("assets/sounds/stop_beep.wav")
	}

	fmt.Println("All done! Good job!")
}

func startTimer(duration time.Duration) {
	timer := time.NewTimer(duration)

	<-timer.C
	fmt.Println("Time's up!")
}

func playSound(soundFile string) {
	f, err := os.Open(soundFile)
	if err != nil {
		panic(err)
	}

	s, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		panic(err)
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		done <- true
	})))
	<-done
}
