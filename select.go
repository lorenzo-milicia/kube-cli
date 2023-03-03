package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func checkbox(label string, checked bool) string {
	if checked {
		return "[x] " + label
	}
	return fmt.Sprintf("[ ] %s", label)
}

func choiceUpdate(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.choice++
			if m.choice > 2 {
				m.choice = 2
			}
		case "k", "up":
			m.choice--
			if m.choice < 0 {
				m.choice = 0
			}
		case "enter", "q", "esc":
			m.chosen = true
			ns, err := m.getNamespaces()
			if err != nil {
				panic(err)
			}
			for _, ns := range ns.Items {
				m.items = append(m.items, ns.Name)
			}
			return m, textinput.Blink
		}
	}
	return m, nil
}
