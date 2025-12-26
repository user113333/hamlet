package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	directoryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

type model struct {
	cwd           string
	cwdDirs       []string
	cwdFiles      []string
	cursor        int
	cursorHistory map[string]int
}

func modelInit() model {
	currentDir, _ := filepath.Abs(".") // TODO: err
	m := model{
		currentDir,
		[]string{},
		[]string{},
		0,
		make(map[string]int),
	}
	m.updateEntries()
	return m
}

func (m *model) updateEntries() {
	m.cwdDirs = nil
	m.cwdFiles = nil
	entries, _ := os.ReadDir(m.cwd) // TODO: err
	for _, entry := range entries {
		if entry.IsDir() {
			m.cwdDirs = append(m.cwdDirs, entry.Name())
		} else {
			m.cwdFiles = append(m.cwdFiles, entry.Name())
		}
	}
}

func (m model) entriesLength() int {
	return len(m.cwdDirs) + len(m.cwdFiles)
}

func (m *model) rememberCursor() {
	m.cursorHistory[m.cwd] = m.cursor
}

func (m *model) restoreCursor() {
	m.cursor = m.cursorHistory[m.cwd]
}

func (m *model) Up() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *model) Down() {
	if m.cursor < m.entriesLength()-1 {
		m.cursor++
	}
}

func (m *model) Left() {
	m.rememberCursor()
	prevCwd := filepath.Base(m.cwd)

	m.cwd = filepath.Join(m.cwd, "..")

	m.updateEntries()
	m.restoreCursor()
	if m.cursor == 0 {
		for i, entry := range m.cwdDirs {
			if prevCwd == entry {
				m.cursor = i
			}
		}
	}
}

func (m *model) Right() {
	if m.cursor >= m.entriesLength() {
		return
	}

	m.rememberCursor()

	if m.cursor < len(m.cwdDirs) {
		m.cwd = filepath.Join(m.cwd, m.cwdDirs[m.cursor])
	} else if m.cursor < m.entriesLength() {
		i := m.cursor - len(m.cwdDirs)
		m.cwd = filepath.Join(m.cwd, m.cwdFiles[i])
	}

	m.restoreCursor()
	m.updateEntries()
}

func main() {
	p := tea.NewProgram(modelInit(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
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

var i int = 0

func (m model) View() string {
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
