package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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

var (
	green  = lipgloss.AdaptiveColor{Light: "#688060", Dark: "#688060"}
	red    = lipgloss.AdaptiveColor{Light: "#DCA3A3", Dark: "#DCA3A3"}
	blue   = lipgloss.AdaptiveColor{Light: "#8CD0D3", Dark: "#8CD0D3"}
	yellow = lipgloss.AdaptiveColor{Light: "#F0DFAF", Dark: "#F0DFAF"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(yellow).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(green).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(blue)
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type state int

const (
	statusNormal state = iota
	stateDone
)

func main() {
	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}

type Model struct {
	state     state
	lg        *lipgloss.Renderer
	styles    *Styles
	form      *huh.Form
	width     int
	worktime  string
	breaktime string
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("QOMOBORO").
				Description("pomodoro\n\n\n\n").
				Next(true).
				NextLabel("next"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("worktime").
				Options(huh.NewOptions("15 min", "20 min", "30 min")...).
				Title("choose work duration").
				Description("how long before a break?").
				Value(&m.worktime),

			huh.NewSelect[string]().
				Key("breaktime").
				Options(huh.NewOptions("5", "10", "15")...).
				Title("choose break duration").
				Description("how long to regain focus?").
				Value(&m.breaktime),

			huh.NewConfirm().
				Key("done").
				Title("ready?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("welp, finish up then")
					}
					return nil
				}).
				Affirmative("yes").
				Negative("no"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := m.styles

	var worktime string
	if m.form.GetString("worktime") != "" {
		worktime = "work for " + m.form.GetString("worktime")
	}
	var breaktime string
	if m.form.GetString("breaktime") != "" {
		breaktime = "rest for " + m.form.GetString("breaktime") + " minutes"
	}

	switch m.form.State {
	case huh.StateCompleted:
		var b strings.Builder
		fmt.Fprintf(&b, "Imagine working for %s.\n", worktime)
		fmt.Fprintf(&b, "Now imagine resting for %s.\n\nCongratulations!\nYou've imagined yourself useful.\n", breaktime)

		// shaketime, strErr := strconv.Atoi(breaktime)
		// if strErr != nil {
		// 	panic(strErr)
		// }

		// t := pomodoro.New(defaultPomoType, shaketime, defaultPomoStart, defaultPomoStop)
		// t.Repeat(3)
		return s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String()) + "\n\n"
	default:

		// Form (left side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		// Status (right side)
		var status string
		{
			var (
				qomoInfo = "(None)"
			)

			if m.form.GetString("worktime") != "" {
				qomoInfo = fmt.Sprintf("%s\n%s", worktime, breaktime)
			}

			const statusWidth = 28
			statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
			status = s.Status.
				Height(lipgloss.Height(form)).
				Width(statusWidth).
				MarginLeft(statusMarginLeft).
				Render(s.StatusHeader.Render("qurrent qomo qonfig") + "\n" + qomoInfo)
		}

		errors := m.form.Errors()
		header := m.appBoundaryView("qomoboro")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(green),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}
