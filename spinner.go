package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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

type spinmodel struct {
	index   int
	spinner spinner.Model

	color int

	textStyle    lipgloss.Style
	spinnerStyle lipgloss.Style
	helpStyle    lipgloss.Style
}

func (m spinmodel) Init() tea.Cmd {
	textStyle = m.textStyle.Render
	spinnerStyle = m.spinnerStyle
	helpStyle = m.helpStyle.Render

	m.color = 69

	return m.spinner.Tick
}

func (m spinmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.resetSpinner()
			return m, m.spinner.Tick
		case "l", "right":
			m.index++
			if m.index >= len(spinners) {
				m.index = 0
			}
			m.resetSpinner()
			return m, m.spinner.Tick
		default:
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd

		m.color = (m.color + 1) % 255
		log.Info("Color", "m.color", m.color)
		m.spinner.Style.Foreground(lipgloss.Color(strconv.Itoa(m.color)))

		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *spinmodel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

func (m spinmodel) View() (s string) {
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
