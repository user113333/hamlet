package main

import (
	"os"
	"path/filepath"
	"slices"
)

type ModelApp struct {
	cwd           string
	cwdDirs       []string
	cwdFiles      []string
	cursor        int
	cursorHistory map[string]string

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
			make(map[string]string),
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
		return true // TODO: throw error (maybe return)
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
	m.cursor = 0
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
	m.cursorHistory[m.cwd] = m.getEntry(m.cursor)
}

func (m *Model) restoreCursor() {
	target := m.cursorHistory[m.cwd]
	i := slices.Index(m.cwdDirs, target)
	if i == -1 {
		i = slices.Index(m.cwdFiles, target)
		if i != -1 {
			i += len(m.cwdDirs)
		}
	}
	if i != -1 {
		m.cursor = i
	}
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

	// If no cwd restored, position cursor to dir we are comming from
	if m.cursor == 0 {
		i := slices.Index(m.cwdDirs, prevCwd)
		if i >= 0 {
			m.cursor = i
		}
	}
}

func (m *Model) Right() {
	if m.cursor >= m.entriesLength() {
		return
	}

	m.rememberCursor()

	m.cwd = filepath.Join(m.cwd, m.getEntry(m.cursor))

	m.updateEntries()
	m.restoreCursor()
}
