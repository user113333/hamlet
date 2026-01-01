package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
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
		case "H":
			m.showHiddenItems = !m.showHiddenItems
			m.rememberCursor()
			m.updateEntries()
			m.restoreCursor()
		}
	case tea.WindowSizeMsg:
		model.windowWidth = msg.Width
		model.windowHeight = msg.Height
	}

	return m, nil
}

var (
	directoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

// Appends to strings.Builder cwd files,dir
// Stops at limit
// returns: has been limit reached
func (m Model) AppendCWDItemsToBuilder(builder *strings.Builder, limit int) bool {
	i := 0
	for _, entry := range m.cwdDirs {
		if i == limit {
			return true
		}
		if i == m.cursor {
			builder.WriteString(selectedStyle.Render(entry))
		} else {
			builder.WriteString(directoryStyle.Render(entry))
		}
		builder.WriteString("\n")
		i++
	}
	for _, entry := range m.cwdFiles {
		if i == limit {
			return true
		}
		if i == m.cursor {
			builder.WriteString(selectedStyle.Render(entry))
		} else {
			builder.WriteString(entry)
		}
		builder.WriteString("\n")
		i++
	}
	return false
}

func (m Model) View() string {
	var builder strings.Builder
	builder.WriteString(m.cwd)
	builder.WriteString("\n")
	isLimitReached := m.AppendCWDItemsToBuilder(&builder, model.windowHeight-2)
	if isLimitReached {
		builder.WriteString("...")
	}
	return builder.String()
}
