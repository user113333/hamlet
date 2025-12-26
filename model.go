package main

import (
	"os"
	"path/filepath"
)

type Model struct {
	cwd           string
	cwdDirs       []string
	cwdFiles      []string
	cursor        int
	cursorHistory map[string]int
}

func ModelNew() Model {
	currentDir, _ := filepath.Abs(".") // TODO: err
	m := Model{
		currentDir,
		[]string{},
		[]string{},
		0,
		make(map[string]int),
	}
	m.updateEntries()
	return m
}

func (m *Model) updateEntries() {
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

func (m Model) entriesLength() int {
	return len(m.cwdDirs) + len(m.cwdFiles)
}

func (m *Model) rememberCursor() {
	m.cursorHistory[m.cwd] = m.cursor
}

func (m *Model) restoreCursor() {
	m.cursor = m.cursorHistory[m.cwd]
}

func (m *Model) Up() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) Down() {
	if m.cursor < m.entriesLength()-1 {
		m.cursor++
	}
}

func (m *Model) Left() {
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

func (m *Model) Right() {
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
