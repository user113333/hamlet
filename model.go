package main

import (
	"os"
	"path/filepath"
)

type ModelApp struct {
	cwd           string
	cwdDirs       []string
	cwdFiles      []string
	cursor        int
	cursorHistory map[string]int

	showHiddenItems bool
}

type ModelTea struct {
	windowWidth  int
	windowHeight int
}

type Model struct {
	ModelApp
	ModelTea
}

func ModelNew() Model {
	currentDir, _ := filepath.Abs(".") // TODO: err
	m := Model{
		ModelApp{
			currentDir,
			[]string{},
			[]string{},
			0,
			make(map[string]int),
			false,
		},
		ModelTea{
			0,
			0,
		},
	}
	m.updateEntries()
	return m
}

func isItemHidden(name string) bool {
	if len(name) <= 0 {
		return true // ?
	}
	if name[0] == '.' {
		return true
	}
	return false
}

func (m *Model) getEntry(i int) string {
	if i < len(m.cwdDirs) {
		return m.cwdDirs[i]
	} else if i < m.entriesLength() {
		return m.cwdFiles[i-len(m.cwdDirs)]
	}
	return "" // TODO: throw error (maybe return)
}

func (m *Model) updateEntries() {
	m.cwdDirs = nil
	m.cwdFiles = nil
	entries, _ := os.ReadDir(m.cwd) // TODO: err
	for _, entry := range entries {
		if !m.showHiddenItems {
			if isItemHidden(entry.Name()) {
				continue
			}
		}
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

	m.cwd = filepath.Join(m.cwd, m.getEntry(m.cursor))

	m.restoreCursor()
	m.updateEntries()
}
