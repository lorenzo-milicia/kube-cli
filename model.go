package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type errMsg error

type model struct {
	spinner   spinner.Model
	connected bool
	clientset *kubernetes.Clientset
	choice    int
	chosen    bool
	items     []string
	textInput textinput.Model
	quitting  bool
	err       error
}

func initialModel() model {
	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// Text input
	ti := textinput.New()
	ti.Placeholder = "Select a namespace"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{spinner: s, connected: false, chosen: false, choice: 0, textInput: ti}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, connectCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.chosen {
		return testSearchUpdate(msg, m)
	}

	if !m.connected {
		return spinnerUpdate(msg, m)
	}
	if m.connected {
		return choiceUpdate(msg, m)
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	var output string
	if !m.connected {
		output += fmt.Sprintf("%s Getting Kubernetes namespaces...\n", m.spinner.View())
	}
	if m.connected {
		output += "✅ Successfully connected to the cluster\n"
	}
	if m.connected && !m.chosen {
		c := m.choice

		tpl := "What to do today?\n"
		tpl += "%s\n"
		choices := fmt.Sprintf(
			"%s\n%s\n%s\n",
			checkbox("Get namespaces list", c == 0),
			checkbox("Select a namespace", c == 1),
			checkbox("I don't know...", c == 2),
		)

		output += fmt.Sprintf(tpl, choices)
	}
	if m.connected && m.chosen {
		var searchedText = m.textInput.Value()
		output += fmt.Sprintf(
			"Select a namespace\n%s\n",
			m.textInput.View(),
		) + "\n"

		var visibleNamespaces []string
		for _, ns := range m.items {
			if searchedText == "" {
				visibleNamespaces = append(visibleNamespaces, fmt.Sprintf("%v\n", ns))
			} else {
				if strings.Contains(ns, searchedText) {
					visibleNamespaces = append(visibleNamespaces, fmt.Sprintf("%v\n", ns))
				}

			}
		}
		for _, n := range visibleNamespaces {
			output += n
		}
	}
	return output
}
