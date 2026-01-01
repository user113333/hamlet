package main

import (
	"os"
	"path/filepath"
	"slices"
)

type ModelApp struct {
	cwd           string
	cwdEntries    []Entry
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

type Entry struct {
	name  string
	isDir bool
}

func ModelNew() Model {
	currentDir, _ := filepath.Abs(".") // TODO: err
	m := Model{
		ModelApp{
			currentDir,
			[]Entry{},
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

func (m *Model) getEntry(i int) Entry {
	if i < len(m.cwdEntries) {
		return m.cwdEntries[i]
	}
	return Entry{} // TODO: throw error (maybe return)
}

func (m *Model) updateEntries() {
	m.cwdEntries = nil
	m.cursor = 0
	entries, _ := os.ReadDir(m.cwd) // TODO: err
	for _, entry := range entries {
		if !m.showHiddenItems {
			if isItemHidden(entry.Name()) {
				continue
			}
		}
		m.cwdEntries = append(m.cwdEntries, Entry{entry.Name(), entry.IsDir()})
	}
}

func (m *Model) rememberCursor() {
	m.cursorHistory[m.cwd] = m.getEntry(m.cursor).name
}

func (m *Model) restoreCursor() {
	target := m.cursorHistory[m.cwd]
	i := slices.IndexFunc(m.cwdEntries, func(e Entry) bool {
		return e.name == target
	})
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
	if m.cursor < len(m.cwdEntries)-1 {
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
		i := slices.IndexFunc(m.cwdEntries, func(e Entry) bool {
			return prevCwd == e.name && e.isDir
		})
		if i >= 0 {
			m.cursor = i
		}
	}
}

func (m *Model) Right() {
	if m.cursor >= len(m.cwdEntries) {
		return
	}

	m.rememberCursor()

	m.cwd = filepath.Join(m.cwd, m.getEntry(m.cursor).name)

	m.updateEntries()
	m.restoreCursor()
}
