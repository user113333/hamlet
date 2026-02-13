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
		// TODO: move switch into m.ProcessKey
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		m.ProcessKey(msg.String())
	case tea.WindowSizeMsg:
		model.windowWidth = msg.Width
		model.windowHeight = msg.Height
	}

	return m, nil
}

func (m *Model) ProcessKey(key string) {
	if len(m.shiftkey) > 0 {
		switch m.shiftkey {
		case "g":
			m.ProcessKeyAfterg(key)
		}
		m.shiftkey = ""
		return
	}
	switch key {
	case "k", "up":
		m.Up()
	case "j", "down":
		m.Down()
	case "h", "left":
		m.Left()
	case "l", "right":
		m.Right()
	case "g":
		m.shiftkey = "g"
	case "H":
		m.showHiddenItems = !m.showHiddenItems
		m.rememberCursor()
		m.updateEntries()
		m.restoreCursor()
	}
}

func (m *Model) ProcessKeyAfterg(key string) {
	switch key {
	case "h":
		m.ChangeWDHome()
	}
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
	// TODO: cleanup
	for _, entry := range m.cwdEntries {
		if !entry.isDir {
			continue
		}
		if i == limit {
			return true
		}
		if i == m.cursor {
			builder.WriteString(selectedStyle.Render(entry.name))
		} else {
			builder.WriteString(directoryStyle.Render(entry.name))
		}
		builder.WriteString("\n")
		i++
	}
	for _, entry := range m.cwdEntries {
		if entry.isDir {
			continue
		}
		if i == limit {
			return true
		}
		if i == m.cursor {
			builder.WriteString(selectedStyle.Render(entry.name))
		} else {
			builder.WriteString(entry.name)
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
