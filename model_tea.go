package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "k", "up":
			m.Up()
		case "j", "down":
			m.Down()
		case "h", "left":
			m.Left()
		case "l", "right":
			m.Right()
		}
	}

	return m, nil
}

var (
	directoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

func (m Model) View() string {
	var builder strings.Builder
	builder.WriteString(m.cwd)
	builder.WriteString("\n")
	i := 0
	for _, entry := range m.cwdDirs {
		if i == m.cursor {
			builder.WriteString(selectedStyle.Render(entry))
		} else {
			builder.WriteString(directoryStyle.Render(entry))
		}
		builder.WriteString("\n")
		i++
	}
	for _, entry := range m.cwdFiles {
		if i == m.cursor {
			builder.WriteString(selectedStyle.Render(entry))
		} else {
			builder.WriteString(entry)
		}
		builder.WriteString("\n")
		i++
	}
	return builder.String()
}
