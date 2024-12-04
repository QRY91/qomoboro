package main

import (
	"fmt"
	"qomoboro/pomodoro"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding             = 2
	maxWidth            = 80
	defaultPomoType     = "work"
	defaultPomoDuration = 2
	defaultPomoStart    = "assets/sounds/start_beep.wav"
	defaultPomoStop     = "assets/sounds/stop_beep.wav"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func main() {
	t := pomodoro.New(defaultPomoType, defaultPomoDuration, defaultPomoStart, defaultPomoStop)
	t.Repeat(3)
	// p := tea.NewProgram(initialModel())
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}

type (
	errMsg error
)

type Model struct {
	// progress  progress.Model
	// textInput textinput.Model
	// err error
}

func initialModel() Model {
	// pr := progress.New(progress.WithDefaultGradient())
	// ti := textinput.New()
	// ti.Placeholder = "Pikachu"
	// ti.Focus()
	// ti.CharLimit = 156
	// ti.Width = 20

	return Model{
		// progress:  pr,
		// textInput: ti,
		// err:       nil,
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
