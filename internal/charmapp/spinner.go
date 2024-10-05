package charmapp

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var benSpinner = spinner.Spinner{
	Frames: []string{"n", "e", "v", "e", "r", " ", "g", "o", "n", "n", "a", " ", "g", "i", "v", "e", " ", "y", "o", "u", " ", "u", "p"},
	FPS: time.Second / 10,
}

var reverseDot = spinner.Spinner{
	Frames: []string{"⣷ ", "⣯ ", "⣟ ", "⡿ ", "⢿ ", "⣻ ", "⣽ ", "⣾ "},
	FPS:    spinner.Dot.FPS,
}

var (
	// Available spinners
	spinners = []spinner.Spinner{
		benSpinner,
		spinner.Line,
		spinner.Dot,
		reverseDot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
		spinner.Meter,
	}

	textStyle    func(strs ...string) string // lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	helpStyle    func(strs ...string) string // lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type Spinmodel struct {
	index   int
	spinner spinner.Model

	color int

	TextStyle    lipgloss.Style
	SpinnerStyle lipgloss.Style
	HelpStyle    lipgloss.Style
}

func (m Spinmodel) Init() tea.Cmd {
	textStyle = m.TextStyle.Render
	spinnerStyle = m.SpinnerStyle
	helpStyle = m.HelpStyle.Render

	m.color = 69

	return m.spinner.Tick
}

func (m Spinmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "h", "left":
			m.index--
			if m.index < 0 {
				m.index = len(spinners) - 1
			}
			m.ResetSpinner()
			return m, m.spinner.Tick
		case "l", "right":
			m.index++
			if m.index >= len(spinners) {
				m.index = 0
			}
			m.ResetSpinner()
			return m, m.spinner.Tick
		default:
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *Spinmodel) ResetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

func (m Spinmodel) View() (s string) {
	var gap string
	switch m.index {
	case 1:
		gap = ""
	default:
		gap = " "
	}

	s += fmt.Sprintf("\n %s%s%s\n\n", m.spinner.View(), gap, textStyle("Spinning..."))
	s += helpStyle("h/l, ←/→ : change spinner • q: exit\n")
	return
}
