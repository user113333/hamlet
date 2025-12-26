package main

// A simple program that opens the alternate screen buffer then counts down
// from 5 and then exits.

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
	cwd      string
	cwdDirs  []string
	cwdFiles []string
	cursor   int
}

func modelInit() model {
	currentDir, _ := filepath.Abs(".") // TODO: err
	m := model{
		currentDir,
		[]string{},
		[]string{},
		0,
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
		case "j", "down":
			if m.cursor < m.entriesLength()-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "l", "right":
			if m.cursor < len(m.cwdDirs) {
				m.cwd = filepath.Join(m.cwd, m.cwdDirs[m.cursor])
				m.cursor = 0
				m.updateEntries()
			} else if m.cursor < m.entriesLength() {
				i := m.cursor - len(m.cwdDirs)
				m.cwd = filepath.Join(m.cwd, m.cwdFiles[i])
				m.cursor = 0
				m.updateEntries()
			}
		case "h", "left":
			m.cwd = filepath.Join(m.cwd, "..")
			m.cursor = 0
			m.updateEntries()
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
