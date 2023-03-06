package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {

	var ops = []NamespaceOperation{
		PodsOperation,
	}

	p := tea.NewProgram(initialModel(ops))
	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
