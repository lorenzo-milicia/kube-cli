package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	operations "go.lorenzomilicia.dev/kube-gum-cli/ops"
	"log"
	"os"
	"plugin"
)

func main() {
	//	args := os.Args
	ops := importPlugins([]string{"./plugin/pods.so"})

	p := tea.NewProgram(initialModel(ops))
	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func importPlugins(paths []string) []operations.NamespaceOperation {
	var ops []operations.NamespaceOperation
	for _, path := range paths {
		p, err := plugin.Open(path)
		if err != nil {
			log.Fatalf("Unable to open plugin %s. error: %v", path, err)
		}
		o, err := p.Lookup("Operation")
		if err != nil {
			log.Fatalf("Unable to find function. error: %v", err)
		}
		ops = append(ops, *o.(*operations.NamespaceOperation))
	}
	return ops
}
