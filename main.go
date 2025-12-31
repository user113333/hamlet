package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var model Model

var cFlag = flag.Bool("c", false, "on exit save cwd into ~/.config/hamlet/cwd")

func cFlagHandle() {
	if cFlag == nil {
		return
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error getting config directory: %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(configDir+"/hamlet/cwd", []byte(model.cwd+"\n"), 0644)
	if err != nil {
		fmt.Printf("Error writing cwd file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	p := tea.NewProgram(ModelNew(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	model = m.(Model)

	cFlagHandle()
}
