package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var cFlag = flag.Bool("c", false, "on exit print cwd")

func main() {
	var model Model

	flag.Parse()
	p := tea.NewProgram(ModelNew(), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	model = m.(Model)

	if *cFlag {
		fmt.Printf("%s\n", model.cwd)
	}
}
